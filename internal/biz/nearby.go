package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"ito-deposit/internal/pkg/baidumap"
	"ito-deposit/internal/pkg/geo"
	"math"
)

// NearbyLockerPoint 附近寄存点信息
type NearbyLockerPoint struct {
	ID        int32   `json:"id"`        // 寄存点ID
	Name      string  `json:"name"`      // 寄存点名称
	Address   string  `json:"address"`   // 地址
	Distance  float64 `json:"distance"`  // 距离（米或公里）
	Longitude float64 `json:"longitude"` // 经度
	Latitude  float64 `json:"latitude"`  // 纬度
}

// MapPointInfo 地图点位信息
type MapPointInfo struct {
	ID             int32   `json:"id"`              // 寄存点ID
	Name           string  `json:"name"`            // 寄存点名称
	Address        string  `json:"address"`         // 地址
	Longitude      float64 `json:"longitude"`       // 经度
	Latitude       float64 `json:"latitude"`        // 纬度
	TotalAvailable int32   `json:"total_available"` // 总可用柜数量
	Status         string  `json:"status"`          // 状态
}

// ClusterInfo 聚合点信息
type ClusterInfo struct {
	Longitude      float64 `json:"longitude"`       // 聚合点经度
	Latitude       float64 `json:"latitude"`        // 聚合点纬度
	Count          int32   `json:"count"`           // 聚合点数量
	TotalAvailable int32   `json:"total_available"` // 聚合点总可用柜数量
	PointIds       []int32 `json:"point_ids"`       // 聚合的寄存点ID列表
}

// MapData 地图数据
type MapData struct {
	Points      []*MapPointInfo `json:"points"`       // 详细点位列表
	Clusters    []*ClusterInfo  `json:"clusters"`     // 聚合点列表
	TotalCount  int32           `json:"total_count"`  // 总寄存点数量
	ZoomLevel   int32           `json:"zoom_level"`   // 当前缩放级别
	IsClustered bool            `json:"is_clustered"` // 是否返回聚合数据
}

// NearbyRepo 附近寄存点数据仓库接口
type NearbyRepo interface {
	// GetLockerPoints 获取所有寄存点
	GetLockerPoints(ctx context.Context) ([]*LockerPoint, error)
	// GetLockerPointByID 根据ID获取寄存点详情
	GetLockerPointByID(ctx context.Context, id int32) (*LockerPoint, error)
	// SearchLockerPointsInCity 搜索指定城市内的寄存点
	SearchLockerPointsInCity(ctx context.Context, cityName string, keyword string, page, pageSize int64) ([]*LockerPoint, int64, error)
	// GetAllLockerPoints 获取所有寄存点（不依赖城市表）
	GetAllLockerPoints(ctx context.Context, keyword string, page, pageSize int64) ([]*LockerPoint, int64, error)
	// GetLockerPointsInBounds 获取指定边界内的寄存点
	GetLockerPointsInBounds(ctx context.Context, cityName string, northLat, southLat, eastLng, westLng float64) ([]*LockerPoint, error)
}

// NearbyUsecase 附近寄存点用例
type NearbyUsecase struct {
	repo        NearbyRepo
	geoSvc      *geo.GeoService
	baiduMap    *baidumap.BaiduMapClient
	cityUsecase *CityUsecase
	log         *log.Helper
}

// NewNearbyUsecase 创建附近寄存点用例实例
func NewNearbyUsecase(repo NearbyRepo, geoSvc *geo.GeoService, cityUsecase *CityUsecase, logger log.Logger) *NearbyUsecase {
	// 创建百度地图客户端，使用提供的AK
	baiduMapClient := baidumap.NewBaiduMapClient("7pzoTHchDdMRK7jmpCr1sugjv3hfoxz5")

	return &NearbyUsecase{
		repo:        repo,
		geoSvc:      geoSvc,
		baiduMap:    baiduMapClient,
		cityUsecase: cityUsecase,
		log:         log.NewHelper(logger),
	}
}

// InitLockerPointsGeo 初始化寄存点地理位置数据
func (uc *NearbyUsecase) InitLockerPointsGeo(ctx context.Context) error {
	// 获取所有寄存点
	lockerPoints, err := uc.repo.GetLockerPoints(ctx)
	if err != nil {
		uc.log.Errorf("获取寄存点失败: %v", err)
		return err
	}

	// 转换为GEO服务需要的格式
	geoLockers := make([]geo.LockerPointInfo, 0, len(lockerPoints))
	for _, point := range lockerPoints {
		geoLockers = append(geoLockers, geo.LockerPointInfo{
			ID:        point.Id,
			Name:      point.Name,
			Address:   point.Address,
			Longitude: point.Longitude,
			Latitude:  point.Latitude,
		})
	}

	// 批量保存到Redis GEO
	err = uc.geoSvc.SaveAllLockerPoints(ctx, geoLockers)
	if err != nil {
		uc.log.Errorf("保存寄存点位置信息失败: %v", err)
		return err
	}

	uc.log.Infof("成功初始化 %d 个寄存点的地理位置数据", len(geoLockers))
	return nil
}

// FindNearbyLockerPoints 查找附近的寄存点
func (uc *NearbyUsecase) FindNearbyLockerPoints(ctx context.Context, longitude, latitude float64, radius float64, limit int64) ([]*NearbyLockerPoint, error) {
	// 使用GEO服务查询附近的寄存点
	locations, err := uc.geoSvc.FindNearbyLockerPoints(ctx, longitude, latitude, radius, "km", limit)
	if err != nil {
		uc.log.Errorf("查询附近寄存点失败: %v", err)
		return nil, err
	}

	// 转换为业务实体
	result := make([]*NearbyLockerPoint, 0, len(locations))
	for _, loc := range locations {
		// 获取寄存点详情
		point, err := uc.repo.GetLockerPointByID(ctx, loc.ID)
		if err != nil {
			uc.log.Warnf("获取寄存点详情失败: %v", err)
			continue
		}

		// 创建附近寄存点对象
		nearby := &NearbyLockerPoint{
			ID:        loc.ID,
			Name:      loc.Name,
			Address:   point.Address,
			Distance:  loc.Distance,
			Longitude: loc.Longitude,
			Latitude:  loc.Latitude,
		}

		result = append(result, nearby)
	}

	return result, nil
}

// GetLocationByIP 根据IP地址获取位置信息
func (uc *NearbyUsecase) GetLocationByIP(ctx context.Context, ip string) (float64, float64, string, error) {
	// 调用百度地图API获取位置信息
	location, err := uc.baiduMap.GetLocation(ip)
	if err != nil {
		uc.log.Errorf("获取位置信息失败: %v", err)
		return 0, 0, "", err
	}

	// 返回经纬度和城市名称
	return location.Content.Point.X, location.Content.Point.Y, location.Content.AddressDetail.City, nil
}

// GetRealtimeLocation 获取实时位置信息
func (uc *NearbyUsecase) GetRealtimeLocation(ctx context.Context, cityCode, ip string) (float64, float64, string, error) {
	// 调用百度地图API获取实时位置信息
	location, err := uc.baiduMap.GetRealtimeLocation(cityCode, ip, "")
	if err != nil {
		uc.log.Errorf("获取实时位置信息失败: %v", err)
		return 0, 0, "", err
	}

	// 返回经纬度和城市名称
	return location.Result.Location.Lng, location.Result.Location.Lat, location.Result.City, nil
}

// GetLocationByCityName 根据城市名称获取位置信息
func (uc *NearbyUsecase) GetLocationByCityName(ctx context.Context, cityName string) (float64, float64, error) {
	// 根据城市名称获取城市
	city, err := uc.GetCityByName(ctx, cityName)
	if err != nil {
		uc.log.Errorf("获取城市信息失败: %v", err)
		return 0, 0, err
	}

	// 返回城市中心点经纬度
	return city.Longitude, city.Latitude, nil
}

// GetUserLocationInCity 获取用户在指定城市内的位置
func (uc *NearbyUsecase) GetUserLocationInCity(ctx context.Context, cityName string, userLongitude, userLatitude float64, ip string, useRealtime bool) (float64, float64, error) {
	// 如果提供了用户位置，则直接使用
	if userLongitude != 0 && userLatitude != 0 {
		// 验证用户位置是否在城市范围内
		// 这里可以添加验证逻辑，但为简单起见，我们直接返回用户位置
		return userLongitude, userLatitude, nil
	}

	// 获取城市信息，用于获取城市编码
	city, err := uc.GetCityByName(ctx, cityName)
	if err != nil {
		uc.log.Errorf("获取城市信息失败: %v", err)
		return 0, 0, err
	}

	cityCode := city.Code

	// 如果启用了实时定位且提供了IP地址，则使用实时定位
	if useRealtime && ip != "" {
		lng, lat, _, err := uc.GetRealtimeLocation(ctx, cityCode, ip)
		if err == nil {
			// 保存用户实时位置到Redis
			userID := "user:" + ip // 使用IP作为用户ID
			if err := uc.geoSvc.SaveUserLocation(ctx, userID, lng, lat); err != nil {
				uc.log.Warnf("保存用户位置失败: %v", err)
			}
			return lng, lat, nil
		}
		// 如果实时定位失败，记录日志但继续使用城市中心点
		uc.log.Warnf("实时定位失败，将使用城市中心点: %v", err)
	} else if ip != "" {
		// 如果没有启用实时定位但提供了IP，尝试使用IP定位
		lng, lat, _, err := uc.GetLocationByIP(ctx, ip)
		if err == nil {
			// 保存用户位置到Redis
			userID := "user:" + ip // 使用IP作为用户ID
			if err := uc.geoSvc.SaveUserLocation(ctx, userID, lng, lat); err != nil {
				uc.log.Warnf("保存用户位置失败: %v", err)
			}
			return lng, lat, nil
		}
		// 如果IP定位失败，记录日志但继续使用城市中心点
		uc.log.Warnf("IP定位失败，将使用城市中心点: %v", err)
	}

	// 如果没有提供用户位置或定位失败，则使用城市中心点
	return uc.GetLocationByCityName(ctx, cityName)
}

// GetCityByName 根据城市名称获取城市信息
func (uc *NearbyUsecase) GetCityByName(ctx context.Context, cityName string) (*City, error) {
	// 调用城市服务获取城市信息
	return uc.cityUsecase.GetUserCityByName(ctx, cityName)
}

// UserLocationInfo 用户位置信息
type UserLocationInfo struct {
	Longitude    float64 `json:"longitude"`     // 用户经度
	Latitude     float64 `json:"latitude"`      // 用户纬度
	Address      string  `json:"address"`       // 详细地址
	City         string  `json:"city"`          // 城市名称
	District     string  `json:"district"`      // 区县
	Province     string  `json:"province"`      // 省份
	LocationType string  `json:"location_type"` // 定位类型
}

// MyNearbyInfo 我的附近信息
type MyNearbyInfo struct {
	UserLocation *UserLocationInfo    `json:"user_location"` // 用户位置信息
	NearbyPoints []*NearbyLockerPoint `json:"nearby_points"` // 附近寄存点列表
	TotalCount   int32                `json:"total_count"`   // 附近寄存点总数
	SearchRadius float64              `json:"search_radius"` // 实际搜索半径
	BaiduMapAK   string               `json:"baidu_map_ak"`  // 百度地图AK
}

// GetMyNearbyInfo 获取我的附近信息（实时位置和附近寄存点）
func (uc *NearbyUsecase) GetMyNearbyInfo(ctx context.Context, ip string, userLng, userLat, radius float64, limit int64) (*MyNearbyInfo, error) {
	// 设置默认值
	if radius <= 0 {
		radius = 5.0 // 默认5公里
	}
	if limit <= 0 {
		limit = 20 // 默认20个
	}

	var userLocation *UserLocationInfo
	var err error

	// 获取用户位置信息
	if userLng != 0 && userLat != 0 {
		// 如果提供了GPS坐标，优先使用
		userLocation = &UserLocationInfo{
			Longitude:    userLng,
			Latitude:     userLat,
			Address:      "GPS定位",
			City:         "",
			District:     "",
			Province:     "",
			LocationType: "gps",
		}
	} else if ip != "" {
		// 使用IP定位
		userLocation, err = uc.getUserLocationByIP(ctx, ip)
		if err != nil {
			uc.log.Warnf("IP定位失败: %v，将使用默认位置", err)
			// 如果IP定位失败，使用北京作为默认位置
			userLocation = uc.getDefaultLocation()
		}
	} else {
		// 没有提供任何位置信息，使用默认位置
		userLocation = uc.getDefaultLocation()
	}

	// 查找附近的寄存点
	nearbyPoints, err := uc.FindNearbyLockerPoints(ctx, userLocation.Longitude, userLocation.Latitude, radius, limit)
	if err != nil {
		uc.log.Errorf("查找附近寄存点失败: %v", err)
		return nil, err
	}

	// 构建响应
	myNearbyInfo := &MyNearbyInfo{
		UserLocation: userLocation,
		NearbyPoints: nearbyPoints,
		TotalCount:   int32(len(nearbyPoints)),
		SearchRadius: radius,
		BaiduMapAK:   "7pzoTHchDdMRK7jmpCr1sugjv3hfoxz5", // 百度地图AK，前端使用
	}

	return myNearbyInfo, nil
}

// getUserLocationByIP 通过IP获取用户位置信息
func (uc *NearbyUsecase) getUserLocationByIP(ctx context.Context, ip string) (*UserLocationInfo, error) {
	// 调用百度地图API获取位置信息
	location, err := uc.baiduMap.GetLocation(ip)
	if err != nil {
		return nil, err
	}

	return &UserLocationInfo{
		Longitude:    location.Content.Point.X,
		Latitude:     location.Content.Point.Y,
		Address:      location.Content.Address,
		City:         location.Content.AddressDetail.City,
		District:     location.Content.AddressDetail.District,
		Province:     location.Content.AddressDetail.Province,
		LocationType: "ip",
	}, nil
}

// getDefaultLocation 获取默认位置（北京天安门）
func (uc *NearbyUsecase) getDefaultLocation() *UserLocationInfo {
	return &UserLocationInfo{
		Longitude:    116.397428, // 天安门经度
		Latitude:     39.90923,   // 天安门纬度
		Address:      "北京市东城区天安门广场",
		City:         "北京市",
		District:     "东城区",
		Province:     "北京市",
		LocationType: "default",
	}
}

// SearchLockerPointsInCity 搜索指定城市内的寄存点
func (uc *NearbyUsecase) SearchLockerPointsInCity(ctx context.Context, cityName string, keyword string, page, pageSize int64) ([]*LockerPoint, int64, error) {
	// 参数验证
	if cityName == "" {
		return nil, 0, fmt.Errorf("城市名称不能为空")
	}

	// 设置默认值
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 验证城市是否存在
	_, err := uc.GetCityByName(ctx, cityName)
	if err != nil {
		uc.log.Errorf("获取城市信息失败: %v", err)
		return nil, 0, err
	}

	// 调用数据仓库搜索寄存点
	return uc.repo.SearchLockerPointsInCity(ctx, cityName, keyword, page, pageSize)
}

// GetAllLockerPoints 获取所有寄存点（不依赖城市表）
func (uc *NearbyUsecase) GetAllLockerPoints(ctx context.Context, keyword string, page, pageSize int64) ([]*LockerPoint, int64, error) {
	// 设置默认值
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	
	uc.log.Infof("获取所有寄存点，关键词: %s, 页码: %d, 页大小: %d", keyword, page, pageSize)
	
	// 调用数据仓库获取所有寄存点
	return uc.repo.GetAllLockerPoints(ctx, keyword, page, pageSize)
}

// GetCityLockerPointsMap 获取城市寄存点分布图数据
func (uc *NearbyUsecase) GetCityLockerPointsMap(ctx context.Context, cityName string, northLat, southLat, eastLng, westLng float64, zoomLevel int32, enableCluster bool) (*MapData, error) {
	// 参数验证
	if cityName == "" {
		return nil, fmt.Errorf("城市名称不能为空")
	}

	if northLat <= southLat {
		return nil, fmt.Errorf("北纬度必须大于南纬度")
	}

	if eastLng <= westLng {
		return nil, fmt.Errorf("东经度必须大于西经度")
	}

	if zoomLevel < 1 || zoomLevel > 20 {
		return nil, fmt.Errorf("缩放级别必须在1-20之间")
	}

	// 验证城市是否存在
	_, err := uc.GetCityByName(ctx, cityName)
	if err != nil {
		uc.log.Errorf("获取城市信息失败: %v", err)
		return nil, err
	}

	// 获取边界内的寄存点
	lockerPoints, err := uc.repo.GetLockerPointsInBounds(ctx, cityName, northLat, southLat, eastLng, westLng)
	if err != nil {
		uc.log.Errorf("获取边界内寄存点失败: %v", err)
		return nil, err
	}

	// 转换为地图点位信息
	mapPoints := make([]*MapPointInfo, 0, len(lockerPoints))
	for _, point := range lockerPoints {
		totalAvailable := point.AvailableLarge + point.AvailableMedium + point.AvailableSmall
		status := uc.getLockerPointStatus(totalAvailable)

		mapPoints = append(mapPoints, &MapPointInfo{
			ID:             point.Id,
			Name:           point.Name,
			Address:        point.Address,
			Longitude:      point.Longitude,
			Latitude:       point.Latitude,
			TotalAvailable: totalAvailable,
			Status:         status,
		})
	}

	mapData := &MapData{
		TotalCount:  int32(len(mapPoints)),
		ZoomLevel:   zoomLevel,
		IsClustered: false,
	}

	// 根据缩放级别和聚合设置决定是否进行聚合
	if enableCluster && uc.shouldCluster(zoomLevel, len(mapPoints)) {
		// 进行聚合
		clusters := uc.clusterPoints(mapPoints, zoomLevel)
		mapData.Clusters = clusters
		mapData.IsClustered = true
	} else {
		// 返回详细点位
		mapData.Points = mapPoints
		mapData.IsClustered = false
	}

	return mapData, nil
}

// getLockerPointStatus 获取寄存点状态
func (uc *NearbyUsecase) getLockerPointStatus(totalAvailable int32) string {
	if totalAvailable == 0 {
		return "full"
	} else if totalAvailable <= 5 {
		return "busy"
	} else {
		return "available"
	}
}

// shouldCluster 判断是否应该进行聚合
func (uc *NearbyUsecase) shouldCluster(zoomLevel int32, pointCount int) bool {
	// 根据缩放级别和点位数量决定是否聚合
	// 缩放级别越低（数字越小），越容易聚合
	// 点位数量越多，越容易聚合

	if zoomLevel <= 10 {
		return pointCount > 10
	} else if zoomLevel <= 15 {
		return pointCount > 50
	} else {
		return pointCount > 100
	}
}

// clusterPoints 对点位进行聚合
func (uc *NearbyUsecase) clusterPoints(points []*MapPointInfo, zoomLevel int32) []*ClusterInfo {
	if len(points) == 0 {
		return []*ClusterInfo{}
	}

	// 根据缩放级别确定聚合距离（度数）
	clusterDistance := uc.getClusterDistance(zoomLevel)

	clusters := make([]*ClusterInfo, 0)
	used := make([]bool, len(points))

	for i, point := range points {
		if used[i] {
			continue
		}

		// 创建新的聚合点
		cluster := &ClusterInfo{
			Longitude:      point.Longitude,
			Latitude:       point.Latitude,
			Count:          1,
			TotalAvailable: point.TotalAvailable,
			PointIds:       []int32{point.ID},
		}

		used[i] = true

		// 查找附近的点位进行聚合
		for j := i + 1; j < len(points); j++ {
			if used[j] {
				continue
			}

			// 计算距离（使用度数差）
			distance := uc.calculateDistanceInDegrees(point.Longitude, point.Latitude, points[j].Longitude, points[j].Latitude)

			if distance <= clusterDistance {
				// 加入聚合点
				cluster.Count++
				cluster.TotalAvailable += points[j].TotalAvailable
				cluster.PointIds = append(cluster.PointIds, points[j].ID)

				// 更新聚合点中心位置（简单平均）
				cluster.Longitude = (cluster.Longitude*float64(cluster.Count-1) + points[j].Longitude) / float64(cluster.Count)
				cluster.Latitude = (cluster.Latitude*float64(cluster.Count-1) + points[j].Latitude) / float64(cluster.Count)

				used[j] = true
			}
		}

		clusters = append(clusters, cluster)
	}

	return clusters
}

// getClusterDistance 根据缩放级别获取聚合距离
func (uc *NearbyUsecase) getClusterDistance(zoomLevel int32) float64 {
	// 缩放级别越低，聚合距离越大
	switch {
	case zoomLevel <= 5:
		return 0.1 // 约11公里
	case zoomLevel <= 10:
		return 0.05 // 约5.5公里
	case zoomLevel <= 15:
		return 0.01 // 约1.1公里
	default:
		return 0.005 // 约550米
	}
}

// calculateDistance 计算两点之间的距离（使用Haversine公式，返回公里）
func (uc *NearbyUsecase) calculateDistance(lng1, lat1, lng2, lat2 float64) float64 {
	const earthRadius = 6371 // 地球半径，单位：公里

	// 将度数转换为弧度
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLatRad := (lat2 - lat1) * math.Pi / 180
	deltaLngRad := (lng2 - lng1) * math.Pi / 180

	// Haversine公式
	a := math.Sin(deltaLatRad/2)*math.Sin(deltaLatRad/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLngRad/2)*math.Sin(deltaLngRad/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// calculateDistanceInDegrees 计算两点之间的距离（度数差，用于聚合判断）
func (uc *NearbyUsecase) calculateDistanceInDegrees(lng1, lat1, lng2, lat2 float64) float64 {
	dlng := lng1 - lng2
	dlat := lat1 - lat2
	return math.Sqrt(dlng*dlng + dlat*dlat)
}

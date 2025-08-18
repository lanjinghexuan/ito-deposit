package biz

import (
	"context"
	"fmt"
	"ito-deposit/internal/pkg/baidumap"
	"ito-deposit/internal/pkg/geo"
	"math"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
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
		// 验证坐标是否在合理范围内（中国境内）
		if uc.isValidCoordinate(userLongitude, userLatitude) {
			return userLongitude, userLatitude, nil
		} else {
			uc.log.Warnf("用户提供的坐标超出合理范围: lng=%f, lat=%f", userLongitude, userLatitude)
		}
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
		if err == nil && uc.isValidCoordinate(lng, lat) {
			// 验证定位结果是否在目标城市附近（50公里范围内）
			if uc.isNearCity(lng, lat, city.Longitude, city.Latitude, 50.0) {
				// 保存用户实时位置到Redis
				userID := "user:" + ip // 使用IP作为用户ID
				if err := uc.geoSvc.SaveUserLocation(ctx, userID, lng, lat); err != nil {
					uc.log.Warnf("保存用户位置失败: %v", err)
				}
				return lng, lat, nil
			} else {
				uc.log.Warnf("实时定位结果距离目标城市过远: lng=%f, lat=%f, 城市: %s", lng, lat, cityName)
			}
		}
		// 如果实时定位失败，记录日志但继续使用城市中心点
		uc.log.Warnf("实时定位失败，将使用城市中心点: %v", err)
	} else if ip != "" {
		// 如果没有启用实时定位但提供了IP，尝试使用IP定位
		lng, lat, _, err := uc.GetLocationByIP(ctx, ip)
		if err == nil && uc.isValidCoordinate(lng, lat) {
			// 验证定位结果是否在目标城市附近
			if uc.isNearCity(lng, lat, city.Longitude, city.Latitude, 100.0) {
				// 保存用户位置到Redis
				userID := "user:" + ip // 使用IP作为用户ID
				if err := uc.geoSvc.SaveUserLocation(ctx, userID, lng, lat); err != nil {
					uc.log.Warnf("保存用户位置失败: %v", err)
				}
				return lng, lat, nil
			} else {
				uc.log.Warnf("IP定位结果距离目标城市过远: lng=%f, lat=%f, 城市: %s", lng, lat, cityName)
			}
		}
		// 如果IP定位失败，记录日志但继续使用城市中心点
		uc.log.Warnf("IP定位失败，将使用城市中心点: %v", err)
	}

	// 如果没有提供用户位置或定位失败，则使用城市中心点
	return uc.GetLocationByCityName(ctx, cityName)
}

// GetCityByName 根据城市名称获取城市信息
func (uc *NearbyUsecase) GetCityByName(ctx context.Context, cityName string) (*City, error) {
	// 首先尝试精确匹配
	city, err := uc.cityUsecase.GetUserCityByName(ctx, cityName)
	if err == nil {
		return city, nil
	}

	// 如果精确匹配失败，尝试各种变体
	cityVariants := uc.generateCityVariants(cityName)
	uc.log.Infof("尝试城市名称变体匹配: %v", cityVariants)

	for _, variant := range cityVariants {
		if variant != cityName { // 避免重复尝试原始名称
			city, err = uc.cityUsecase.GetUserCityByName(ctx, variant)
			if err == nil {
				uc.log.Infof("城市名称变体匹配成功: %s -> %s", cityName, variant)
				return city, nil
			}
		}
	}

	// 所有尝试都失败，返回原始错误
	uc.log.Errorf("无法找到城市: %s，尝试的变体: %v", cityName, cityVariants)
	return nil, fmt.Errorf("城市 %s 不存在", cityName)
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
	uc.log.Infof("开始获取城市 %s 边界内的寄存点，边界: 北纬=%f, 南纬=%f, 东经=%f, 西经=%f",
		cityName, northLat, southLat, eastLng, westLng)

	lockerPoints, err := uc.repo.GetLockerPointsInBounds(ctx, cityName, northLat, southLat, eastLng, westLng)
	if err != nil {
		uc.log.Errorf("获取边界内寄存点失败: %v", err)
		return nil, err
	}

	uc.log.Infof("边界查询获取到 %d 个寄存点", len(lockerPoints))

	// 如果边界查询没有结果，尝试获取该城市的所有寄存点作为备用
	if len(lockerPoints) == 0 {
		uc.log.Warnf("边界查询无结果，尝试获取城市 %s 的所有寄存点", cityName)
		allPoints, _, err := uc.repo.SearchLockerPointsInCity(ctx, cityName, "", 1, 100)
		if err != nil {
			uc.log.Errorf("获取城市所有寄存点失败: %v", err)
		} else {
			uc.log.Infof("备用查询获取到 %d 个寄存点", len(allPoints))
			// 转换为LockerPoint类型
			for _, point := range allPoints {
				lockerPoints = append(lockerPoints, &LockerPoint{
					Id:              point.Id,
					Name:            point.Name,
					Address:         point.Address,
					Longitude:       point.Longitude,
					Latitude:        point.Latitude,
					AvailableLarge:  point.AvailableLarge,
					AvailableMedium: point.AvailableMedium,
					AvailableSmall:  point.AvailableSmall,
					OpenTime:        point.OpenTime,
					Mobile:          point.Mobile,
				})
			}
		}
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

// isValidCoordinate 验证坐标是否在合理范围内（中国境内）
func (uc *NearbyUsecase) isValidCoordinate(lng, lat float64) bool {
	// 中国大陆经纬度范围 (WGS84坐标系)
	// 经度：73°33′E 至 135°05′E
	// 纬度：3°51′N 至 53°33′N
	if lng < 73.0 || lng > 135.0 {
		uc.log.Warnf("经度超出中国范围: %f (应在73-135之间)", lng)
		return false
	}
	if lat < 3.0 || lat > 54.0 {
		uc.log.Warnf("纬度超出中国范围: %f (应在3-54之间)", lat)
		return false
	}

	// 特别检查是否在海上（常见的坐标系错误）
	if (lng > 100 && lng < 125) && (lat > 0 && lat < 25) {
		// 这个范围可能是南海，需要更精确的验证
		uc.log.Infof("坐标在南海区域，需要验证: lng=%f, lat=%f", lng, lat)
	}

	return true
}

// isNearCity 判断坐标是否在城市附近指定范围内
func (uc *NearbyUsecase) isNearCity(lng, lat, cityLng, cityLat, maxDistanceKm float64) bool {
	distance := uc.calculateDistance(lng, lat, cityLng, cityLat)
	return distance <= maxDistanceKm
}

// convertBD09ToWGS84 将百度坐标系(BD09)转换为WGS84坐标系
func (uc *NearbyUsecase) convertBD09ToWGS84(bdLng, bdLat float64) (float64, float64) {
	// BD09 -> GCJ02
	x := bdLng - 0.0065
	y := bdLat - 0.006
	z := math.Sqrt(x*x+y*y) - 0.00002*math.Sin(y*math.Pi)
	theta := math.Atan2(y, x) - 0.000003*math.Cos(x*math.Pi)
	gcjLng := z * math.Cos(theta)
	gcjLat := z * math.Sin(theta)

	// GCJ02 -> WGS84
	dlat := uc.transformLat(gcjLng-105.0, gcjLat-35.0)
	dlng := uc.transformLng(gcjLng-105.0, gcjLat-35.0)
	radlat := gcjLat / 180.0 * math.Pi
	magic := math.Sin(radlat)
	magic = 1 - 0.00669342162296594323*magic*magic
	sqrtmagic := math.Sqrt(magic)
	dlat = (dlat * 180.0) / ((6378245.0 * (1 - 0.00669342162296594323)) / (magic * sqrtmagic) * math.Pi)
	dlng = (dlng * 180.0) / (6378245.0 / sqrtmagic * math.Cos(radlat) * math.Pi)
	mglat := gcjLat - dlat
	mglng := gcjLng - dlng

	return mglng, mglat
}

// transformLat 纬度转换辅助函数
func (uc *NearbyUsecase) transformLat(lng, lat float64) float64 {
	ret := -100.0 + 2.0*lng + 3.0*lat + 0.2*lat*lat + 0.1*lng*lat + 0.2*math.Sqrt(math.Abs(lng))
	ret += (20.0*math.Sin(6.0*lng*math.Pi) + 20.0*math.Sin(2.0*lng*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lat*math.Pi) + 40.0*math.Sin(lat/3.0*math.Pi)) * 2.0 / 3.0
	ret += (160.0*math.Sin(lat/12.0*math.Pi) + 320*math.Sin(lat*math.Pi/30.0)) * 2.0 / 3.0
	return ret
}

// transformLng 经度转换辅助函数
func (uc *NearbyUsecase) transformLng(lng, lat float64) float64 {
	ret := 300.0 + lng + 2.0*lat + 0.1*lng*lng + 0.1*lng*lat + 0.1*math.Sqrt(math.Abs(lng))
	ret += (20.0*math.Sin(6.0*lng*math.Pi) + 20.0*math.Sin(2.0*lng*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lng*math.Pi) + 40.0*math.Sin(lng/3.0*math.Pi)) * 2.0 / 3.0
	ret += (150.0*math.Sin(lng/12.0*math.Pi) + 300.0*math.Sin(lng/30.0*math.Pi)) * 2.0 / 3.0
	return ret
}

// generateCityVariants 生成城市名称的各种变体用于匹配
func (uc *NearbyUsecase) generateCityVariants(cityName string) []string {
	variants := []string{cityName} // 原始名称

	// 去除常见的后缀和前缀
	cleanName := cityName

	// 去除省份后缀
	if len(cleanName) > 1 && (cleanName[len(cleanName)-3:] == "省" || cleanName[len(cleanName)-3:] == "市") {
		cleanName = cleanName[:len(cleanName)-3]
		variants = append(variants, cleanName)
	}

	// 添加市后缀
	if !strings.HasSuffix(cityName, "市") {
		variants = append(variants, cityName+"市")
	}

	// 处理特殊情况
	switch cleanName {
	case "郑州":
		variants = append(variants, "郑州市")
	case "北京":
		variants = append(variants, "北京市")
	case "上海":
		variants = append(variants, "上海市")
	case "广州":
		variants = append(variants, "广州市")
	case "深圳":
		variants = append(variants, "深圳市")
	case "杭州":
		variants = append(variants, "杭州市")
	case "南京":
		variants = append(variants, "南京市")
	case "武汉":
		variants = append(variants, "武汉市")
	case "成都":
		variants = append(variants, "成都市")
	case "西安":
		variants = append(variants, "西安市")
	}

	// 去重
	uniqueVariants := make([]string, 0)
	seen := make(map[string]bool)
	for _, variant := range variants {
		if !seen[variant] {
			uniqueVariants = append(uniqueVariants, variant)
			seen[variant] = true
		}
	}

	return uniqueVariants
}

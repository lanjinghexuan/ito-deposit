package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "ito-deposit/api/helloworld/v1"
	"ito-deposit/internal/biz"
)

// NearbyService 附近服务实现
type NearbyService struct {
	v1.UnimplementedNearbyServer

	uc  *biz.NearbyUsecase
	log *log.Helper
}

// NewNearbyService 创建附近服务实例
func NewNearbyService(uc *biz.NearbyUsecase, logger log.Logger) *NearbyService {
	return &NearbyService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// InitLockerPointsGeo 初始化寄存点地理位置数据
func (s *NearbyService) InitLockerPointsGeo(ctx context.Context, req *v1.InitLockerPointsGeoRequest) (*v1.InitLockerPointsGeoReply, error) {
	// 调用业务逻辑初始化寄存点地理位置数据
	err := s.uc.InitLockerPointsGeo(ctx)
	if err != nil {
		s.log.Errorf("初始化寄存点地理位置数据失败: %v", err)
		return nil, v1.ErrorInternalError("初始化寄存点地理位置数据失败: %v", err)
	}

	// 构建响应
	return &v1.InitLockerPointsGeoReply{
		Success: true,
	}, nil
}

// 构建寄存点响应的辅助函数
func buildLockerPointsResponse(lockerPoints []*biz.NearbyLockerPoint) []*v1.NearbyLockerPointInfo {
	items := make([]*v1.NearbyLockerPointInfo, 0, len(lockerPoints))
	for _, point := range lockerPoints {
		items = append(items, &v1.NearbyLockerPointInfo{
			Id:        point.ID,
			Name:      point.Name,
			Address:   point.Address,
			Distance:  float32(point.Distance),
			Longitude: float32(point.Longitude),
			Latitude:  float32(point.Latitude),
		})
	}
	return items
}

// 获取默认参数值的辅助函数
func getDefaultParams(radius float64, limit int64) (float64, int64) {
	if radius <= 0 {
		radius = 5.0 // 默认5公里
	}
	if limit <= 0 {
		limit = 10 // 默认10个
	}
	return radius, limit
}

// 查询附近寄存点的核心方法，减少代码重复
func (s *NearbyService) findNearbyLockerPointsCore(ctx context.Context, longitude, latitude float64, radius float64, limit int64) (*v1.FindNearbyLockerPointsReply, error) {
	// 调用业务逻辑查询附近的寄存点
	lockerPoints, err := s.uc.FindNearbyLockerPoints(ctx, longitude, latitude, radius, limit)
	if err != nil {
		s.log.Errorf("查询附近寄存点失败: %v", err)
		return nil, v1.ErrorInternalError("查询附近寄存点失败: %v", err)
	}

	// 构建响应
	return &v1.FindNearbyLockerPointsReply{
		Items: buildLockerPointsResponse(lockerPoints),
	}, nil
}

// FindNearbyLockerPoints 查找附近的寄存点
func (s *NearbyService) FindNearbyLockerPoints(ctx context.Context, req *v1.FindNearbyLockerPointsRequest) (*v1.FindNearbyLockerPointsReply, error) {
	var longitude, latitude float64
	var err error

	// 根据请求类型获取位置信息
	if req.CityName != "" {
		// 根据城市名称获取位置信息
		longitude, latitude, err = s.uc.GetLocationByCityName(ctx, req.CityName)
		if err != nil {
			s.log.Errorf("根据城市名称获取位置信息失败: %v", err)
			return nil, v1.ErrorInternalError("获取位置信息失败: %v", err)
		}
	} else if req.Ip != "" {
		// 根据IP地址获取位置信息
		longitude, latitude, _, err = s.uc.GetLocationByIP(ctx, req.Ip)
		if err != nil {
			s.log.Errorf("根据IP地址获取位置信息失败: %v", err)
			return nil, v1.ErrorInternalError("获取位置信息失败: %v", err)
		}
	} else {
		// 使用请求中提供的经纬度
		longitude = req.Longitude
		latitude = req.Latitude
	}

	// 设置默认值
	radius, limit := getDefaultParams(req.Radius, req.Limit)

	// 使用核心方法查询
	return s.findNearbyLockerPointsCore(ctx, longitude, latitude, radius, limit)
}

// FindNearbyLockerPointsInCity 查找用户在指定城市内附近的寄存点
func (s *NearbyService) FindNearbyLockerPointsInCity(ctx context.Context, req *v1.FindNearbyLockerPointsInCityRequest) (*v1.FindNearbyLockerPointsReply, error) {
	// 参数验证
	if req.CityName == "" {
		return nil, v1.ErrorBadRequest("城市名称不能为空")
	}

	// 获取用户在城市内的位置
	longitude, latitude, err := s.uc.GetUserLocationInCity(ctx, req.CityName, req.Longitude, req.Latitude, req.Ip, req.UseRealtime)
	if err != nil {
		s.log.Errorf("获取用户在城市内的位置失败: %v", err)
		return nil, v1.ErrorInternalError("获取位置信息失败: %v", err)
	}

	// 设置默认值
	radius, limit := getDefaultParams(req.Radius, req.Limit)

	// 使用核心方法查询
	return s.findNearbyLockerPointsCore(ctx, longitude, latitude, radius, limit)
}

// FindMyNearbyLockerPoints 使用实时定位查找我的附近寄存点
func (s *NearbyService) FindMyNearbyLockerPoints(ctx context.Context, req *v1.FindMyNearbyLockerPointsRequest) (*v1.FindNearbyLockerPointsReply, error) {
	// 参数验证
	if req.CityName == "" {
		return nil, v1.ErrorBadRequest("城市名称不能为空")
	}

	// 创建一个FindNearbyLockerPointsInCityRequest，复用现有功能
	inCityReq := &v1.FindNearbyLockerPointsInCityRequest{
		CityName:    req.CityName,
		Longitude:   0, // 不提供经纬度，使用实时定位
		Latitude:    0,
		Radius:      req.Radius,
		Limit:       req.Limit,
		Ip:          req.Ip,
		UseRealtime: true, // 强制使用实时定位
	}

	// 调用已有的方法处理请求
	return s.FindNearbyLockerPointsInCity(ctx, inCityReq)
}

// SearchLockerPointsInCity 搜索指定城市内的寄存点
func (s *NearbyService) SearchLockerPointsInCity(ctx context.Context, req *v1.SearchLockerPointsInCityRequest) (*v1.SearchLockerPointsInCityReply, error) {
	// 参数验证
	if req.CityName == "" {
		return nil, v1.ErrorBadRequest("城市名称不能为空")
	}

	// 设置默认值
	page := req.Page
	if page <= 0 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 调用业务逻辑搜索寄存点
	lockerPoints, total, err := s.uc.SearchLockerPointsInCity(ctx, req.CityName, req.Keyword, page, pageSize)
	if err != nil {
		s.log.Warnf("按城市搜索寄存点失败: %v，尝试获取所有寄存点", err)

		// 如果城市搜索失败，降级到获取所有寄存点
		lockerPoints, total, err = s.uc.GetAllLockerPoints(ctx, req.Keyword, page, pageSize)
		if err != nil {
			s.log.Errorf("获取所有寄存点失败: %v", err)
			return nil, v1.ErrorInternalError("获取寄存点失败: %v", err)
		}

		s.log.Infof("降级成功，获取到 %d 个寄存点", len(lockerPoints))
	}

	// 构建响应
	items := make([]*v1.LockerPointDetail, 0, len(lockerPoints))
	for _, point := range lockerPoints {
		items = append(items, &v1.LockerPointDetail{
			Id:              point.Id,
			Name:            point.Name,
			Address:         point.Address,
			Longitude:       float32(point.Longitude),
			Latitude:        float32(point.Latitude),
			AvailableLarge:  point.AvailableLarge,
			AvailableMedium: point.AvailableMedium,
			AvailableSmall:  point.AvailableSmall,
			OpenTime:        point.OpenTime,
			Mobile:          point.Mobile,
		})
	}

	return &v1.SearchLockerPointsInCityReply{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetCityLockerPointsMap 获取城市寄存点分布图数据
func (s *NearbyService) GetCityLockerPointsMap(ctx context.Context, req *v1.GetCityLockerPointsMapRequest) (*v1.GetCityLockerPointsMapReply, error) {
	// 参数验证
	if req.CityName == "" {
		return nil, v1.ErrorBadRequest("城市名称不能为空")
	}

	if req.NorthLat <= req.SouthLat {
		return nil, v1.ErrorBadRequest("北纬度必须大于南纬度")
	}

	if req.EastLng <= req.WestLng {
		return nil, v1.ErrorBadRequest("东经度必须大于西经度")
	}

	if req.ZoomLevel < 1 || req.ZoomLevel > 20 {
		return nil, v1.ErrorBadRequest("缩放级别必须在1-20之间")
	}

	// 调用业务逻辑获取地图数据
	mapData, err := s.uc.GetCityLockerPointsMap(ctx, req.CityName, req.NorthLat, req.SouthLat, req.EastLng, req.WestLng, req.ZoomLevel, req.EnableCluster)
	if err != nil {
		s.log.Errorf("获取城市寄存点分布图数据失败: %v", err)
		return nil, v1.ErrorInternalError("获取地图数据失败: %v", err)
	}

	// 构建响应
	reply := &v1.GetCityLockerPointsMapReply{
		TotalCount:  mapData.TotalCount,
		ZoomLevel:   mapData.ZoomLevel,
		IsClustered: mapData.IsClustered,
	}

	// 如果是聚合数据
	if mapData.IsClustered {
		clusters := make([]*v1.ClusterInfo, 0, len(mapData.Clusters))
		for _, cluster := range mapData.Clusters {
			clusters = append(clusters, &v1.ClusterInfo{
				Longitude:      cluster.Longitude,
				Latitude:       cluster.Latitude,
				Count:          cluster.Count,
				TotalAvailable: cluster.TotalAvailable,
				PointIds:       cluster.PointIds,
			})
		}
		reply.Clusters = clusters
	} else {
		// 详细点位数据
		points := make([]*v1.MapPointInfo, 0, len(mapData.Points))
		for _, point := range mapData.Points {
			points = append(points, &v1.MapPointInfo{
				Id:             point.ID,
				Name:           point.Name,
				Address:        point.Address,
				Longitude:      point.Longitude,
				Latitude:       point.Latitude,
				TotalAvailable: point.TotalAvailable,
				Status:         point.Status,
			})
		}
		reply.Points = points
	}

	return reply, nil
}

// GetMyNearbyInfo 获取我的附近信息（实时位置和附近寄存点）
func (s *NearbyService) GetMyNearbyInfo(ctx context.Context, req *v1.GetMyNearbyInfoRequest) (*v1.GetMyNearbyInfoReply, error) {
	// 设置默认值
	radius := req.Radius
	if radius <= 0 {
		radius = 5.0 // 默认5公里
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 20 // 默认20个
	}

	// 调用业务逻辑获取我的附近信息
	myNearbyInfo, err := s.uc.GetMyNearbyInfo(ctx, req.Ip, req.Longitude, req.Latitude, radius, limit)
	if err != nil {
		s.log.Errorf("获取我的附近信息失败: %v", err)
		return nil, v1.ErrorInternalError("获取我的附近信息失败: %v", err)
	}

	// 构建用户位置信息响应
	userLocationReply := &v1.UserLocationInfo{
		Longitude:    myNearbyInfo.UserLocation.Longitude,
		Latitude:     myNearbyInfo.UserLocation.Latitude,
		Address:      myNearbyInfo.UserLocation.Address,
		City:         myNearbyInfo.UserLocation.City,
		District:     myNearbyInfo.UserLocation.District,
		Province:     myNearbyInfo.UserLocation.Province,
		LocationType: myNearbyInfo.UserLocation.LocationType,
	}

	// 构建附近寄存点列表响应
	nearbyPointsReply := make([]*v1.NearbyLockerPointInfo, 0, len(myNearbyInfo.NearbyPoints))
	for _, point := range myNearbyInfo.NearbyPoints {
		nearbyPointsReply = append(nearbyPointsReply, &v1.NearbyLockerPointInfo{
			Id:        point.ID,
			Name:      point.Name,
			Address:   point.Address,
			Distance:  float32(point.Distance),
			Longitude: float32(point.Longitude),
			Latitude:  float32(point.Latitude),
		})
	}

	// 构建响应
	reply := &v1.GetMyNearbyInfoReply{
		UserLocation: userLocationReply,
		NearbyPoints: nearbyPointsReply,
		TotalCount:   myNearbyInfo.TotalCount,
		SearchRadius: myNearbyInfo.SearchRadius,
		BaiduMapAk:   myNearbyInfo.BaiduMapAK,
	}

	return reply, nil
}

// Get AllLockerPoints 获取所有寄存点（不依赖城市）
func (s *NearbyService) GetAllLockerPoints(ctx context.Context, req *v1.GetAllLockerPointsRequest) (*v1.SearchLockerPointsInCityReply, error) {
	// 设置默认值
	page := req.Page
	if page <= 0 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize <= 100 {
		pageSize = 100
	}

	s.log.Infof("获取所有寄存点，关键词: %s, 页码: %d, 页大小: %d", req.Keyword, page, pageSize)

	// 调用业务逻辑获取所有寄存点
	lockerPoints, total, err := s.uc.GetAllLockerPoints(ctx, req.Keyword, page, pageSize)
	if err != nil {
		s.log.Errorf("获取所有寄存点失败: %v", err)
		return nil, v1.ErrorInternalError("获取寄存点失败: %v", err)
	}

	s.log.Infof("成功获取 %d 个寄存点，总数: %d", len(lockerPoints), total)

	// 构建响应
	items := make([]*v1.LockerPointDetail, 0, len(lockerPoints))
	for _, point := range lockerPoints {
		items = append(items, &v1.LockerPointDetail{
			Id:              point.Id,
			Name:            point.Name,
			Address:         point.Address,
			Longitude:       float32(point.Longitude),
			Latitude:        float32(point.Latitude),
			AvailableLarge:  point.AvailableLarge,
			AvailableMedium: point.AvailableMedium,
			AvailableSmall:  point.AvailableSmall,
			OpenTime:        point.OpenTime,
			Mobile:          point.Mobile,
		})
	}

	return &v1.SearchLockerPointsInCityReply{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// TestDatabaseConnection 测试数据库连接和查询（临时调试用）
func (s *NearbyService) TestDatabaseConnection(ctx context.Context) error {
	s.log.Info("=== 开始测试数据库连接 ===")

	// 直接调用业务逻辑层
	lockerPoints, total, err := s.uc.GetAllLockerPoints(ctx, "", 1, 10)
	if err != nil {
		s.log.Errorf("测试查询失败: %v", err)
		return err
	}

	s.log.Infof("测试结果: 总数=%d, 返回数量=%d", total, len(lockerPoints))

	for i, point := range lockerPoints {
		if i < 5 { // 只显示前5个
			s.log.Infof("寄存点 %d: ID=%d, Name=%s, Address=%s", i+1, point.Id, point.Name, point.Address)
		}
	}

	s.log.Info("=== 数据库连接测试完成 ===")
	return nil
}

package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "ito-deposit/api/helloworld/v1"
	"ito-deposit/internal/biz"
)

// CityService 城市服务实现
type CityService struct {
	v1.UnimplementedCityServer

	uc  *biz.CityUsecase
	log *log.Helper
}

// NewCityService 创建城市服务实例
func NewCityService(uc *biz.CityUsecase, logger log.Logger) *CityService {
	return &CityService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// CreateCity 创建城市
func (s *CityService) CreateCity(ctx context.Context, req *v1.CreateCityRequest) (*v1.CreateCityReply, error) {
	// 参数验证
	if req.Name == "" {
		return nil, v1.ErrorBadRequest("城市名称不能为空")
	}

	// 转换为业务实体
	city := &biz.City{
		Name:   req.Name,
		Status: int8(req.Status),
	}

	// 调用业务逻辑创建城市
	result, err := s.uc.CreateCity(ctx, city)
	if err != nil {
		s.log.Errorf("创建城市失败: %v", err)
		return nil, v1.ErrorInternalError("创建城市失败: %v", err)
	}

	// 构建响应
	return &v1.CreateCityReply{
		Id:        result.ID,
		Name:      result.Name,
		Code:      result.Code,
		Latitude:  result.Latitude,
		Longitude: result.Longitude,
		Status:    int32(result.Status),
	}, nil
}

// UpdateCity 更新城市
func (s *CityService) UpdateCity(ctx context.Context, req *v1.UpdateCityRequest) (*v1.UpdateCityReply, error) {
	// 参数验证
	if req.Id <= 0 {
		return nil, v1.ErrorBadRequest("城市ID无效")
	}
	if req.Name == "" {
		return nil, v1.ErrorBadRequest("城市名称不能为空")
	}

	// 转换为业务实体
	city := &biz.City{
		ID:     req.Id,
		Name:   req.Name,
		Status: int8(req.Status),
	}

	// 调用业务逻辑更新城市
	result, err := s.uc.UpdateCity(ctx, city)
	if err != nil {
		s.log.Errorf("更新城市失败: %v", err)
		return nil, v1.ErrorInternalError("更新城市失败: %v", err)
	}

	// 构建响应
	return &v1.UpdateCityReply{
		Id:        result.ID,
		Name:      result.Name,
		Code:      result.Code,
		Latitude:  result.Latitude,
		Longitude: result.Longitude,
		Status:    int32(result.Status),
	}, nil
}

// GetCity 获取城市详情
func (s *CityService) GetCity(ctx context.Context, req *v1.GetCityRequest) (*v1.GetCityReply, error) {
	// 参数验证
	if req.Id <= 0 {
		return nil, v1.ErrorBadRequest("城市ID无效")
	}

	// 调用业务逻辑获取城市
	city, err := s.uc.GetCityByID(ctx, req.Id)
	if err != nil {
		s.log.Errorf("获取城市失败: %v", err)
		return nil, v1.ErrorInternalError("获取城市失败: %v", err)
	}

	// 构建响应
	return &v1.GetCityReply{
		Id:        city.ID,
		Name:      city.Name,
		Code:      city.Code,
		Latitude:  city.Latitude,
		Longitude: city.Longitude,
		Status:    int32(city.Status),
	}, nil
}

// ListCities 获取城市列表
func (s *CityService) ListCities(ctx context.Context, req *v1.ListCitiesRequest) (*v1.ListCitiesReply, error) {
	// 参数验证和默认值设置
	page := req.Page
	if page <= 0 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 始终查询所有状态的城市
	status := int8(-1) // -1表示查询所有状态

	// 打印请求参数
	s.log.Infof("获取城市列表请求: page=%d, pageSize=%d, status=%d", page, pageSize, status)

	// 调用业务逻辑获取城市列表
	cities, total, err := s.uc.ListCities(ctx, page, pageSize, status)
	if err != nil {
		s.log.Errorf("获取城市列表失败: %v", err)
		return nil, v1.ErrorInternalError("获取城市列表失败: %v", err)
	}

	// 打印结果信息
	s.log.Infof("获取城市列表成功: 总数=%d, 当前页数量=%d", total, len(cities))

	// 构建响应
	items := make([]*v1.CityInfo, 0, len(cities))
	for i, city := range cities {
		if city != nil {
			cityInfo := &v1.CityInfo{
				Id:        city.ID,
				Name:      city.Name,
				Code:      city.Code,
				Latitude:  city.Latitude,
				Longitude: city.Longitude,
				Status:    int32(city.Status),
			}
			items = append(items, cityInfo)
			s.log.Infof("添加城市[%d]到响应: ID=%d, 名称=%s, 编码=%s", i, city.ID, city.Name, city.Code)
		} else {
			s.log.Warnf("城市[%d]为空", i)
		}
	}

	// 即使没有数据，也返回空数组而不是null
	if items == nil {
		items = make([]*v1.CityInfo, 0)
	}

	reply := &v1.ListCitiesReply{
		Total: total,
		Items: items,
	}
	s.log.Infof("返回响应: total=%d, items.length=%d", reply.Total, len(reply.Items))
	
	return reply, nil
}

// UpdateCityStatus 更新城市状态
func (s *CityService) UpdateCityStatus(ctx context.Context, req *v1.UpdateCityStatusRequest) (*v1.UpdateCityStatusReply, error) {
	// 参数验证
	if req.Id <= 0 {
		return nil, v1.ErrorBadRequest("城市ID无效")
	}
	if req.Status != 0 && req.Status != 1 {
		return nil, v1.ErrorBadRequest("状态值无效，只能为0或1")
	}

	// 调用业务逻辑更新城市状态
	err := s.uc.UpdateCityStatus(ctx, req.Id, int8(req.Status))
	if err != nil {
		s.log.Errorf("更新城市状态失败: %v", err)
		return nil, v1.ErrorInternalError("更新城市状态失败: %v", err)
	}

	// 构建响应
	return &v1.UpdateCityStatusReply{
		Success: true,
	}, nil
}

// ListUserCities 获取用户端城市列表
func (s *CityService) ListUserCities(ctx context.Context, req *v1.ListUserCitiesRequest) (*v1.ListUserCitiesReply, error) {
	// 参数验证和默认值设置
	page := req.Page
	if page <= 0 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 打印请求参数
	s.log.Infof("获取用户端城市列表请求: page=%d, pageSize=%d", page, pageSize)

	// 调用业务逻辑获取城市列表
	cities, total, err := s.uc.ListUserCities(ctx, page, pageSize)
	if err != nil {
		s.log.Errorf("获取用户端城市列表失败: %v", err)
		return nil, v1.ErrorInternalError("获取城市列表失败: %v", err)
	}

	// 打印结果信息
	s.log.Infof("获取用户端城市列表成功: 总数=%d, 当前页数量=%d", total, len(cities))

	// 构建响应
	items := make([]*v1.UserCityInfo, 0, len(cities))
	for i, city := range cities {
		if city != nil {
			cityInfo := &v1.UserCityInfo{
				Id:        city.ID,
				Name:      city.Name,
				Code:      city.Code,
				Latitude:  city.Latitude,
				Longitude: city.Longitude,
			}
			items = append(items, cityInfo)
			s.log.Infof("添加城市[%d]到用户端响应: ID=%d, 名称=%s", i, city.ID, city.Name)
		}
	}

	// 即使没有数据，也返回空数组而不是null
	if items == nil {
		items = make([]*v1.UserCityInfo, 0)
	}

	reply := &v1.ListUserCitiesReply{
		Total: total,
		Items: items,
	}
	
	return reply, nil
}

// SearchCities 搜索城市
func (s *CityService) SearchCities(ctx context.Context, req *v1.SearchCitiesRequest) (*v1.SearchCitiesReply, error) {
	// 参数验证和默认值设置
	page := req.Page
	if page <= 0 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 打印请求参数
	s.log.Infof("搜索城市请求: keyword=%s, page=%d, pageSize=%d", req.Keyword, page, pageSize)

	// 调用业务逻辑搜索城市
	cities, total, err := s.uc.SearchCities(ctx, req.Keyword, page, pageSize)
	if err != nil {
		s.log.Errorf("搜索城市失败: %v", err)
		return nil, v1.ErrorInternalError("搜索城市失败: %v", err)
	}

	// 打印结果信息
	s.log.Infof("搜索城市成功: 总数=%d, 当前页数量=%d", total, len(cities))

	// 构建响应
	items := make([]*v1.UserCityInfo, 0, len(cities))
	for i, city := range cities {
		if city != nil {
			cityInfo := &v1.UserCityInfo{
				Id:        city.ID,
				Name:      city.Name,
				Code:      city.Code,
				Latitude:  city.Latitude,
				Longitude: city.Longitude,
			}
			items = append(items, cityInfo)
			s.log.Infof("添加城市[%d]到搜索结果: ID=%d, 名称=%s", i, city.ID, city.Name)
		}
	}

	// 即使没有数据，也返回空数组而不是null
	if items == nil {
		items = make([]*v1.UserCityInfo, 0)
	}

	reply := &v1.SearchCitiesReply{
		Total: total,
		Items: items,
	}
	
	return reply, nil
}

// GetUserCity 获取用户端城市详情
func (s *CityService) GetUserCity(ctx context.Context, req *v1.GetUserCityRequest) (*v1.GetUserCityReply, error) {
	// 参数验证
	if req.Id <= 0 {
		return nil, v1.ErrorBadRequest("城市ID无效")
	}

	// 调用业务逻辑获取城市
	city, err := s.uc.GetUserCity(ctx, req.Id)
	if err != nil {
		s.log.Errorf("获取用户端城市详情失败: %v", err)
		return nil, v1.ErrorInternalError("获取城市详情失败: %v", err)
	}

	// 构建响应
	return &v1.GetUserCityReply{
		Id:        city.ID,
		Name:      city.Name,
		Code:      city.Code,
		Latitude:  city.Latitude,
		Longitude: city.Longitude,
	}, nil
}

// GetCityByCode 根据城市编码获取城市
func (s *CityService) GetCityByCode(ctx context.Context, req *v1.GetCityByCodeRequest) (*v1.GetUserCityReply, error) {
	// 参数验证
	if req.Code == "" {
		return nil, v1.ErrorBadRequest("城市编码不能为空")
	}

	// 调用业务逻辑获取城市
	city, err := s.uc.GetUserCityByCode(ctx, req.Code)
	if err != nil {
		s.log.Errorf("根据编码获取城市失败: %v", err)
		return nil, v1.ErrorInternalError("获取城市详情失败: %v", err)
	}

	// 构建响应
	return &v1.GetUserCityReply{
		Id:        city.ID,
		Name:      city.Name,
		Code:      city.Code,
		Latitude:  city.Latitude,
		Longitude: city.Longitude,
	}, nil
}

// GetHotCities 获取热门城市
func (s *CityService) GetHotCities(ctx context.Context, req *v1.GetHotCitiesRequest) (*v1.ListUserCitiesReply, error) {
	// 设置默认限制
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	// 调用业务逻辑获取热门城市
	cities, total, err := s.uc.GetHotCities(ctx, limit)
	if err != nil {
		s.log.Errorf("获取热门城市失败: %v", err)
		return nil, v1.ErrorInternalError("获取热门城市失败: %v", err)
	}

	// 打印结果信息
	s.log.Infof("获取热门城市成功: 总数=%d, 返回数量=%d", total, len(cities))

	// 构建响应
	items := make([]*v1.UserCityInfo, 0, len(cities))
	for i, city := range cities {
		if city != nil {
			cityInfo := &v1.UserCityInfo{
				Id:        city.ID,
				Name:      city.Name,
				Code:      city.Code,
				Latitude:  city.Latitude,
				Longitude: city.Longitude,
			}
			items = append(items, cityInfo)
			s.log.Infof("添加热门城市[%d]到响应: ID=%d, 名称=%s", i, city.ID, city.Name)
		}
	}

	// 即使没有数据，也返回空数组而不是null
	if items == nil {
		items = make([]*v1.UserCityInfo, 0)
	}

	reply := &v1.ListUserCitiesReply{
		Total: total,
		Items: items,
	}
	
	return reply, nil
}

package biz

import (
	"context"
	"errors"
	"ito-deposit/internal/pkg/baidumap"

	"github.com/go-kratos/kratos/v2/log"
)

// City 城市业务实体
type City struct {
	ID        int32   // 城市ID
	Name      string  // 城市名称
	Code      string  // 城市编码
	Latitude  float64 // 纬度
	Longitude float64 // 经度
	Status    int8    // 状态(1:启用,0:禁用)
}

// CityRepo 城市数据仓库接口
type CityRepo interface {
	// CreateCity 创建城市
	CreateCity(ctx context.Context, city *City) (*City, error)
	// UpdateCity 更新城市
	UpdateCity(ctx context.Context, city *City) (*City, error)
	// GetCityByID 根据ID获取城市
	GetCityByID(ctx context.Context, id int32) (*City, error)
	// GetCityByCode 根据城市编码获取城市
	GetCityByCode(ctx context.Context, code string) (*City, error)
	// GetCityByName 根据城市名称获取城市
	GetCityByName(ctx context.Context, name string) (*City, error)
	// ListCities 获取城市列表
	ListCities(ctx context.Context, page, pageSize int32, status int8) ([]*City, int64, error)
	// UpdateCityStatus 更新城市状态
	UpdateCityStatus(ctx context.Context, id int32, status int8) error
	// SearchCities 搜索城市
	SearchCities(ctx context.Context, keyword string, page, pageSize int32) ([]*City, int64, error)
	// GetHotCities 获取热门城市
	GetHotCities(ctx context.Context, limit int32) ([]*City, int64, error)
}

// CityUsecase 城市用例
type CityUsecase struct {
	repo     CityRepo
	baiduMap *baidumap.BaiduMapClient
	log      *log.Helper
}

// NewCityUsecase 创建城市用例实例
func NewCityUsecase(repo CityRepo, logger log.Logger) *CityUsecase {
	// 创建百度地图客户端，使用提供的AK
	baiduMapClient := baidumap.NewBaiduMapClient("7pzoTHchDdMRK7jmpCr1sugjv3hfoxz5")

	return &CityUsecase{
		repo:     repo,
		baiduMap: baiduMapClient,
		log:      log.NewHelper(logger),
	}
}

// CreateCity 创建城市
func (uc *CityUsecase) CreateCity(ctx context.Context, city *City) (*City, error) {
	// 调用百度地图API获取城市地理编码信息
	geocodeResp, err := uc.baiduMap.Geocode(city.Name)
	if err != nil {
		uc.log.Errorf("获取城市地理编码失败: %v", err)
		return nil, err
	}

	// 设置城市编码
	// 优先使用6位数字行政区划编码(adcode)，其次使用字母缩写(citycode)
	if geocodeResp.AdCode != "" && len(geocodeResp.AdCode) == 6 {
		// 使用6位数字行政区划编码
		city.Code = geocodeResp.AdCode
		uc.log.Infof("使用6位数字行政区划编码(adcode): %s", city.Code)
	} else if geocodeResp.CityCode != "" {
		// 使用字母缩写城市编码
		city.Code = geocodeResp.CityCode
		uc.log.Infof("使用字母缩写城市编码(citycode): %s", city.Code)
	} else {
		// 如果没有获取到编码，使用城市名称作为编码
		// 这确保了即使API没有返回编码，我们也能创建城市记录
		city.Code = city.Name
		uc.log.Infof("未获取到编码，使用城市名称作为编码: %s", city.Code)
	}

	// 设置经纬度
	city.Latitude = geocodeResp.Latitude
	city.Longitude = geocodeResp.Longitude

	// 打印调试信息
	uc.log.Infof("百度地图API返回 - 城市: %s, 编码: %s, 经度: %f, 纬度: %f", city.Name, city.Code, city.Longitude, city.Latitude)

	// 验证坐标是否合理
	if city.Longitude < 70 || city.Longitude > 140 || city.Latitude < 10 || city.Latitude > 60 {
		uc.log.Errorf("警告：城市 %s 的坐标可能不正确 - 经度: %f, 纬度: %f", city.Name, city.Longitude, city.Latitude)
	}

	// 调用数据仓库创建城市
	return uc.repo.CreateCity(ctx, city)
}

// UpdateCity 更新城市
func (uc *CityUsecase) UpdateCity(ctx context.Context, city *City) (*City, error) {
	// 如果城市名称发生变化，需要重新获取地理编码信息
	existCity, err := uc.repo.GetCityByID(ctx, city.ID)
	if err != nil {
		return nil, err
	}

	if existCity.Name != city.Name {
		// 调用百度地图API获取城市地理编码信息
		geocodeResp, err := uc.baiduMap.Geocode(city.Name)
		if err != nil {
			uc.log.Errorf("获取城市地理编码失败: %v", err)
			return nil, err
		}

		// 设置城市编码
		// 优先使用6位数字行政区划编码(adcode)，其次使用字母缩写(citycode)
		if geocodeResp.AdCode != "" && len(geocodeResp.AdCode) == 6 {
			// 使用6位数字行政区划编码
			city.Code = geocodeResp.AdCode
			uc.log.Infof("使用6位数字行政区划编码(adcode): %s", city.Code)
		} else if geocodeResp.CityCode != "" {
			// 使用字母缩写城市编码
			city.Code = geocodeResp.CityCode
			uc.log.Infof("使用字母缩写城市编码(citycode): %s", city.Code)
		} else {
			// 如果没有获取到编码，使用城市名称作为编码
			// 这确保了即使API没有返回编码，我们也能更新城市记录
			city.Code = city.Name
			uc.log.Infof("未获取到编码，使用城市名称作为编码: %s", city.Code)
		}

		// 设置经纬度
		city.Latitude = geocodeResp.Latitude
		city.Longitude = geocodeResp.Longitude

		// 打印调试信息
		uc.log.Infof("更新城市编码: %s, 经度: %f, 纬度: %f", city.Code, city.Longitude, city.Latitude)
	}

	// 调用数据仓库更新城市
	return uc.repo.UpdateCity(ctx, city)
}

// GetCityByID 根据ID获取城市
func (uc *CityUsecase) GetCityByID(ctx context.Context, id int32) (*City, error) {
	return uc.repo.GetCityByID(ctx, id)
}

// ListCities 获取城市列表
func (uc *CityUsecase) ListCities(ctx context.Context, page, pageSize int32, status int8) ([]*City, int64, error) {
	return uc.repo.ListCities(ctx, page, pageSize, status)
}

// UpdateCityStatus 更新城市状态
func (uc *CityUsecase) UpdateCityStatus(ctx context.Context, id int32, status int8) error {
	return uc.repo.UpdateCityStatus(ctx, id, status)
}

// ListUserCities 获取用户端城市列表（只返回启用状态的城市）
func (uc *CityUsecase) ListUserCities(ctx context.Context, page, pageSize int32) ([]*City, int64, error) {
	// 只查询状态为启用(1)的城市
	return uc.repo.ListCities(ctx, page, pageSize, 1)
}

// SearchCities 搜索城市（只搜索启用状态的城市）
func (uc *CityUsecase) SearchCities(ctx context.Context, keyword string, page, pageSize int32) ([]*City, int64, error) {
	return uc.repo.SearchCities(ctx, keyword, page, pageSize)
}

// GetUserCity 获取用户端城市详情（只返回启用状态的城市）
func (uc *CityUsecase) GetUserCity(ctx context.Context, id int32) (*City, error) {
	city, err := uc.repo.GetCityByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 检查城市状态是否为启用
	if city.Status != 1 {
		return nil, errors.New("城市不存在或已禁用")
	}

	return city, nil
}

// GetUserCityByCode 根据城市编码获取城市（只返回启用状态的城市）
func (uc *CityUsecase) GetUserCityByCode(ctx context.Context, code string) (*City, error) {
	city, err := uc.repo.GetCityByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	// 检查城市状态是否为启用
	if city.Status != 1 {
		return nil, errors.New("城市不存在或已禁用")
	}

	return city, nil
}

// GetUserCityByName 根据城市名称获取城市（只返回启用状态的城市）
func (uc *CityUsecase) GetUserCityByName(ctx context.Context, name string) (*City, error) {
	city, err := uc.repo.GetCityByName(ctx, name)
	if err != nil {
		return nil, err
	}

	// 检查城市状态是否为启用
	if city.Status != 1 {
		return nil, errors.New("城市不存在或已禁用")
	}

	return city, nil
}

// GetHotCities 获取热门城市（只返回启用状态的城市）
func (uc *CityUsecase) GetHotCities(ctx context.Context, limit int32) ([]*City, int64, error) {
	return uc.repo.GetHotCities(ctx, limit)
}

package data

import (
	"context"
	"errors"
	"fmt"
	"ito-deposit/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

// cityRepo 城市数据仓库实现
type cityRepo struct {
	data *Data
	log  *log.Helper
}

// NewCityRepo 创建城市数据仓库实例
func NewCityRepo(data *Data, logger log.Logger) biz.CityRepo {
	return &cityRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// 将数据库模型转换为业务实体
func (r *cityRepo) convertToBizCity(city *City) *biz.City {
	if city == nil {
		return nil
	}
	return &biz.City{
		ID:        city.ID,
		Name:      city.Name,
		Code:      city.Code,
		Latitude:  city.Latitude,
		Longitude: city.Longitude,
		Status:    city.Status,
	}
}

// 将业务实体转换为数据库模型
func (r *cityRepo) convertToDataCity(city *biz.City) *City {
	return &City{
		ID:        city.ID,
		Name:      city.Name,
		Code:      city.Code,
		Latitude:  city.Latitude,
		Longitude: city.Longitude,
		Status:    city.Status,
	}
}

// CreateCity 创建城市
func (r *cityRepo) CreateCity(ctx context.Context, city *biz.City) (*biz.City, error) {
	dataCity := r.convertToDataCity(city)
	var existCity City
	err := r.data.DB.Where("name = ?", dataCity.Name).First(&existCity).Error
	if err == nil {
		return nil, errors.New("城市名称已存在") // 优先提示名称冲突
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	fmt.Println(dataCity.Code, 1111111)
	// 只有编码非空时才检查唯一性
	if dataCity.Code != "" {
		err := r.data.DB.Where("code = ?", dataCity.Code).First(&existCity).Error
		if err == nil {
			return nil, errors.New("城市编码已存在")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Errorf("查询城市编码失败: %v", err)
			return nil, err
		}
	}

	// 创建记录
	if err := r.data.DB.Create(dataCity).Error; err != nil {
		r.log.Errorf("创建城市失败: %v", err)
		return nil, err
	}
	return r.convertToBizCity(dataCity), nil
}

// UpdateCity 更新城市
func (r *cityRepo) UpdateCity(ctx context.Context, city *biz.City) (*biz.City, error) {
	// 将业务实体转换为数据库模型
	dataCity := r.convertToDataCity(city)

	// 检查城市是否存在
	var existCity City
	if err := r.data.DB.First(&existCity, dataCity.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("城市不存在")
		}
		r.log.Errorf("查询城市失败: %v", err)
		return nil, err
	}

	// 更新城市信息
	if err := r.data.DB.Model(&City{}).Where("id = ?", dataCity.ID).Updates(map[string]interface{}{
		"name":      dataCity.Name,
		"code":      dataCity.Code,
		"latitude":  dataCity.Latitude,
		"longitude": dataCity.Longitude,
		"status":    dataCity.Status,
	}).Error; err != nil {
		r.log.Errorf("更新城市失败: %v", err)
		return nil, err
	}

	// 重新获取更新后的城市信息
	if err := r.data.DB.First(&dataCity, dataCity.ID).Error; err != nil {
		r.log.Errorf("获取更新后的城市信息失败: %v", err)
		return nil, err
	}

	// 将数据库模型转换回业务实体
	return r.convertToBizCity(dataCity), nil
}

// GetCityByID 根据ID获取城市
func (r *cityRepo) GetCityByID(ctx context.Context, id int32) (*biz.City, error) {
	var city City
	if err := r.data.DB.First(&city, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("城市不存在")
		}
		r.log.Errorf("获取城市失败: %v", err)
		return nil, err
	}

	// 将数据库模型转换为业务实体
	return r.convertToBizCity(&city), nil
}

// GetCityByCode 根据城市编码获取城市
func (r *cityRepo) GetCityByCode(ctx context.Context, code string) (*biz.City, error) {
	var city City
	if err := r.data.DB.Where("code = ?", code).First(&city).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("城市不存在")
		}
		r.log.Errorf("获取城市失败: %v", err)
		return nil, err
	}

	// 将数据库模型转换为业务实体
	return r.convertToBizCity(&city), nil
}

// GetCityByName 根据城市名称获取城市
func (r *cityRepo) GetCityByName(ctx context.Context, name string) (*biz.City, error) {
	var city City

	// 首先尝试精确匹配
	err := r.data.DB.Where("name = ?", name).First(&city).Error
	if err == nil {
		return r.convertToBizCity(&city), nil
	}

	// 如果精确匹配失败且是记录不存在错误，尝试模糊匹配
	if errors.Is(err, gorm.ErrRecordNotFound) {
		r.log.Infof("精确匹配失败，尝试模糊匹配城市: %s", name)

		// 尝试模糊匹配
		err = r.data.DB.Where("name LIKE ? OR name LIKE ?", "%"+name+"%", name+"%").First(&city).Error
		if err == nil {
			r.log.Infof("模糊匹配成功: %s -> %s", name, city.Name)
			return r.convertToBizCity(&city), nil
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("城市不存在")
		}
	}

	r.log.Errorf("获取城市失败: %v", err)
	return nil, err
}

// ListCities 获取城市列表
func (r *cityRepo) ListCities(ctx context.Context, page, pageSize int32, status int8) ([]*biz.City, int64, error) {
	var cities []City
	var total int64

	// 打印SQL语句
	db := r.data.DB.Debug()

	// 构建查询条件
	query := db.Model(&City{})
	if status != -1 {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.log.Errorf("获取城市总数失败: %v", err)
		return nil, 0, err
	}

	// 打印调试信息
	r.log.Infof("城市总数: %d", total)

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&cities).Error; err != nil {
		r.log.Errorf("获取城市列表失败: %v", err)
		return nil, 0, err
	}

	// 打印调试信息
	r.log.Infof("查询到城市数量: %d", len(cities))
	for i, city := range cities {
		r.log.Infof("数据库城市[%d]: ID=%d, 名称=%s, 编码=%s, 状态=%d", i, city.ID, city.Name, city.Code, city.Status)
	}

	// 将数据库模型转换为业务实体
	bizCities := make([]*biz.City, 0, len(cities))
	for i := range cities {
		bizCity := r.convertToBizCity(&cities[i])
		bizCities = append(bizCities, bizCity)
		r.log.Infof("业务实体城市: ID=%d, 名称=%s, 编码=%s, 状态=%d", bizCity.ID, bizCity.Name, bizCity.Code, bizCity.Status)
	}

	return bizCities, total, nil
}

// UpdateCityStatus 更新城市状态
func (r *cityRepo) UpdateCityStatus(ctx context.Context, id int32, status int8) error {
	// 检查城市是否存在
	var city City
	if err := r.data.DB.First(&city, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("城市不存在")
		}
		r.log.Errorf("查询城市失败: %v", err)
		return err
	}

	// 更新状态
	if err := r.data.DB.Model(&City{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		r.log.Errorf("更新城市状态失败: %v", err)
		return err
	}
	return nil
}

// SearchCities 搜索城市
func (r *cityRepo) SearchCities(ctx context.Context, keyword string, page, pageSize int32) ([]*biz.City, int64, error) {
	var cities []City
	var total int64

	// 打印SQL语句
	db := r.data.DB.Debug()

	// 构建查询条件
	query := db.Model(&City{}).Where("status = ?", 1) // 只查询启用状态的城市

	// 如果有关键词，添加模糊查询条件
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.log.Errorf("获取城市总数失败: %v", err)
		return nil, 0, err
	}

	// 打印调试信息
	r.log.Infof("搜索城市总数: %d, 关键词: %s", total, keyword)

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&cities).Error; err != nil {
		r.log.Errorf("搜索城市失败: %v", err)
		return nil, 0, err
	}

	// 打印调试信息
	r.log.Infof("搜索到城市数量: %d", len(cities))
	for i, city := range cities {
		r.log.Infof("搜索结果城市[%d]: ID=%d, 名称=%s, 编码=%s", i, city.ID, city.Name, city.Code)
	}

	// 将数据库模型转换为业务实体
	bizCities := make([]*biz.City, 0, len(cities))
	for i := range cities {
		bizCity := r.convertToBizCity(&cities[i])
		bizCities = append(bizCities, bizCity)
	}

	return bizCities, total, nil
}

// GetHotCities 获取热门城市
func (r *cityRepo) GetHotCities(ctx context.Context, limit int32) ([]*biz.City, int64, error) {
	var cities []City
	var total int64

	// 打印SQL语句
	db := r.data.DB.Debug()

	// 构建查询条件：只查询启用状态的城市
	query := db.Model(&City{}).Where("status = ?", 1)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.log.Errorf("获取热门城市总数失败: %v", err)
		return nil, 0, err
	}

	// 打印调试信息
	r.log.Infof("热门城市总数: %d", total)

	// 设置默认限制
	if limit <= 0 {
		limit = 20
	}

	// 查询热门城市（这里简单地按ID排序，实际应用中可能需要根据访问量或其他指标排序）
	if err := query.Order("id ASC").Limit(int(limit)).Find(&cities).Error; err != nil {
		r.log.Errorf("获取热门城市失败: %v", err)
		return nil, 0, err
	}

	// 打印调试信息
	r.log.Infof("获取到热门城市数量: %d", len(cities))
	for i, city := range cities {
		r.log.Infof("热门城市[%d]: ID=%d, 名称=%s, 编码=%s", i, city.ID, city.Name, city.Code)
	}

	// 将数据库模型转换为业务实体
	bizCities := make([]*biz.City, 0, len(cities))
	for i := range cities {
		bizCity := r.convertToBizCity(&cities[i])
		bizCities = append(bizCities, bizCity)
	}

	return bizCities, total, nil
}

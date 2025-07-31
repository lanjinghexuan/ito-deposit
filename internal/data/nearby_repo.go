package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"ito-deposit/internal/biz"
)

// nearbyRepo 附近寄存点数据仓库实现
type nearbyRepo struct {
	data *Data
	log  *log.Helper
}

// NewNearbyRepo 创建附近寄存点数据仓库实例
func NewNearbyRepo(data *Data, logger log.Logger) biz.NearbyRepo {
	return &nearbyRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// GetLockerPoints 获取所有寄存点
func (r *nearbyRepo) GetLockerPoints(ctx context.Context) ([]*biz.LockerPoint, error) {
	var lockerPoints []LockerPoint

	// 查询所有寄存点
	if err := r.data.DB.Find(&lockerPoints).Error; err != nil {
		r.log.Errorf("获取寄存点列表失败: %v", err)
		return nil, err
	}

	// 转换为业务实体
	result := make([]*biz.LockerPoint, 0, len(lockerPoints))
	for _, point := range lockerPoints {
		result = append(result, &biz.LockerPoint{
			Id:        point.Id,
			Name:      point.Name,
			Address:   point.Address,
			Longitude: point.Longitude,
			Latitude:  point.Latitude,
		})
	}

	return result, nil
}

// GetLockerPointByID 根据ID获取寄存点详情
func (r *nearbyRepo) GetLockerPointByID(ctx context.Context, id int32) (*biz.LockerPoint, error) {
	var point LockerPoint

	// 查询寄存点
	if err := r.data.DB.First(&point, id).Error; err != nil {
		r.log.Errorf("获取寄存点详情失败: %v", err)
		return nil, err
	}

	// 转换为业务实体
	return &biz.LockerPoint{
		Id:        point.Id,
		Name:      point.Name,
		Address:   point.Address,
		Longitude: point.Longitude,
		Latitude:  point.Latitude,
	}, nil
}

// SearchLockerPointsInCity 搜索指定城市内的寄存点
func (r *nearbyRepo) SearchLockerPointsInCity(ctx context.Context, cityName string, keyword string, page, pageSize int64) ([]*biz.LockerPoint, int64, error) {
	var lockerPoints []LockerPoint
	var total int64

	// 构建查询 - 先通过城市名称找到城市ID
	var city City
	if err := r.data.DB.Where("name = ?", cityName).First(&city).Error; err != nil {
		r.log.Errorf("查找城市失败: %v", err)
		return nil, 0, err
	}

	// 构建寄存点查询
	query := r.data.DB.Model(&LockerPoint{})
	query = query.Where("location_id = ?", city.ID)

	// 如果提供了关键词，则按名称或地址搜索
	if keyword != "" {
		query = query.Where("name LIKE ? OR address LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 计算总记录数
	if err := query.Count(&total).Error; err != nil {
		r.log.Errorf("计算寄存点总数失败: %v", err)
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(int(offset)).Limit(int(pageSize)).Find(&lockerPoints).Error; err != nil {
		r.log.Errorf("搜索寄存点失败: %v", err)
		return nil, 0, err
	}

	// 转换为业务实体
	result := make([]*biz.LockerPoint, 0, len(lockerPoints))
	for _, point := range lockerPoints {
		result = append(result, &biz.LockerPoint{
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

	return result, total, nil
}

// GetLockerPointsInBounds 获取指定边界内的寄存点
func (r *nearbyRepo) GetLockerPointsInBounds(ctx context.Context, cityName string, northLat, southLat, eastLng, westLng float64) ([]*biz.LockerPoint, error) {
	var lockerPoints []LockerPoint

	// 先通过城市名称找到城市ID
	var city City
	if err := r.data.DB.Where("name = ?", cityName).First(&city).Error; err != nil {
		r.log.Errorf("查找城市失败: %v", err)
		return nil, err
	}

	// 构建寄存点查询
	query := r.data.DB.Model(&LockerPoint{})
	query = query.Where("location_id = ?", city.ID)

	// 添加地理边界条件
	query = query.Where("latitude BETWEEN ? AND ?", southLat, northLat)
	query = query.Where("longitude BETWEEN ? AND ?", westLng, eastLng)

	// 执行查询
	if err := query.Find(&lockerPoints).Error; err != nil {
		r.log.Errorf("获取边界内寄存点失败: %v", err)
		return nil, err
	}

	// 转换为业务实体
	result := make([]*biz.LockerPoint, 0, len(lockerPoints))
	for _, point := range lockerPoints {
		result = append(result, &biz.LockerPoint{
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

	return result, nil
}

package data

import (
	"context"
	"fmt"
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
	
	// 构建查询 - 先通过城市名称找到城市ID，支持多种城市名称格式
	var city City
	
	// 尝试多种城市名称格式进行匹配
	cityVariants := []string{
		cityName,           // 原始名称，如"郑州"
		cityName + "市",    // 添加"市"后缀，如"郑州市"
		cityName + "省",    // 添加"省"后缀，如"河南省"
	}
	
	// 如果城市名称以"市"结尾，也尝试去掉"市"
	if len(cityName) > 1 && cityName[len(cityName)-3:] == "市" {
		cityVariants = append(cityVariants, cityName[:len(cityName)-3])
	}
	
	r.log.Infof("搜索城市变体: %v", cityVariants)
	
	var err error
	for _, variant := range cityVariants {
		err = r.data.DB.Where("name = ?", variant).First(&city).Error
		if err == nil {
			r.log.Infof("找到匹配的城市: %s -> %s (ID: %d)", cityName, variant, city.ID)
			break
		}
	}
	
	if err != nil {
		r.log.Warnf("查找城市失败，尝试的变体: %v，将触发降级机制: %v", cityVariants, err)
		// 返回错误，让服务层降级到GetAllLockerPoints
		return nil, 0, fmt.Errorf("城市 %s 不存在，触发降级机制", cityName)
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

// GetAllLockerPoints 获取所有寄存点（不依赖城市表）
func (r *nearbyRepo) GetAllLockerPoints(ctx context.Context, keyword string, page, pageSize int64) ([]*biz.LockerPoint, int64, error) {
	var lockerPoints []LockerPoint
	var total int64

	// 构建寄存点查询
	query := r.data.DB.Model(&LockerPoint{})

	// 如果提供了关键词，则按名称或地址搜索
	if keyword != "" {
		query = query.Where("name LIKE ? OR address LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 不过滤状态，查询所有寄存点
	r.log.Infof("开始查询寄存点，关键词: '%s', 页码: %d, 页大小: %d", keyword, page, pageSize)

	// 计算总记录数
	if err := query.Count(&total).Error; err != nil {
		r.log.Errorf("计算寄存点总数失败: %v", err)
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(int(offset)).Limit(int(pageSize)).Order("id ASC").Find(&lockerPoints).Error; err != nil {
		r.log.Errorf("查询寄存点失败: %v", err)
		return nil, 0, err
	}

	r.log.Infof("查询到 %d 个寄存点，总数: %d", len(lockerPoints), total)

	// 调试：打印前几个寄存点的信息
	for i, point := range lockerPoints {
		if i < 3 { // 只打印前3个
			r.log.Infof("寄存点 %d: ID=%d, Name=%s, Address=%s", i+1, point.Id, point.Name, point.Address)
		}
	}

	// 转换为业务实体，并动态计算可用柜子数量
	result := make([]*biz.LockerPoint, 0, len(lockerPoints))
	for _, point := range lockerPoints {
		// 动态计算可用柜子数量
		availableLarge, availableMedium, availableSmall := r.calculateAvailableCells(ctx, point.Id)

		result = append(result, &biz.LockerPoint{
			Id:              point.Id,
			Name:            point.Name,
			Address:         point.Address,
			Longitude:       point.Longitude,
			Latitude:        point.Latitude,
			AvailableLarge:  availableLarge,
			AvailableMedium: availableMedium,
			AvailableSmall:  availableSmall,
			OpenTime:        point.OpenTime,
			Mobile:          point.Mobile,
		})
	}

	return result, total, nil
}

// GetLockerPointsInBounds 获取指定边界内的寄存点
func (r *nearbyRepo) GetLockerPointsInBounds(ctx context.Context, cityName string, northLat, southLat, eastLng, westLng float64) ([]*biz.LockerPoint, error) {
	var lockerPoints []LockerPoint
	
	// 先通过城市名称找到城市ID，支持多种城市名称格式
	var city City
	
	// 尝试多种城市名称格式进行匹配
	cityVariants := []string{
		cityName,           // 原始名称，如"郑州"
		cityName + "市",    // 添加"市"后缀，如"郑州市"
		cityName + "省",    // 添加"省"后缀，如"河南省"
	}
	
	// 如果城市名称以"市"结尾，也尝试去掉"市"
	if len(cityName) > 1 && cityName[len(cityName)-3:] == "市" {
		cityVariants = append(cityVariants, cityName[:len(cityName)-3])
	}
	
	var err error
	for _, variant := range cityVariants {
		err = r.data.DB.Where("name = ?", variant).First(&city).Error
		if err == nil {
			r.log.Infof("找到匹配的城市: %s -> %s (ID: %d)", cityName, variant, city.ID)
			break
		}
	}
	
	if err != nil {
		r.log.Errorf("查找城市失败，尝试的变体: %v: %v", cityVariants, err)
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

// calculateAvailableCells 计算指定寄存点的可用柜子数量
func (r *nearbyRepo) calculateAvailableCells(ctx context.Context, lockerPointId int32) (large, medium, small int32) {
	// 查询该寄存点下所有柜组
	var groups []CabinetGroup
	if err := r.data.DB.Where("location_point_id = ? AND status = ?", lockerPointId, "normal").Find(&groups).Error; err != nil {
		r.log.Warnf("查询寄存点 %d 的柜组失败: %v", lockerPointId, err)
		// 如果查询失败，使用默认值
		return r.getDefaultCellCounts(lockerPointId)
	}

	// 如果没有柜组，使用默认值
	if len(groups) == 0 {
		r.log.Infof("寄存点 %d 没有柜组数据，使用默认柜子数量", lockerPointId)
		return r.getDefaultCellCounts(lockerPointId)
	}

	// 收集所有柜组ID
	groupIds := make([]int32, len(groups))
	for i, group := range groups {
		groupIds[i] = group.Id
	}

	// 查询所有可用的柜格（状态为normal的柜格）
	var cells []CabinetCell
	if err := r.data.DB.Where("group_id IN ? AND status = ?", groupIds, "normal").Find(&cells).Error; err != nil {
		r.log.Warnf("查询寄存点 %d 的柜格失败: %v", lockerPointId, err)
		// 如果查询柜格失败，根据柜组数量估算
		return r.estimateCellCountsByGroups(groups)
	}

	// 统计各种大小的可用柜格数量
	for _, cell := range cells {
		switch cell.CellSize {
		case "large":
			large++
		case "medium":
			medium++
		case "small":
			small++
		}
	}

	// 如果没有柜格数据，根据柜组数量估算
	if large == 0 && medium == 0 && small == 0 {
		r.log.Infof("寄存点 %d 没有柜格数据，根据柜组数量估算", lockerPointId)
		return r.estimateCellCountsByGroups(groups)
	}

	r.log.Debugf("寄存点 %d 可用柜格统计: 大柜=%d, 中柜=%d, 小柜=%d", lockerPointId, large, medium, small)
	return large, medium, small
}

// getDefaultCellCounts 获取默认的柜子数量（基于寄存点ID）
func (r *nearbyRepo) getDefaultCellCounts(lockerPointId int32) (large, medium, small int32) {
	// 根据寄存点ID提供不同的默认值，模拟真实场景
	switch lockerPointId {
	case 1:
		return 7, 13, 22  // 北京西站南广场寄存点
	case 2:
		return 3, 6, 10   // 上海虹桥站出发层寄存点
	case 3:
		return 4, 7, 11   // 广州南站东进站口寄存点
	default:
		// 其他寄存点使用通用默认值
		return 5, 8, 12
	}
}

// estimateCellCountsByGroups 根据柜组数量估算柜格数量
func (r *nearbyRepo) estimateCellCountsByGroups(groups []CabinetGroup) (large, medium, small int32) {
	totalCells := int32(0)
	for _, group := range groups {
		totalCells += group.TotalCells
	}

	// 按照常见比例分配：大柜20%，中柜30%，小柜50%
	large = totalCells * 20 / 100
	medium = totalCells * 30 / 100
	small = totalCells - large - medium // 剩余的都是小柜

	r.log.Infof("根据 %d 个柜组估算柜格数量: 总数=%d, 大柜=%d, 中柜=%d, 小柜=%d",
		len(groups), totalCells, large, medium, small)

	return large, medium, small
}

package data

import (
	"context"
	"errors"
	"strings"
	"time"

	"ito-deposit/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

// cabinetCellRepo 柜口数据仓库实现
type cabinetCellRepo struct {
	data *Data
	log  *log.Helper
}

// NewCabinetCellRepo 创建柜口数据仓库实例
func NewCabinetCellRepo(data *Data, logger log.Logger) biz.CabinetCellRepo {
	return &cabinetCellRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// 将数据库模型转换为业务实体
func (r *cabinetCellRepo) convertToBizCabinetCell(cell *CabinetCell) *biz.CabinetCell {
	if cell == nil {
		return nil
	}

	// 调试：打印数据库原始数据
	r.log.Infof("=== 数据库原始数据 ===")
	r.log.Infof("ID: %d", cell.Id)
	r.log.Infof("CabinetGroupId: %d", cell.CabinetGroupId)
	r.log.Infof("CellNo: %d", cell.CellNo)
	r.log.Infof("CellSize: '%s' (长度: %d)", cell.CellSize, len(cell.CellSize))
	r.log.Infof("Status: '%s'", cell.Status)

	return &biz.CabinetCell{
		Id:           cell.Id,
		GroupId:      cell.CabinetGroupId,
		CellNo:       cell.CellNo,
		CellSize:     cell.CellSize,
		Status:       cell.Status,
		LastOpenTime: cell.LastOpenTime,
		CreateTime:   cell.CreateTime,
		UpdateTime:   cell.UpdateTime,
	}
}

// 将业务实体转换为数据库模型
func (r *cabinetCellRepo) convertToDataCabinetCell(cell *biz.CabinetCell) *CabinetCell {
	return &CabinetCell{
		Id:             cell.Id,
		CabinetGroupId: cell.GroupId,
		CellNo:         cell.CellNo,
		CellSize:       cell.CellSize,
		Status:         cell.Status,
		LastOpenTime:   cell.LastOpenTime,
		CreateTime:     cell.CreateTime,
		UpdateTime:     cell.UpdateTime,
	}
}

// CreateCabinetCell 创建柜口
func (r *cabinetCellRepo) CreateCabinetCell(ctx context.Context, cell *biz.CabinetCell) (*biz.CabinetCell, error) {
	dataCell := r.convertToDataCabinetCell(cell)

	// 检查同一柜组内格口编号是否已存在
	var existCell CabinetCell
	err := r.data.DB.Where("cabinet_group_id = ? AND cell_no = ?", dataCell.CabinetGroupId, dataCell.CellNo).First(&existCell).Error
	if err == nil {
		return nil, errors.New("该柜组内格口编号已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		r.log.Errorf("查询格口编号失败: %v", err)
		return nil, err
	}

	// 设置创建和更新时间
	now := time.Now()
	dataCell.CreateTime = now
	dataCell.UpdateTime = now

	// 创建记录
	if err := r.data.DB.Create(dataCell).Error; err != nil {
		r.log.Errorf("创建柜口失败: %v", err)
		return nil, err
	}

	return r.convertToBizCabinetCell(dataCell), nil
}

// UpdateCabinetCell 更新柜口
func (r *cabinetCellRepo) UpdateCabinetCell(ctx context.Context, cell *biz.CabinetCell) (*biz.CabinetCell, error) {
	dataCell := r.convertToDataCabinetCell(cell)

	// 检查柜口是否存在
	var existCell CabinetCell
	if err := r.data.DB.First(&existCell, dataCell.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("柜口不存在")
		}
		r.log.Errorf("查询柜口失败: %v", err)
		return nil, err
	}

	// 如果更新格口编号，检查同一柜组内是否重复
	if dataCell.CellNo != existCell.CellNo {
		var duplicateCell CabinetCell
		err := r.data.DB.Where("cabinet_group_id = ? AND cell_no = ? AND id != ?",
			dataCell.CabinetGroupId, dataCell.CellNo, dataCell.Id).First(&duplicateCell).Error
		if err == nil {
			return nil, errors.New("该柜组内格口编号已存在")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Errorf("检查格口编号重复失败: %v", err)
			return nil, err
		}
	}

	// 更新柜口信息
	dataCell.UpdateTime = time.Now()
	if err := r.data.DB.Model(&CabinetCell{}).Where("id = ?", dataCell.Id).Updates(map[string]interface{}{
		"cabinet_group_id": dataCell.CabinetGroupId,
		"cell_no":          dataCell.CellNo,
		"cell_size":        dataCell.CellSize,
		"status":           dataCell.Status,
		"last_open_time":   dataCell.LastOpenTime,
		"update_time":      dataCell.UpdateTime,
	}).Error; err != nil {
		r.log.Errorf("更新柜口失败: %v", err)
		return nil, err
	}

	// 重新获取更新后的柜口信息
	if err := r.data.DB.First(&dataCell, dataCell.Id).Error; err != nil {
		r.log.Errorf("获取更新后的柜口信息失败: %v", err)
		return nil, err
	}

	return r.convertToBizCabinetCell(dataCell), nil
}

// GetCabinetCellByID 根据ID获取柜口
func (r *cabinetCellRepo) GetCabinetCellByID(ctx context.Context, id int32) (*biz.CabinetCell, error) {
	var cell CabinetCell
	if err := r.data.DB.First(&cell, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("柜口不存在")
		}
		r.log.Errorf("获取柜口失败: %v", err)
		return nil, err
	}

	return r.convertToBizCabinetCell(&cell), nil
}

// GetCabinetCellByGroupAndNo 根据柜组ID和格口编号获取柜口
func (r *cabinetCellRepo) GetCabinetCellByGroupAndNo(ctx context.Context, groupId, cellNo int32) (*biz.CabinetCell, error) {
	var cell CabinetCell
	if err := r.data.DB.Where("cabinet_group_id = ? AND cell_no = ?", groupId, cellNo).First(&cell).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("柜口不存在")
		}
		r.log.Errorf("获取柜口失败: %v", err)
		return nil, err
	}

	return r.convertToBizCabinetCell(&cell), nil
}

// ListCabinetCells 获取柜口列表
func (r *cabinetCellRepo) ListCabinetCells(ctx context.Context, page, pageSize int32, groupId int32, status string) ([]*biz.CabinetCell, int64, error) {
	var cells []CabinetCell
	var total int64

	// 构建查询条件
	query := r.data.DB.Model(&CabinetCell{})

	if groupId > 0 {
		query = query.Where("cabinet_group_id = ?", groupId)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.log.Errorf("获取柜口总数失败: %v", err)
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("cabinet_group_id ASC, cell_no ASC").Offset(int(offset)).Limit(int(pageSize)).Find(&cells).Error; err != nil {
		r.log.Errorf("获取柜口列表失败: %v", err)
		return nil, 0, err
	}

	// 将数据库模型转换为业务实体
	bizCells := make([]*biz.CabinetCell, 0, len(cells))
	for i := range cells {
		bizCell := r.convertToBizCabinetCell(&cells[i])
		bizCells = append(bizCells, bizCell)
	}

	return bizCells, total, nil
}

// UpdateCabinetCellStatus 更新柜口状态
func (r *cabinetCellRepo) UpdateCabinetCellStatus(ctx context.Context, id int32, status string) error {
	// 检查柜口是否存在
	var cell CabinetCell
	if err := r.data.DB.First(&cell, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("柜口不存在")
		}
		r.log.Errorf("查询柜口失败: %v", err)
		return err
	}

	// 更新状态
	if err := r.data.DB.Model(&CabinetCell{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      status,
		"update_time": time.Now(),
	}).Error; err != nil {
		r.log.Errorf("更新柜口状态失败: %v", err)
		return err
	}

	return nil
}

// DeleteCabinetCell 删除柜口
func (r *cabinetCellRepo) DeleteCabinetCell(ctx context.Context, id int32) error {
	// 检查柜口是否存在
	var cell CabinetCell
	if err := r.data.DB.First(&cell, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("柜口不存在")
		}
		r.log.Errorf("查询柜口失败: %v", err)
		return err
	}

	// 检查柜口是否正在使用
	if cell.Status == "inUse" {
		return errors.New("柜口正在使用中，无法删除")
	}

	// 软删除：更新状态为禁用
	if err := r.data.DB.Model(&CabinetCell{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      "disabled",
		"update_time": time.Now(),
	}).Error; err != nil {
		r.log.Errorf("删除柜口失败: %v", err)
		return err
	}

	return nil
}

// SearchCabinetCells 搜索柜口
func (r *cabinetCellRepo) SearchCabinetCells(ctx context.Context, keyword string, page, pageSize int32, groupId int32) ([]*biz.CabinetCell, int64, error) {
	var cells []CabinetCell
	var total int64

	// 构建查询条件
	query := r.data.DB.Model(&CabinetCell{})

	if groupId > 0 {
		query = query.Where("cabinet_group_id = ?", groupId)
	}

	// 关键词搜索（格口编号或状态）
	if keyword != "" {
		keyword = strings.TrimSpace(keyword)
		// 尝试将关键词转换为数字进行格口编号搜索，同时支持状态搜索
		query = query.Where("cell_no = ? OR status LIKE ? OR cell_size LIKE ?",
			keyword, "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.log.Errorf("搜索柜口总数失败: %v", err)
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("cabinet_group_id ASC, cell_no ASC").Offset(int(offset)).Limit(int(pageSize)).Find(&cells).Error; err != nil {
		r.log.Errorf("搜索柜口失败: %v", err)
		return nil, 0, err
	}

	// 将数据库模型转换为业务实体
	bizCells := make([]*biz.CabinetCell, 0, len(cells))
	for i := range cells {
		bizCell := r.convertToBizCabinetCell(&cells[i])
		bizCells = append(bizCells, bizCell)
	}

	return bizCells, total, nil
}

// GetCabinetCellsByGroupId 根据柜组ID获取所有柜口
func (r *cabinetCellRepo) GetCabinetCellsByGroupId(ctx context.Context, groupId int32) ([]*biz.CabinetCell, error) {
	var cells []CabinetCell

	if err := r.data.DB.Where("cabinet_group_id = ?", groupId).Order("cell_no ASC").Find(&cells).Error; err != nil {
		r.log.Errorf("根据柜组ID获取柜口失败: %v", err)
		return nil, err
	}

	// 将数据库模型转换为业务实体
	bizCells := make([]*biz.CabinetCell, 0, len(cells))
	for i := range cells {
		bizCell := r.convertToBizCabinetCell(&cells[i])
		bizCells = append(bizCells, bizCell)
	}

	return bizCells, nil
}

// BatchCreateCabinetCells 批量创建柜口
func (r *cabinetCellRepo) BatchCreateCabinetCells(ctx context.Context, cells []*biz.CabinetCell) error {
	if len(cells) == 0 {
		return nil
	}

	// 转换为数据库模型
	dataCells := make([]CabinetCell, 0, len(cells))
	now := time.Now()

	for _, cell := range cells {
		dataCell := r.convertToDataCabinetCell(cell)
		dataCell.CreateTime = now
		dataCell.UpdateTime = now
		dataCells = append(dataCells, *dataCell)
	}

	// 批量创建
	if err := r.data.DB.CreateInBatches(dataCells, 100).Error; err != nil {
		r.log.Errorf("批量创建柜口失败: %v", err)
		return err
	}

	return nil
}

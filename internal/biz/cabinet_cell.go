package biz

import (
	"context"
	"time"
	
	"github.com/go-kratos/kratos/v2/log"
)

// CabinetCell 柜口业务实体
type CabinetCell struct {
	Id           int32     // 格口ID
	GroupId      int32     // 所属柜组ID
	CellNo       int32     // 格口编号
	CellSize     string    // 格口大小
	Status       string    // 格口状态
	LastOpenTime time.Time // 最后开启时间
	CreateTime   time.Time // 创建时间
	UpdateTime   time.Time // 更新时间
}

// CabinetCellRepo 柜口数据仓库接口
type CabinetCellRepo interface {
	// CreateCabinetCell 创建柜口
	CreateCabinetCell(ctx context.Context, cell *CabinetCell) (*CabinetCell, error)
	// UpdateCabinetCell 更新柜口
	UpdateCabinetCell(ctx context.Context, cell *CabinetCell) (*CabinetCell, error)
	// GetCabinetCellByID 根据ID获取柜口
	GetCabinetCellByID(ctx context.Context, id int32) (*CabinetCell, error)
	// GetCabinetCellByGroupAndNo 根据柜组ID和格口编号获取柜口
	GetCabinetCellByGroupAndNo(ctx context.Context, groupId, cellNo int32) (*CabinetCell, error)
	// ListCabinetCells 获取柜口列表
	ListCabinetCells(ctx context.Context, page, pageSize int32, groupId int32, status string) ([]*CabinetCell, int64, error)
	// UpdateCabinetCellStatus 更新柜口状态
	UpdateCabinetCellStatus(ctx context.Context, id int32, status string) error
	// DeleteCabinetCell 删除柜口
	DeleteCabinetCell(ctx context.Context, id int32) error
	// SearchCabinetCells 搜索柜口
	SearchCabinetCells(ctx context.Context, keyword string, page, pageSize int32, groupId int32) ([]*CabinetCell, int64, error)
	// GetCabinetCellsByGroupId 根据柜组ID获取所有柜口
	GetCabinetCellsByGroupId(ctx context.Context, groupId int32) ([]*CabinetCell, error)
	// BatchCreateCabinetCells 批量创建柜口
	BatchCreateCabinetCells(ctx context.Context, cells []*CabinetCell) error
}

// CabinetCellUsecase 柜口用例
type CabinetCellUsecase struct {
	repo CabinetCellRepo
	log  *log.Helper
}

// NewCabinetCellUsecase 创建柜口用例实例
func NewCabinetCellUsecase(repo CabinetCellRepo, logger log.Logger) *CabinetCellUsecase {
	return &CabinetCellUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// CreateCabinetCell 创建柜口
func (uc *CabinetCellUsecase) CreateCabinetCell(ctx context.Context, cell *CabinetCell) (*CabinetCell, error) {
	// 设置默认值
	if cell.CellSize == "" {
		cell.CellSize = "medium"
	}
	if cell.Status == "" {
		cell.Status = "normal"
	}
	
	return uc.repo.CreateCabinetCell(ctx, cell)
}

// UpdateCabinetCell 更新柜口
func (uc *CabinetCellUsecase) UpdateCabinetCell(ctx context.Context, cell *CabinetCell) (*CabinetCell, error) {
	return uc.repo.UpdateCabinetCell(ctx, cell)
}

// GetCabinetCellByID 根据ID获取柜口
func (uc *CabinetCellUsecase) GetCabinetCellByID(ctx context.Context, id int32) (*CabinetCell, error) {
	return uc.repo.GetCabinetCellByID(ctx, id)
}

// GetCabinetCellByGroupAndNo 根据柜组ID和格口编号获取柜口
func (uc *CabinetCellUsecase) GetCabinetCellByGroupAndNo(ctx context.Context, groupId, cellNo int32) (*CabinetCell, error) {
	return uc.repo.GetCabinetCellByGroupAndNo(ctx, groupId, cellNo)
}

// ListCabinetCells 获取柜口列表
func (uc *CabinetCellUsecase) ListCabinetCells(ctx context.Context, page, pageSize int32, groupId int32, status string) ([]*CabinetCell, int64, error) {
	return uc.repo.ListCabinetCells(ctx, page, pageSize, groupId, status)
}

// UpdateCabinetCellStatus 更新柜口状态
func (uc *CabinetCellUsecase) UpdateCabinetCellStatus(ctx context.Context, id int32, status string) error {
	return uc.repo.UpdateCabinetCellStatus(ctx, id, status)
}

// DeleteCabinetCell 删除柜口
func (uc *CabinetCellUsecase) DeleteCabinetCell(ctx context.Context, id int32) error {
	return uc.repo.DeleteCabinetCell(ctx, id)
}

// SearchCabinetCells 搜索柜口
func (uc *CabinetCellUsecase) SearchCabinetCells(ctx context.Context, keyword string, page, pageSize int32, groupId int32) ([]*CabinetCell, int64, error) {
	return uc.repo.SearchCabinetCells(ctx, keyword, page, pageSize, groupId)
}

// GetCabinetCellsByGroupId 根据柜组ID获取所有柜口
func (uc *CabinetCellUsecase) GetCabinetCellsByGroupId(ctx context.Context, groupId int32) ([]*CabinetCell, error) {
	return uc.repo.GetCabinetCellsByGroupId(ctx, groupId)
}

// BatchCreateCabinetCells 批量创建柜口
func (uc *CabinetCellUsecase) BatchCreateCabinetCells(ctx context.Context, cells []*CabinetCell) error {
	// 为每个柜口设置默认值
	for _, cell := range cells {
		if cell.CellSize == "" {
			cell.CellSize = "medium"
		}
		if cell.Status == "" {
			cell.Status = "normal"
		}
	}
	
	return uc.repo.BatchCreateCabinetCells(ctx, cells)
}

// OpenCabinetCell 开启柜口（更新最后开启时间和状态）
func (uc *CabinetCellUsecase) OpenCabinetCell(ctx context.Context, id int32) error {
	// 获取柜口信息
	cell, err := uc.repo.GetCabinetCellByID(ctx, id)
	if err != nil {
		return err
	}
	
	// 更新最后开启时间和状态
	cell.LastOpenTime = time.Now()
	if cell.Status == "normal" {
		cell.Status = "inUse"
	}
	
	_, err = uc.repo.UpdateCabinetCell(ctx, cell)
	return err
}

// CloseCabinetCell 关闭柜口（恢复为正常状态）
func (uc *CabinetCellUsecase) CloseCabinetCell(ctx context.Context, id int32) error {
	// 获取柜口信息
	cell, err := uc.repo.GetCabinetCellByID(ctx, id)
	if err != nil {
		return err
	}
	
	// 恢复为正常状态
	if cell.Status == "inUse" {
		cell.Status = "normal"
	}
	
	_, err = uc.repo.UpdateCabinetCell(ctx, cell)
	return err
}
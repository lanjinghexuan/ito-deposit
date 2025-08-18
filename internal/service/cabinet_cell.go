package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	pb "ito-deposit/api/helloworld/v1"
	"ito-deposit/internal/biz"
	"ito-deposit/internal/data"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// CabinetCellService 柜口服务结构体
type CabinetCellService struct {
	pb.UnimplementedCabinetCellServer
	RedisDb *redis.Client           // Redis客户端，用于缓存
	DB      *gorm.DB                // 数据库连接
	uc      *biz.CabinetCellUsecase // 业务用例
}

// NewCabinetCellService 创建新的柜口服务实例
func NewCabinetCellService(datas *data.Data, uc *biz.CabinetCellUsecase) *CabinetCellService {
	return &CabinetCellService{
		RedisDb: datas.Redis,
		DB:      datas.DB,
		uc:      uc,
	}
}

// 将业务实体转换为protobuf消息
func (s *CabinetCellService) convertToPbCabinetCell(cell *biz.CabinetCell) *pb.CabinetCellInfo {
	if cell == nil {
		return nil
	}

	// 详细调试：打印柜格完整信息
	log.Printf("=== 转换柜格数据 ===")
	log.Printf("ID: %d", cell.Id)
	log.Printf("GroupId: %d", cell.GroupId)
	log.Printf("CellNo: %d", cell.CellNo)
	log.Printf("CellSize: '%s' (长度: %d)", cell.CellSize, len(cell.CellSize))
	log.Printf("Status: '%s'", cell.Status)
	log.Printf("CreateTime: %v", cell.CreateTime)
	log.Printf("UpdateTime: %v", cell.UpdateTime)

	// 确保cell_size字段不为空，如果为空则设置默认值
	cellSize := cell.CellSize
	if cellSize == "" {
		cellSize = "medium" // 默认为中等尺寸
		log.Printf("警告：柜格 %d 的 CellSize 为空，设置为默认值 'medium'", cell.Id)
	}

	log.Printf("最终使用的 CellSize: '%s'", cellSize)

	return &pb.CabinetCellInfo{
		Id:             cell.Id,
		CabinetGroupId: cell.GroupId,
		CellNo:         cell.CellNo,
		CellSize:       cellSize,
		Status:         cell.Status,
		LastOpenTime:   timeToTimestamp(cell.LastOpenTime),
		CreateTime:     timeToTimestamp(cell.CreateTime),
		UpdateTime:     timeToTimestamp(cell.UpdateTime),
	}
}

// CreateCabinetCell 创建柜口
func (s *CabinetCellService) CreateCabinetCell(ctx context.Context, req *pb.CreateCabinetCellRequest) (*pb.CreateCabinetCellReply, error) {
	// 静默检查表是否存在，如果不存在则创建
	if !s.DB.Migrator().HasTable("cabinet_cells") {
		if err := s.DB.AutoMigrate(&data.CabinetCell{}); err != nil {
			return &pb.CreateCabinetCellReply{
				Code: 500,
				Msg:  "数据库表创建失败: " + err.Error(),
			}, nil
		}
	}

	// 1. 参数验证
	if req.CabinetGroupId == 0 {
		return &pb.CreateCabinetCellReply{
			Code: 400,
			Msg:  "柜组ID不能为空",
		}, nil
	}

	if req.CellNo == 0 {
		return &pb.CreateCabinetCellReply{
			Code: 400,
			Msg:  "格口编号不能为空",
		}, nil
	}

	// 2. 设置默认值
	cellSize := req.CellSize
	if cellSize == "" {
		cellSize = "medium" // 默认为中等大小
	}

	status := req.Status
	if status == "" {
		status = "normal" // 默认状态为正常
	}

	// 3. 创建业务实体
	cell := &biz.CabinetCell{
		GroupId:  req.CabinetGroupId,
		CellNo:   req.CellNo,
		CellSize: cellSize,
		Status:   status,
	}

	// 4. 调用业务层创建柜口
	createdCell, err := s.uc.CreateCabinetCell(ctx, cell)
	if err != nil {
		if strings.Contains(err.Error(), "已存在") {
			return &pb.CreateCabinetCellReply{
				Code: 400,
				Msg:  err.Error(),
			}, nil
		}
		return &pb.CreateCabinetCellReply{
			Code: 500,
			Msg:  "创建柜口失败: " + err.Error(),
		}, nil
	}

	// 5. 返回成功结果
	return &pb.CreateCabinetCellReply{
		Code:     200,
		Msg:      "柜口创建成功",
		CellId:   createdCell.Id,
		CellInfo: s.convertToPbCabinetCell(createdCell),
	}, nil
}

// UpdateCabinetCell 更新柜口
func (s *CabinetCellService) UpdateCabinetCell(ctx context.Context, req *pb.UpdateCabinetCellRequest) (*pb.UpdateCabinetCellReply, error) {
	// 1. 参数验证
	if req.Id == 0 {
		return &pb.UpdateCabinetCellReply{
			Code: 400,
			Msg:  "柜口ID不能为空",
		}, nil
	}

	// 2. 获取现有柜口信息
	existingCell, err := s.uc.GetCabinetCellByID(ctx, req.Id)
	if err != nil {
		if strings.Contains(err.Error(), "不存在") {
			return &pb.UpdateCabinetCellReply{
				Code: 404,
				Msg:  "柜口不存在",
			}, nil
		}
		return &pb.UpdateCabinetCellReply{
			Code: 500,
			Msg:  "查询柜口失败: " + err.Error(),
		}, nil
	}

	// 3. 构建更新数据（只更新非空字段）
	updateCell := &biz.CabinetCell{
		Id:           req.Id,
		GroupId:      existingCell.GroupId,
		CellNo:       existingCell.CellNo,
		CellSize:     existingCell.CellSize,
		Status:       existingCell.Status,
		LastOpenTime: existingCell.LastOpenTime,
		CreateTime:   existingCell.CreateTime,
	}

	if req.CabinetGroupId != 0 {
		updateCell.GroupId = req.CabinetGroupId
	}
	if req.CellNo != 0 {
		updateCell.CellNo = req.CellNo
	}
	if req.CellSize != "" {
		updateCell.CellSize = req.CellSize
	}
	if req.Status != "" {
		updateCell.Status = req.Status
	}
	if req.LastOpenTime != nil {
		updateCell.LastOpenTime = timestampToTime(req.LastOpenTime)
	}

	// 4. 调用业务层更新柜口
	updatedCell, err := s.uc.UpdateCabinetCell(ctx, updateCell)
	if err != nil {
		if strings.Contains(err.Error(), "已存在") {
			return &pb.UpdateCabinetCellReply{
				Code: 400,
				Msg:  err.Error(),
			}, nil
		}
		return &pb.UpdateCabinetCellReply{
			Code: 500,
			Msg:  "更新柜口失败: " + err.Error(),
		}, nil
	}

	// 5. 返回成功结果
	return &pb.UpdateCabinetCellReply{
		Code:     200,
		Msg:      "柜口更新成功",
		CellInfo: s.convertToPbCabinetCell(updatedCell),
	}, nil
}

// DeleteCabinetCell 删除柜口
func (s *CabinetCellService) DeleteCabinetCell(ctx context.Context, req *pb.DeleteCabinetCellRequest) (*pb.DeleteCabinetCellReply, error) {
	// 1. 参数验证
	if req.Id == 0 {
		return &pb.DeleteCabinetCellReply{
			Code: 400,
			Msg:  "柜口ID不能为空",
		}, nil
	}

	// 2. 调用业务层删除柜口
	err := s.uc.DeleteCabinetCell(ctx, req.Id)
	if err != nil {
		if strings.Contains(err.Error(), "不存在") {
			return &pb.DeleteCabinetCellReply{
				Code: 404,
				Msg:  "柜口不存在",
			}, nil
		}
		if strings.Contains(err.Error(), "正在使用") {
			return &pb.DeleteCabinetCellReply{
				Code: 400,
				Msg:  err.Error(),
			}, nil
		}
		return &pb.DeleteCabinetCellReply{
			Code: 500,
			Msg:  "删除柜口失败: " + err.Error(),
		}, nil
	}

	return &pb.DeleteCabinetCellReply{
		Code:    200,
		Msg:     "柜口删除成功",
		Success: true,
	}, nil
}

// GetCabinetCell 获取单个柜口
func (s *CabinetCellService) GetCabinetCell(ctx context.Context, req *pb.GetCabinetCellRequest) (*pb.GetCabinetCellReply, error) {
	// 1. 参数验证
	if req.Id == 0 {
		return &pb.GetCabinetCellReply{
			Code: 400,
			Msg:  "柜口ID不能为空",
		}, nil
	}

	// 2. 调用业务层获取柜口
	cell, err := s.uc.GetCabinetCellByID(ctx, req.Id)
	if err != nil {
		if strings.Contains(err.Error(), "不存在") {
			return &pb.GetCabinetCellReply{
				Code: 404,
				Msg:  "柜口不存在",
			}, nil
		}
		return &pb.GetCabinetCellReply{
			Code: 500,
			Msg:  "查询柜口失败: " + err.Error(),
		}, nil
	}

	// 3. 返回成功结果
	return &pb.GetCabinetCellReply{
		Code:     200,
		Msg:      "查询成功",
		CellInfo: s.convertToPbCabinetCell(cell),
	}, nil
}

// ListCabinetCells 获取柜口列表
func (s *CabinetCellService) ListCabinetCells(ctx context.Context, req *pb.ListCabinetCellsRequest) (*pb.ListCabinetCellsReply, error) {
	// 1. 分页参数处理
	page := req.Page
	if page <= 0 {
		page = 1
	}

	pageSize := req.Size
	if pageSize <= 0 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100 // 限制最大页面大小
	}

	// 2. 调用业务层获取柜口列表
	cells, total, err := s.uc.ListCabinetCells(ctx, int32(page), int32(pageSize), req.CabinetGroupId, req.Status)
	if err != nil {
		return &pb.ListCabinetCellsReply{
			Code: 500,
			Msg:  "查询柜口列表失败: " + err.Error(),
		}, nil
	}

	// 3. 构建返回数据
	var cellInfos []*pb.CabinetCellInfo
	for _, cell := range cells {
		cellInfo := s.convertToPbCabinetCell(cell)
		cellInfos = append(cellInfos, cellInfo)
	}

	return &pb.ListCabinetCellsReply{
		Code:  200,
		Msg:   "查询成功",
		Cells: cellInfos,
		Total: total,
	}, nil
}

// SearchCabinetCells 搜索柜口
func (s *CabinetCellService) SearchCabinetCells(ctx context.Context, req *pb.SearchCabinetCellsRequest) (*pb.SearchCabinetCellsReply, error) {
	// 1. 分页参数处理
	page := req.Page
	if page <= 0 {
		page = 1
	}

	pageSize := req.Size
	if pageSize <= 0 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}

	// 2. 调用业务层搜索柜口
	cells, total, err := s.uc.SearchCabinetCells(ctx, req.Keyword, int32(page), int32(pageSize), req.CabinetGroupId)
	if err != nil {
		return &pb.SearchCabinetCellsReply{
			Code: 500,
			Msg:  "搜索柜口失败: " + err.Error(),
		}, nil
	}

	// 3. 构建返回数据
	var cellInfos []*pb.CabinetCellInfo
	for _, cell := range cells {
		cellInfo := s.convertToPbCabinetCell(cell)
		cellInfos = append(cellInfos, cellInfo)
	}

	return &pb.SearchCabinetCellsReply{
		Code:  200,
		Msg:   "搜索成功",
		Cells: cellInfos,
		Total: total,
	}, nil
}

// GetCabinetCellsByGroup 根据柜组获取所有柜口
func (s *CabinetCellService) GetCabinetCellsByGroup(ctx context.Context, req *pb.GetCabinetCellsByGroupRequest) (*pb.GetCabinetCellsByGroupReply, error) {
	// 1. 参数验证
	if req.CabinetGroupId == 0 {
		return &pb.GetCabinetCellsByGroupReply{
			Code: 400,
			Msg:  "柜组ID不能为空",
		}, nil
	}

	// 2. 调用业务层获取柜口列表
	cells, err := s.uc.GetCabinetCellsByGroupId(ctx, req.CabinetGroupId)
	if err != nil {
		return &pb.GetCabinetCellsByGroupReply{
			Code: 500,
			Msg:  "查询柜口失败: " + err.Error(),
		}, nil
	}

	// 3. 构建返回数据
	var cellInfos []*pb.CabinetCellInfo
	for _, cell := range cells {
		cellInfo := s.convertToPbCabinetCell(cell)
		cellInfos = append(cellInfos, cellInfo)
	}

	return &pb.GetCabinetCellsByGroupReply{
		Code:  200,
		Msg:   "查询成功",
		Cells: cellInfos,
		Total: int64(len(cellInfos)),
	}, nil
}

// BatchCreateCabinetCells 批量创建柜口
func (s *CabinetCellService) BatchCreateCabinetCells(ctx context.Context, req *pb.BatchCreateCabinetCellsRequest) (*pb.BatchCreateCabinetCellsReply, error) {
	// 1. 参数验证
	if req.CabinetGroupId == 0 {
		return &pb.BatchCreateCabinetCellsReply{
			Code: 400,
			Msg:  "柜组ID不能为空",
		}, nil
	}

	if req.StartNo >= req.EndNo {
		return &pb.BatchCreateCabinetCellsReply{
			Code: 400,
			Msg:  "起始编号必须小于结束编号",
		}, nil
	}

	// 2. 设置默认值
	cellSize := req.CellSize
	if cellSize == "" {
		cellSize = "medium"
	}

	// 3. 构建批量创建的柜口列表
	var cells []*biz.CabinetCell
	for cellNo := req.StartNo; cellNo <= req.EndNo; cellNo++ {
		cell := &biz.CabinetCell{
			GroupId:  req.CabinetGroupId,
			CellNo:   cellNo,
			CellSize: cellSize,
			Status:   "normal",
		}
		cells = append(cells, cell)
	}

	// 4. 调用业务层批量创建
	err := s.uc.BatchCreateCabinetCells(ctx, cells)
	if err != nil {
		return &pb.BatchCreateCabinetCellsReply{
			Code: 500,
			Msg:  "批量创建柜口失败: " + err.Error(),
		}, nil
	}

	return &pb.BatchCreateCabinetCellsReply{
		Code:    200,
		Msg:     "批量创建柜口成功",
		Count:   int32(len(cells)),
		Success: true,
	}, nil
}

// OpenCabinetCell 开启柜口
func (s *CabinetCellService) OpenCabinetCell(ctx context.Context, req *pb.OpenCabinetCellRequest) (*pb.OpenCabinetCellReply, error) {
	// 1. 参数验证
	if req.Id == 0 {
		return &pb.OpenCabinetCellReply{
			Code: 400,
			Msg:  "柜口ID不能为空",
		}, nil
	}

	// 2. 调用业务层开启柜口
	err := s.uc.OpenCabinetCell(ctx, req.Id)
	if err != nil {
		if strings.Contains(err.Error(), "不存在") {
			return &pb.OpenCabinetCellReply{
				Code: 404,
				Msg:  "柜口不存在",
			}, nil
		}
		return &pb.OpenCabinetCellReply{
			Code: 500,
			Msg:  "开启柜口失败: " + err.Error(),
		}, nil
	}

	return &pb.OpenCabinetCellReply{
		Code:    200,
		Msg:     "柜口开启成功",
		Success: true,
	}, nil
}

// CloseCabinetCell 关闭柜口
func (s *CabinetCellService) CloseCabinetCell(ctx context.Context, req *pb.CloseCabinetCellRequest) (*pb.CloseCabinetCellReply, error) {
	// 1. 参数验证
	if req.Id == 0 {
		return &pb.CloseCabinetCellReply{
			Code: 400,
			Msg:  "柜口ID不能为空",
		}, nil
	}

	// 2. 调用业务层关闭柜口
	err := s.uc.CloseCabinetCell(ctx, req.Id)
	if err != nil {
		if strings.Contains(err.Error(), "不存在") {
			return &pb.CloseCabinetCellReply{
				Code: 404,
				Msg:  "柜口不存在",
			}, nil
		}
		return &pb.CloseCabinetCellReply{
			Code: 500,
			Msg:  "关闭柜口失败: " + err.Error(),
		}, nil
	}

	return &pb.CloseCabinetCellReply{
		Code:    200,
		Msg:     "柜口关闭成功",
		Success: true,
	}, nil
}
func (s *CabinetCellService) CellStatus(ctx context.Context, req *pb.CellStatusReq) (*pb.CellStatusRes, error) {
	var cell []data.CabinetCell
	err := s.DB.Debug().Where("status = ?", "abnormal").Find(&cell).Error
	if err != nil {
		return &pb.CellStatusRes{
			Code: 500,
			Msg:  "监控状态失败",
		}, nil
	}
	if len(cell) == 0 {
		return &pb.CellStatusRes{
			Code: 200,
			Msg:  "当前无异常格口",
		}, nil
	}

	content := fmt.Sprintf("❗ 快递柜异常告警：发现 %d 个格口状态为 abnormal", len(cell))
	webhookURL := "https://open.feishu.cn/open-apis/bot/v2/hook/f298817f-d96b-4db0-9597-00ffbeb99c9f"
	err = sendFeishuAlert(webhookURL, content)
	if err != nil {
		log.Printf("发送飞书告警失败: %v", err)
	}
	return &pb.CellStatusRes{
		Code: 200,
		Msg:  "监控状态",
	}, nil
}

// 飞书告警消息发送（直接复制你 main.go 的函数）
func sendFeishuAlert(webhookURL, content string) error {
	msg := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]string{
			"text": content,
		},
	}
	jsonData, _ := json.Marshal(msg) // 简单处理，实际可加 error 判断

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("飞书返回状态码异常: %d", resp.StatusCode)
	}
	return nil
}

package service

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	pb "ito-deposit/api/helloworld/v1"
	"ito-deposit/internal/data"
	"strings"
	"time"
)

// timeToTimestamp 将 time.Time 转换为 protobuf Timestamp
func timeToTimestamp(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}

// timestampToTime 将 protobuf Timestamp 转换为 time.Time
func timestampToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

// GroupService 柜组服务结构体
type GroupService struct {
	pb.UnimplementedGroupServer
	RedisDb *redis.Client // Redis客户端，用于缓存
	DB      *gorm.DB      // 数据库连接
}

// NewGroupService 创建新的柜组服务实例
func NewGroupService(datas *data.Data) *GroupService {
	return &GroupService{
		RedisDb: datas.Redis,
		DB:      datas.DB,
	}
}

// CreateGroup 创建柜组
func (s *GroupService) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.CreateGroupReply, error) {
	// 静默检查表是否存在，如果不存在则创建
	if !s.DB.Migrator().HasTable("cabinet_groups") {
		if err := s.DB.AutoMigrate(&data.CabinetGroup{}); err != nil {
			return &pb.CreateGroupReply{
				Code: 500,
				Msg:  "数据库表创建失败: " + err.Error(),
			}, nil
		}
	}

	// 1. 参数验证
	if req.LocationPointId == 0 {
		return &pb.CreateGroupReply{
			Code: 400,
			Msg:  "寄存点ID不能为空",
		}, nil
	}

	if req.GroupName == "" {
		return &pb.CreateGroupReply{
			Code: 400,
			Msg:  "柜组名称不能为空",
		}, nil
	}

	if req.GroupCode == "" {
		return &pb.CreateGroupReply{
			Code: 400,
			Msg:  "柜组编码不能为空",
		}, nil
	}

	if req.TotalCells <= 0 {
		return &pb.CreateGroupReply{
			Code: 400,
			Msg:  "总格口数必须大于0",
		}, nil
	}

	if req.StartNo >= req.EndNo {
		return &pb.CreateGroupReply{
			Code: 400,
			Msg:  "起始编号必须小于结束编号",
		}, nil
	}

	// 2. 检查柜组编码是否已存在
	var existingGroup data.CabinetGroup
	err := s.DB.Where("group_code = ?", req.GroupCode).First(&existingGroup).Error
	if err == nil {
		return &pb.CreateGroupReply{
			Code: 400,
			Msg:  "柜组编码已存在",
		}, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return &pb.CreateGroupReply{
			Code: 500,
			Msg:  "数据库查询失败: " + err.Error(),
		}, nil
	}

	// 3. 解析安装时间
	var installTime time.Time
	if req.InstallTime != nil {
		installTime = timestampToTime(req.InstallTime)
	}

	// 4. 设置默认值
	groupType := req.GroupType
	if groupType == "" {
		groupType = "standard" // 默认为标准类型
	}

	// 5. 创建柜组记录（ID由数据库自动生成）
	group := data.CabinetGroup{
		LocationPointId: req.LocationPointId,
		GroupName:       req.GroupName,
		GroupCode:       req.GroupCode,
		GroupType:       groupType,
		Status:          "normal", // 默认状态为正常
		TotalCells:      req.TotalCells,
		StartNo:         req.StartNo,
		EndNo:           req.EndNo,
		InstallTime:     installTime,
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
	}

	// 6. 保存到数据库
	err = s.DB.Create(&group).Error
	if err != nil {
		return &pb.CreateGroupReply{
			Code: 500,
			Msg:  "创建柜组失败: " + err.Error(),
		}, nil
	}

	// 7. 返回成功结果
	return &pb.CreateGroupReply{
		Code:    200,
		Msg:     "柜组创建成功",
		GroupId: group.Id, // 返回自动生成的ID
	}, nil
}

// UpdateGroup 更新柜组
func (s *GroupService) UpdateGroup(ctx context.Context, req *pb.UpdateGroupRequest) (*pb.UpdateGroupReply, error) {
	// 1. 参数验证
	if req.Id == 0 {
		return &pb.UpdateGroupReply{
			Code: 400,
			Msg:  "柜组ID不能为空",
		}, nil
	}

	// 2. 查询柜组是否存在
	var group data.CabinetGroup
	err := s.DB.Where("id = ?", req.Id).First(&group).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &pb.UpdateGroupReply{
			Code: 404,
			Msg:  "柜组不存在",
		}, nil
	} else if err != nil {
		return &pb.UpdateGroupReply{
			Code: 500,
			Msg:  "查询柜组失败: " + err.Error(),
		}, nil
	}

	// 3. 如果更新柜组编码，检查是否与其他柜组重复
	if req.GroupCode != "" && req.GroupCode != group.GroupCode {
		var existingGroup data.CabinetGroup
		err = s.DB.Where("group_code = ? AND id != ?", req.GroupCode, req.Id).First(&existingGroup).Error
		if err == nil {
			return &pb.UpdateGroupReply{
				Code: 400,
				Msg:  "柜组编码已被其他柜组使用",
			}, nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.UpdateGroupReply{
				Code: 500,
				Msg:  "检查柜组编码失败: " + err.Error(),
			}, nil
		}
	}

	// 4. 验证编号范围
	startNo := req.StartNo
	endNo := req.EndNo
	if startNo == 0 {
		startNo = group.StartNo
	}
	if endNo == 0 {
		endNo = group.EndNo
	}
	if startNo >= endNo {
		return &pb.UpdateGroupReply{
			Code: 400,
			Msg:  "起始编号必须小于结束编号",
		}, nil
	}

	// 5. 构建更新数据
	updateData := map[string]interface{}{
		"update_time": time.Now(),
	}

	// 只更新非空字段
	if req.LocationPointId != 0 {
		updateData["location_point_id"] = req.LocationPointId
	}
	if req.GroupName != "" {
		updateData["group_name"] = req.GroupName
	}
	if req.GroupCode != "" {
		updateData["group_code"] = req.GroupCode
	}
	if req.GroupType != "" {
		updateData["group_type"] = req.GroupType
	}
	if req.Status != "" {
		updateData["status"] = req.Status
	}
	if req.TotalCells > 0 {
		updateData["total_cells"] = req.TotalCells
	}
	if req.StartNo > 0 {
		updateData["start_no"] = req.StartNo
	}
	if req.EndNo > 0 {
		updateData["end_no"] = req.EndNo
	}
	if req.InstallTime != nil {
		installTime := timestampToTime(req.InstallTime)
		updateData["install_time"] = installTime
	}

	// 6. 执行更新
	err = s.DB.Model(&group).Updates(updateData).Error
	if err != nil {
		return &pb.UpdateGroupReply{
			Code: 500,
			Msg:  "更新柜组失败: " + err.Error(),
		}, nil
	}

	// 7. 重新查询更新后的数据
	err = s.DB.Where("id = ?", req.Id).First(&group).Error
	if err != nil {
		return &pb.UpdateGroupReply{
			Code: 500,
			Msg:  "查询更新后的柜组失败: " + err.Error(),
		}, nil
	}

	// 8. 构建返回数据
	groupInfo := &pb.GroupInfo{
		Id:              group.Id,
		LocationPointId: group.LocationPointId,
		GroupName:       group.GroupName,
		GroupCode:       group.GroupCode,
		GroupType:       group.GroupType,
		Status:          group.Status,
		TotalCells:      group.TotalCells,
		StartNo:         group.StartNo,
		EndNo:           group.EndNo,
		InstallTime:     timeToTimestamp(group.InstallTime),
		CreateTime:      timeToTimestamp(group.CreateTime),
		UpdateTime:      timeToTimestamp(group.UpdateTime),
	}

	return &pb.UpdateGroupReply{
		Code:  200,
		Msg:   "柜组更新成功",
		Group: groupInfo,
	}, nil
}

// DeleteGroup 删除柜组
func (s *GroupService) DeleteGroup(ctx context.Context, req *pb.DeleteGroupRequest) (*pb.DeleteGroupReply, error) {
	// 1. 参数验证
	if req.Id == 0 {
		return &pb.DeleteGroupReply{
			Code: 400,
			Msg:  "柜组ID不能为空",
		}, nil
	}

	// 2. 查询柜组是否存在
	var group data.CabinetGroup
	err := s.DB.Where("id = ?", req.Id).First(&group).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &pb.DeleteGroupReply{
			Code: 404,
			Msg:  "柜组不存在",
		}, nil
	} else if err != nil {
		return &pb.DeleteGroupReply{
			Code: 500,
			Msg:  "查询柜组失败: " + err.Error(),
		}, nil
	}

	// 3. 检查柜组是否已被禁用
	if group.Status == "disabled" {
		return &pb.DeleteGroupReply{
			Code: 400,
			Msg:  "柜组已被禁用",
		}, nil
	}

	// 4. 软删除：更新状态为禁用
	err = s.DB.Model(&group).Updates(map[string]interface{}{
		"status":      "disabled",
		"update_time": time.Now(),
	}).Error
	if err != nil {
		return &pb.DeleteGroupReply{
			Code: 500,
			Msg:  "删除柜组失败: " + err.Error(),
		}, nil
	}

	return &pb.DeleteGroupReply{
		Code:    200,
		Msg:     "柜组删除成功",
		Success: true,
	}, nil
}

// GetGroup 获取单个柜组
func (s *GroupService) GetGroup(ctx context.Context, req *pb.GetGroupRequest) (*pb.GetGroupReply, error) {
	// 1. 参数验证
	if req.Id == 0 {
		return &pb.GetGroupReply{
			Code: 400,
			Msg:  "柜组ID不能为空",
		}, nil
	}

	// 2. 查询柜组
	var group data.CabinetGroup
	err := s.DB.Where("id = ?", req.Id).First(&group).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &pb.GetGroupReply{
			Code: 404,
			Msg:  "柜组不存在",
		}, nil
	} else if err != nil {
		return &pb.GetGroupReply{
			Code: 500,
			Msg:  "查询柜组失败: " + err.Error(),
		}, nil
	}

	// 3. 构建返回数据
	groupInfo := &pb.GroupInfo{
		Id:              group.Id,
		LocationPointId: group.LocationPointId,
		GroupName:       group.GroupName,
		GroupCode:       group.GroupCode,
		GroupType:       group.GroupType,
		Status:          group.Status,
		TotalCells:      group.TotalCells,
		StartNo:         group.StartNo,
		EndNo:           group.EndNo,
		InstallTime:     timeToTimestamp(group.InstallTime),
		CreateTime:      timeToTimestamp(group.CreateTime),
		UpdateTime:      timeToTimestamp(group.UpdateTime),
	}

	return &pb.GetGroupReply{
		Code:  200,
		Msg:   "查询成功",
		Group: groupInfo,
	}, nil
}

// ListGroup 获取柜组列表
func (s *GroupService) ListGroup(ctx context.Context, req *pb.ListGroupRequest) (*pb.ListGroupReply, error) {
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

	offset := (page - 1) * pageSize

	// 2. 构建查询条件
	db := s.DB.Model(&data.CabinetGroup{})

	// 按寄存点ID过滤
	if req.LocationPointId != 0 {
		db = db.Where("location_point_id = ?", req.LocationPointId)
	}

	// 按状态过滤
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	// 按类型过滤
	if req.GroupType != "" {
		db = db.Where("group_type = ?", req.GroupType)
	}

	// 排除已禁用的柜组（除非明确查询禁用状态）
	if req.Status != "disabled" {
		db = db.Where("status != ?", "disabled")
	}

	// 3. 获取总数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return &pb.ListGroupReply{
			Code: 500,
			Msg:  "查询总数失败: " + err.Error(),
		}, nil
	}

	// 4. 获取列表数据
	var groups []data.CabinetGroup
	err = db.Order("create_time DESC").
		Limit(int(pageSize)).
		Offset(int(offset)).
		Find(&groups).Error
	if err != nil {
		return &pb.ListGroupReply{
			Code: 500,
			Msg:  "查询柜组列表失败: " + err.Error(),
		}, nil
	}

	// 5. 构建返回数据
	var groupInfos []*pb.GroupInfo
	for _, group := range groups {
		groupInfo := &pb.GroupInfo{
			Id:              group.Id,
			LocationPointId: group.LocationPointId,
			GroupName:       group.GroupName,
			GroupCode:       group.GroupCode,
			GroupType:       group.GroupType,
			Status:          group.Status,
			TotalCells:      group.TotalCells,
			StartNo:         group.StartNo,
			EndNo:           group.EndNo,
			InstallTime:     timeToTimestamp(group.InstallTime),
			CreateTime:      timeToTimestamp(group.CreateTime),
			UpdateTime:      timeToTimestamp(group.UpdateTime),
		}
		groupInfos = append(groupInfos, groupInfo)
	}

	return &pb.ListGroupReply{
		Code:   200,
		Msg:    "查询成功",
		Groups: groupInfos,
		Total:  total,
	}, nil
}

// SearchGroup 搜索柜组
func (s *GroupService) SearchGroup(ctx context.Context, req *pb.SearchGroupRequest) (*pb.SearchGroupReply, error) {
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

	offset := (page - 1) * pageSize

	// 2. 构建查询条件
	db := s.DB.Model(&data.CabinetGroup{})

	// 关键词搜索（柜组名称或编码）
	if req.Keyword != "" {
		keyword := "%" + strings.TrimSpace(req.Keyword) + "%"
		db = db.Where("group_name LIKE ? OR group_code LIKE ?", keyword, keyword)
	}

	// 按寄存点ID过滤
	if req.LocationPointId != 0 {
		db = db.Where("location_point_id = ?", req.LocationPointId)
	}

	// 按状态过滤
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	// 按类型过滤
	if req.GroupType != "" {
		db = db.Where("group_type = ?", req.GroupType)
	}

	// 排除已禁用的柜组（除非明确查询禁用状态）
	if req.Status != "disabled" {
		db = db.Where("status != ?", "disabled")
	}

	// 3. 获取总数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return &pb.SearchGroupReply{
			Code: 500,
			Msg:  "搜索失败: " + err.Error(),
		}, nil
	}

	// 4. 获取搜索结果
	var groups []data.CabinetGroup
	err = db.Order("create_time DESC").
		Limit(int(pageSize)).
		Offset(int(offset)).
		Find(&groups).Error
	if err != nil {
		return &pb.SearchGroupReply{
			Code: 500,
			Msg:  "搜索失败: " + err.Error(),
		}, nil
	}

	// 5. 构建返回数据
	var groupInfos []*pb.GroupInfo
	for _, group := range groups {
		groupInfo := &pb.GroupInfo{
			Id:              group.Id,
			LocationPointId: group.LocationPointId,
			GroupName:       group.GroupName,
			GroupCode:       group.GroupCode,
			GroupType:       group.GroupType,
			Status:          group.Status,
			TotalCells:      group.TotalCells,
			StartNo:         group.StartNo,
			EndNo:           group.EndNo,
			InstallTime:     timeToTimestamp(group.InstallTime),
			CreateTime:      timeToTimestamp(group.CreateTime),
			UpdateTime:      timeToTimestamp(group.UpdateTime),
		}
		groupInfos = append(groupInfos, groupInfo)
	}

	return &pb.SearchGroupReply{
		Code:   200,
		Msg:    "搜索成功",
		Groups: groupInfos,
		Total:  total,
	}, nil
}

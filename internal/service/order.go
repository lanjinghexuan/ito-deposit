package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	pb "ito-deposit/api/helloworld/v1"
	"ito-deposit/internal/basic/pkg"
	"ito-deposit/internal/data"
	"strconv"
)

type OrderService struct {
	pb.UnimplementedOrderServer
	RedisDb *redis.Client
	DB      *gorm.DB
}

func NewOrderService(datas *data.Data) *OrderService {
	return &OrderService{
		RedisDb: datas.Redis,
		DB:      datas.DB,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderReply, error) {
	//查询柜子的是否已经寄存
	var order data.LockerStorages
	s.DB.Where("cabinet_id  = ?", req.CabinetId).First(&order)

	if order.CabinetId == req.CabinetId {
		return nil, errors.New("point_id already exists")
	}

	//随机生成订单编号
	NewString := uuid.NewString()

	var locker = data.LockerOrders{
		OrderNumber:         NewString,
		UserId:              uint64(req.UserId),
		ScheduledDuration:   int32(req.ScheduledDuration),
		Price:               req.Price,
		Discount:            req.Discount,
		AmountPaid:          req.AmountPaid,
		StorageLocationName: req.StorageLocationName,
		CabinetId:           int32(req.CabinetId),
		Status:              int8(req.Status),
		DepositStatus:       int8(req.DepositStatus),
		Title:               req.Title,
	}
	err := s.DB.Create(&locker).Error
	if err != nil {
		return nil, fmt.Errorf("create order failed: %v", err)
	}
	var lockerpick = data.LockerStorages{
		OrderId:   req.OrderId,
		CabinetId: req.CabinetId,
		Status:    req.Status,
		UserId:    req.UserId,
	}
	err = s.DB.Create(&lockerpick).Error
	if err != nil {
		return nil, fmt.Errorf("create order failed: %v", err)
	}

	return &pb.CreateOrderReply{
		Msg: "订单发送,并记录",
	}, nil

}
func (s *OrderService) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.UpdateOrderReply, error) {
	// 1. 查询订单是否存在
	var lockerOrder data.LockerOrders
	if err := s.DB.Where("id = ?", req.Id).First(&lockerOrder).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "订单不存在")
		}
		return nil, status.Errorf(codes.Internal, "数据库错误: %v", err)
	}

	// 2. 查询规则表获取小时费率
	var pricingRule data.LockerPricingRules
	if err := s.DB.Where("hourly_rate = ?", req.HourlyRate).First(&pricingRule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "费率规则不存在")
		}
		return nil, status.Errorf(codes.Internal, "数据库错误: %v", err)
	}

	// 3. 计算总价格（使用实际时长）
	totalPrice := pricingRule.HourlyRate * float64(lockerOrder.ActualDuration)

	// 4. 更新订单信息
	updateData := data.LockerOrders{
		ActualDuration: int32(req.ActualDuration),
		AmountPaid:     totalPrice,
		Status:         int8(req.Status),
		Price:          req.Price,
		DepositStatus:  int8(req.DepositStatus),
	}

	if err := s.DB.Model(&lockerOrder).Updates(&updateData).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "更新订单失败: %v", err)
	}

	// 5. 生成支付链接（使用计算后的金额）
	totalAmountStr := strconv.FormatFloat(totalPrice, 'f', 2, 64)

	payUrl := pkg.Pay(lockerOrder.Title, lockerOrder.OrderNumber, totalAmountStr)

	// 6. 返回响应
	return &pb.UpdateOrderReply{
		PayUrl: payUrl,
	}, nil
}
func (s *OrderService) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderReply, error) {
	var lockerOrder data.LockerOrders

	// 根据订单 ID 查询
	if err := s.DB.Where("id = ?", req.Id).First(&lockerOrder).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}

	// 只有进行中或已取出的订单可以关闭
	if lockerOrder.Status != 1 && lockerOrder.Status != 2 {
		return nil, status.Errorf(codes.FailedPrecondition, "订单状态不允许关闭")
	}

	// 更新状态为已完成（软删除）
	if err := s.DB.Model(&lockerOrder).Delete(&lockerOrder).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "关闭订单失败: %v", err)
	}

	return &pb.DeleteOrderReply{
		Success: true,
	}, nil
}

func (s *OrderService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderReply, error) {
	return &pb.GetOrderReply{}, nil
}

func (s *OrderService) ListOrder(ctx context.Context, req *pb.ListOrderRequest) (*pb.ListOrderReply, error) {
	var list []data.LockerOrders

	// 分页参数处理
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

	// 构建查询条件
	db := s.DB.Model(&data.LockerOrders{})

	if req.StorageLocationName != "" {
		db = db.Where("storage_location_name LIKE ?", "%"+req.StorageLocationName+"%")
	}

	// 添加排序
	db = db.Order("id DESC")

	// 获取总条数
	var total int64
	db.Count(&total)

	// ✅ 正确使用 Limit & Offset
	if err := db.Limit(int(pageSize)).Offset(int(offset)).Find(&list).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "数据库查询失败: %v", err)
	}

	// 构造返回数据
	var orders []*pb.OrderInfo
	for _, v := range list {
		orders = append(orders, &pb.OrderInfo{
			Id:                  v.Id,
			OrderNumber:         v.OrderNumber,
			UserId:              int64(v.UserId),
			ScheduledDuration:   int64(v.ScheduledDuration),
			ActualDuration:      int64(v.ActualDuration),
			Price:               v.Price,
			Discount:            v.Discount,
			AmountPaid:          v.AmountPaid,
			StorageLocationName: v.StorageLocationName,
			CabinetId:           int64(v.CabinetId),
			Status:              int64(v.Status),
			DepositStatus:       int64(v.DepositStatus),
		})
	}

	// 返回分页结果
	return &pb.ListOrderReply{
		Orders: orders,
		Total:  total,
	}, nil
}

func (s *OrderService) ShowOrder(ctx context.Context, req *pb.ShowOrderRequest) (*pb.ShowOrderReply, error) {
	var order data.LockerOrders

	fmt.Println("收到的订单ID:", req.Id) // 打印实际收到的ID
	if req.Id <= 0 {
		return &pb.ShowOrderReply{
			Msg:   "无效的订单ID",
			Order: nil,
		}, nil
	}

	// 执行查询
	result := s.DB.Where("id = ?", req.Id).First(&order)

	// 判断是否查询到数据
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &pb.ShowOrderReply{
				Msg:   "订单不存在",
				Order: nil,
			}, nil
		}
		// 数据库错误
		return &pb.ShowOrderReply{
			Msg:   "数据库错误：" + result.Error.Error(),
			Order: nil,
		}, nil
	}

	orderInfo := &pb.OrderInfo{
		Id:                  order.Id,
		OrderNumber:         order.OrderNumber,
		UserId:              int64(order.UserId),
		ScheduledDuration:   int64(order.ScheduledDuration),
		ActualDuration:      int64(order.ActualDuration),
		Price:               order.Price,
		Discount:            order.Discount,
		AmountPaid:          order.AmountPaid,
		StorageLocationName: order.StorageLocationName,
		CabinetId:           int64(order.CabinetId),
		Status:              int64(order.Status), // 如果protobuf中是int32类型
		DepositStatus:       int64(order.DepositStatus),
	}

	// 返回订单信息
	return &pb.ShowOrderReply{
		Msg:   "success",
		Order: []*pb.OrderInfo{orderInfo},
	}, nil
}

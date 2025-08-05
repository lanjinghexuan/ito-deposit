package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	pb "ito-deposit/api/helloworld/v1"
	"ito-deposit/internal/basic/pkg"
	"ito-deposit/internal/data"
	"strconv"
	"time"
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
func (s *OrderService) CreateLockerStorage(ctx context.Context, req *pb.CreateLockerStorageRequest) (*pb.CreateLockerStorageReply, error) {
	allowed, err := s.allowWithLeakyBucket(ctx, strconv.FormatInt(req.UserId, 10), 1, 20)

	if err != nil {
		return nil, fmt.Errorf("限流检查失败: %w", err)
	}

	if !allowed {
		return nil, fmt.Errorf("请求过于频繁，请稍后再试")
	}

	// 2. 业务层面的限流检查（如检查用户今日订单数）
	var locker data.LockerOrders
	oneMinuteAgo := time.Now().Add(-1 * time.Minute)

	// 统计该用户1分钟内的订单数量
	var count int64
	if err := s.DB.Model(&locker).
		Where("user_id = ? AND create_time >= ?", req.UserId, oneMinuteAgo).
		Count(&count).Error; err != nil {
		return nil, fmt.Errorf("查询订单数量失败: %w", err)
	}

	// 超过限制返回错误
	if count >= 10 {
		return nil, fmt.Errorf("每分钟最多创建10个订单，请稍后再试")
	}
	// 1. 参数校验
	if req.CabinetId <= 0 || req.ExpireTime <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "无效参数：柜子ID、用户ID、过期时间不能为空")
	}
	// 转换过期时间（假设 req.ExpireTime 为秒级时间戳）
	expireTime := time.Unix(req.ExpireTime, 0)
	if expireTime.Before(time.Now()) {
		return nil, status.Errorf(codes.InvalidArgument, "过期时间不能早于当前时间")
	}
	// 2. 创建寄存记录（写入 locker_storages 表）
	storage := &data.LockerStorages{
		OrderId:   req.OrderId,
		CabinetId: req.CabinetId,
		Status:    int64(req.Status), // 1-寄存中，0-待寄存
		UserId:    (req.UserId),
	}
	if err := s.DB.Create(storage).Error; err != nil {
		log.Errorw(ctx, "创建寄存记录失败", "err", err, "order_id", req.OrderId)
		return nil, status.Errorf(codes.Internal, "创建寄存记录失败：%v", err)
	}

	// 3. 创建关联的定时任务
	if err := s.createTimerTasks(ctx, req.OrderId, int64(storage.Id), expireTime); err != nil {
		// 任务创建失败，回滚寄存记录
		_ = s.DB.Delete(&data.LockerStorages{}, storage.Id)
		log.Errorw(ctx, "创建定时任务失败，回滚寄存记录", "err", err, "storage_id", storage.Id)
		return nil, status.Errorf(codes.Internal, "创建定时任务失败：%v", err)
	}

	return &pb.CreateLockerStorageReply{
		Msg: "寄存时间的创建",
	}, nil
}

// createTimerTasks 为寄存订单创建三类定时任务
func (s *OrderService) createTimerTasks(ctx context.Context, orderID, storageID int64, expireTime time.Time) error {
	// 定义任务类型及执行时间偏移量（相对于过期时间）
	taskDefines := []struct {
		taskType string        // 任务类型：remind-到期提醒，check-超时检查，handle-超时处理
		offset   time.Duration // 时间偏移（提前/延后）
	}{
		{"remind", time.Hour},      // 到期提醒
		{"handle", 12 * time.Hour}, // 到期后24小时执行超时处理
	}

	// 批量创建任务
	for _, def := range taskDefines {
		// 计算任务执行时间
		executeTime := expireTime.Add(def.offset)
		if executeTime.Before(time.Now()) {
			return fmt.Errorf("任务执行时间无效（%s）", def.taskType)
		}

		// 生成唯一任务编码（避免重复）
		taskCode := fmt.Sprintf("%s_%d_%s", def.taskType, orderID, uuid.NewString()[:8])

		// 写入 timer_tasks 表
		task := &data.TimerTasks{
			TaskCode:        taskCode,
			OrderId:         uint64(orderID),
			LockerStorageId: uint64(storageID),
			TaskType:        def.taskType,
			CronExpression:  "",
			ExecuteTime:     executeTime,
			Status:          1, // 1-待执行
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		if err := s.DB.Create(task).Error; err != nil {
			return fmt.Errorf("创建%s任务失败：%v", def.taskType, err)
		}
		log.Debugw(ctx, "定时任务创建成功", "task_code", taskCode, "execute_time", executeTime)
	}
	return nil
}

// 令牌桶算法实现 - 允许一定程度的突发流量
func (s *OrderService) allowWithTokenBucket(ctx context.Context, userId string, rate int, capacity int) (bool, error) {
	// 每个用户一个独立的令牌桶
	key := fmt.Sprintf("token_bucket:order:%s", userId)

	// 使用Lua脚本保证操作原子性
	luaScript := `
		local key = KEYS[1]
		local rate = tonumber(ARGV[1])
		local capacity = tonumber(ARGV[2])
		local now = tonumber(ARGV[3])
		local take = tonumber(ARGV[4])
		
		-- 初始化桶
		local bucket = redis.call('hgetall', key)
		local currentTokens = capacity
		local lastRefillTime = now
		
		-- 如果桶存在，获取当前状态
		if #bucket == 4 then
			currentTokens = tonumber(bucket[2])
			lastRefillTime = tonumber(bucket[4])
		end
		
		-- 计算需要补充的令牌
		local elapsed = now - lastRefillTime
		local newTokens = elapsed * rate
		currentTokens = math.min(currentTokens + newTokens, capacity)
		lastRefillTime = now
		
		-- 尝试获取令牌
		local allowed = 0
		if currentTokens >= take then
			currentTokens = currentTokens - take
			allowed = 1
		end
		
		-- 更新桶状态
		redis.call('hmset', key, 
			'tokens', currentTokens, 
			'lastRefillTime', lastRefillTime)
		-- 设置过期时间，防止内存泄漏
		redis.call('expire', key, 3600)
		
		return allowed
	`

	now := float64(time.Now().UnixNano()) / 1e9 // 当前时间（秒）
	take := 1                                   // 每次请求消耗1个令牌

	// 执行Lua脚本
	result, err := s.RedisDb.Eval(ctx, luaScript, []string{key},
		rate, capacity, now, take).Int()
	if err != nil {
		return false, fmt.Errorf("令牌桶限流检查失败: %w", err)
	}

	return result == 1, nil
}

// 漏桶算法实现 - 严格控制流出速率
func (s *OrderService) allowWithLeakyBucket(ctx context.Context, userId string, rate int, capacity int) (bool, error) {
	// 每个用户一个独立的漏桶
	key := fmt.Sprintf("leaky_bucket:order:%s", userId)

	// 使用Lua脚本保证操作原子性
	luaScript := `
		local key = KEYS[1]
		local rate = tonumber(ARGV[1])  -- 每秒处理速率
		local capacity = tonumber(ARGV[2])  -- 桶容量
		local now = tonumber(ARGV[3])  -- 当前时间（秒）
		
		-- 获取当前桶状态
		local bucket = redis.call('hgetall', key)
		local water = 0  -- 当前水量
		local lastCheck = now  -- 上次检查时间
		
		-- 如果桶存在，获取当前状态
		if #bucket == 4 then
			water = tonumber(bucket[2])
			lastCheck = tonumber(bucket[4])
		end
		
		-- 计算漏水后的水量
		local elapsed = now - lastCheck
		water = math.max(0, water - elapsed * rate)
		lastCheck = now
		
		-- 尝试加水
		if water + 1 <= capacity then
			water = water + 1
			-- 更新桶状态
			redis.call('hmset', key, 
				'water', water, 
				'lastCheck', lastCheck)
			redis.call('expire', key, 3600)
			return 1  -- 允许
		else
			-- 水满了，拒绝
			return 0  -- 拒绝
		end
	`

	now := float64(time.Now().UnixNano()) / 1e9 // 当前时间（秒）

	// 执行Lua脚本
	result, err := s.RedisDb.Eval(ctx, luaScript, []string{key},
		rate, capacity, now).Int()
	if err != nil {
		return false, fmt.Errorf("漏桶限流检查失败: %w", err)
	}

	return result == 1, nil
}

func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderReply, error) {

	//查询柜子的是否已经寄存
	var order data.LockerStorages
	s.DB.Where("cabinet_id  = ?", req.CabinetId).First(&order)

	if order.CabinetId != (req.CabinetId) {
		return nil, errors.New("CabinetId already exists")
	}

	//随机生成订单编号
	NewString := uuid.NewString()
	virtualDuration := time.Duration(req.ScheduledDuration) * time.Minute

	comletetime := time.Now().Add(virtualDuration)

	var LockerOrders = data.LockerOrders{
		OrderNumber:         NewString,
		UserId:              uint64(req.UserId),
		StartTime:           time.Time{},
		ScheduledDuration:   int32(req.ScheduledDuration),
		Price:               req.Price,
		Discount:            req.Discount,
		AmountPaid:          req.AmountPaid,
		StorageLocationName: req.StorageLocationName,
		CabinetId:           int32(req.CabinetId),
		Status:              int8(req.Status),
		CreateTime:          time.Time{},
		UpdateTime:          comletetime,
		DepositStatus:       int8(req.DepositStatus),
		DeletedAt:           gorm.DeletedAt{},
		LockerPointId:       int32(req.LockerPointId),
		Title:               req.Title,
	}
	err := s.DB.Create(&LockerOrders).Error
	if err != nil {
		return nil, fmt.Errorf("create order failed: %v", err)
	}
	var LockerStorages = data.LockerStorages{
		OrderId:    req.OrderId,
		CabinetId:  req.CabinetId,
		Status:     req.Status,
		CreateTime: time.Time{},
		UpdateTime: comletetime,
		UserId:     req.UserId,
	}
	err = s.DB.Create(&LockerStorages).Error
	if err != nil {
		return nil, fmt.Errorf("create order failed: %v", err)
	}

	//创建柜子的状态处理

	Lockers := data.Lockers{
		LockerPointId: int32(req.LockerPointId),
		TypeId:        int32(req.TypeId),
		Status:        int8(req.Status),
	}
	err = s.DB.Create(&Lockers).Error
	if err != nil {
		return nil, fmt.Errorf("create order failed: %v", err)
	}

	return &pb.CreateOrderReply{
		Msg: "订单发送,并记录",
	}, nil

}

func (s *OrderService) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.UpdateOrderReply, error) {
	var payUrl string

	// 1. 启动数据库事务
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// 2. 在事务中查询订单是否存在
		var lockerOrder data.LockerOrders
		if err := tx.Where("id = ?", req.Id).First(&lockerOrder).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.NotFound, "订单不存在")
			}
			return status.Errorf(codes.Internal, "数据库错误: %v", err)
		}

		// 如果订单已经支付或关闭，则不允许再次支付
		if lockerOrder.Status != 1 { // 1-待支付
			pkg.LogError("订单状态异常，无法支付")
			return status.Errorf(codes.FailedPrecondition, "订单状态异常，无法支付")
		}

		// 3. 查询柜子信息，获取柜子类型ID
		var locker data.Lockers
		if err := tx.Where("type_id = ?", req.TypeId).Find(&locker).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.NotFound, "柜子信息不存在")
			}
			return status.Errorf(codes.Internal, "数据库错误: %v", err)
		}

		// 4. 根据网点ID和柜子类型查询计费规则
		var pricingRule data.LockerPricingRules
		if err := tx.Where("locker_type = ?", req.LockerType).First(&pricingRule).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.NotFound, "该类型的储物柜没有找到有效的计费规则")
			}
			return status.Errorf(codes.Internal, "数据库错误: %v", err)
		}

		totalPrice := pricingRule.HourlyRate * float64(lockerOrder.ScheduledDuration)

		// 6. 更新订单信息（只更新必要的字段）
		updateData := map[string]interface{}{
			"actual_duration": req.ActualDuration,
			"price":           totalPrice,
			"amount_paid":     totalPrice,
			"status":          req.Status,
			"deposit_status":  req.DepositStatus,
		}

		if err := tx.Model(&lockerOrder).Updates(updateData).Error; err != nil {
			return status.Errorf(codes.Internal, "更新订单失败: %v", err)
		}

		// 7. 生成支付链接（使用计算后的金额）
		totalAmountStr := strconv.FormatFloat(totalPrice, 'f', 2, 64)
		payUrl = pkg.Pay(lockerOrder.Title, lockerOrder.OrderNumber, totalAmountStr)

		// 事务会自动提交
		return nil
	})

	if err != nil {
		return nil, err // 如果事务过程中有任何错误，则直接返回错误
	}

	// 8. 返回响应
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

	// 根据存储位置名称过滤
	if req.StorageLocationName != "" {
		db = db.Where("storage_location_name LIKE ?", "%"+req.StorageLocationName+"%")
	}

	// 根据订单状态过滤
	if len(req.Status) > 0 {
		// 如果传入了状态列表，则按状态过滤
		db = db.Where("status IN (?)", req.Status)
	}

	// 添加排序
	db = db.Order("id DESC")

	// 获取总条数
	var total int64
	db.Count(&total)

	// 执行查询
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
			Status:              int64((v.Status)),
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
		UserId:              int64((order.UserId)),
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

func (s *OrderService) HandleRemindTask(ctx context.Context, req *pb.HandleRemindTaskRequest) (*pb.HandleRemindTaskReply, error) {

	// 1. 参数校验
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "任务ID不能为空")
	}

	// 2. 查询定时任务详情（验证任务合法性）
	var task data.TimerTasks
	if err := s.DB.Where("id = ? AND task_type = ? AND status = ?",
		req.Id, "remind", req.Status).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "到期提醒任务不存在或状态异常")
		}
		log.Errorw(ctx, "查询定时任务失败", "task_id", req.Id, "err", err)
		return nil, status.Errorf(codes.Internal, "查询任务信息失败")
	}

	// 3. 查询关联的寄存订单
	var storage data.LockerStorages
	if err := s.DB.Where("id = ?", task.LockerStorageId).First(&storage).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "关联的寄存订单不存在")
		}
		log.Errorw(ctx, "查询寄存订单失败", "storage_id", task.LockerStorageId, "err", err)
		return nil, status.Errorf(codes.Internal, "查询订单信息失败")
	}
	//4.验证订单的状态
	// 4. 验证订单状态（仅对"寄存中"的订单发送提醒）
	if storage.Status != 1 { // 1=寄存中
		log.Infow(ctx, "订单状态无需发送提醒",
			"storage_id", storage.Id,
			"order_id", storage.OrderId,
			"status", storage.Status)

		// 更新任务状态为"已完成"（无需执行）
		if err := s.DB.Model(&data.TimerTasks{}).
			Where("id = ?", req.Id).
			Updates(map[string]interface{}{
				"status":     2, // 2=已完成
				"updated_at": time.Now(),
			}).Error; err != nil {
			log.Warnw(ctx, "更新任务状态失败", "task_id", req.Id, "err", err)
		}
		return &pb.HandleRemindTaskReply{Msg: "订单状态无需提醒，任务已结束"}, nil
	}
	// 5. 查询用户信息（获取通知渠道）
	var user data.Users
	if err := s.DB.Where("id = ?", storage.UserId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "关联用户不存在")
		}
		log.Errorw(ctx, "查询用户信息失败", "user_id", storage.UserId, "err", err)
		return nil, status.Errorf(codes.Internal, "查询用户信息失败")
	}
	// 6. 发送到期提醒（多渠道通知）
	remindContent := fmt.Sprintf(
		"您在【%s】的寄存订单（单号：%d）到期，请及时取件或续存。",
		storage.OrderId,
		storage.Id,
	)

	if user.Mobile != "" {

		smsErr := s.RedisDb.Set(context.Background(), remindContent, user.Mobile, time.Hour)
		if smsErr != nil {
			log.Warnw(ctx, "短信提醒发送失败",
				"user_id", user.Id,
				"phone", user.Mobile,
				"err", smsErr)
		} else {
			log.Infow(ctx, "短信提醒发送成功", "user_id", user.Id, "phone", user.Mobile)
		}
	}

	return &pb.HandleRemindTaskReply{
		Msg: "定时任务提示",
	}, nil
}

// 对于超时寄存柜订单的的处理
func (s *OrderService) HandleTimeOutTask(ctx context.Context, req *pb.HandleTimeOutTaskRequest) (*pb.HandleTimeOutTaskReply, error) {

	// 1. 参数校验
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "任务ID不能为空")
	}

	// 2. 查询定时任务详情（验证任务合法性）
	var task data.TimerTasks
	err := s.DB.Where("id = ? AND task_type = ? AND status = ? ", req.Id, "handle", req.Status).First(&task).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "到期提醒任务不存在或状态异常")
		}
		log.Errorw(ctx, "查询定时任务失败", "task_id", req.Id, "err", err)
		return nil, status.Errorf(codes.Internal, "查询任务信息失败")
	}
	var order data.LockerOrders
	err = s.DB.Where("order_id = ? ", req.OrderId).First(&task).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "到期提醒任务不存在或状态异常")
		}
		log.Errorw(ctx, "查询定时任务失败", "task_id", req.Id, "err", err)
		return nil, status.Errorf(codes.Internal, "查询任务信息失败")
	}

	//3.根据当前时间比较订单的时间判断是否超时
	//首先验证任务时间是否超时
	//验证订单的存储时间是否超时

	now := time.Now()
	if !task.ExecuteTime.Before(now) {
		return nil, status.Errorf(codes.FailedPrecondition, "任务未到执行时间，当前时间：%s，执行时间：%s",
			now.Format(time.RFC3339), task.ExecuteTime.Format(time.RFC3339))
	}
	if !order.UpdateTime.Before(now) {
		return nil, status.Errorf(codes.FailedPrecondition, "订单未超时，当前时间：%s，订单过期时间：%s",
			now.Format(time.RFC3339), order.UpdateTime.Format(time.RFC3339))
	}

	Time := int64(order.ScheduledDuration) + req.ScheduledDuration

	//4.超时加钱修改订单表的价钱和实际的时间
	Total := Time * req.HourlyRate
	orders := data.LockerOrders{
		ActualDuration: int32(Time),
		AmountPaid:     float64(Total),
		Status:         int8((req.Status)), // 修改订单状态为"已超时"
	}

	err = s.DB.Where("id = ?", req.Id).Updates(&orders).Error // 更新订单状态为"已超时"

	return &pb.HandleTimeOutTaskReply{
		Msg: "对于超时寄存柜订单的的处理完成",
	}, nil

}

func (s *OrderService) ManageOrderSearch(ctx context.Context, req *pb.ManageOrderSearchRequest) (*pb.ManageOrderSearchReply, error) {

	//查询用户是否是管理员
	var admin data.Admin
	err := s.DB.Where("id = ? AND role = ?", req.Id, req.Role).Find(&admin).Error
	if err != nil {
		return nil, fmt.Errorf("管理员端报错")
	}

	if admin.Status != 1 && admin.Role != 1 {
		return nil, fmt.Errorf("管理员端报错")
	}

	//管理员对于订单的查询
	var search data.LockerOrders
	s.DB.Where("order_number  LIKE ? and status  LIKE  ?", "%"+req.OrderNumber+"%", req.Status).Find(&search)
	// 如果订单已经支付或关闭，则不允许再次支付
	if search.Status != 1 { // 1-待支付
		pkg.LogError("订单状态异常，无法支付")
	}

	searchr := &pb.OrderInfo{
		Id:                  search.Id,
		OrderNumber:         search.OrderNumber,
		UserId:              int64((search.UserId)),
		ScheduledDuration:   int64(search.ScheduledDuration),
		ActualDuration:      int64(search.ActualDuration),
		Price:               search.Price,
		Discount:            search.Discount,
		AmountPaid:          search.AmountPaid,
		StorageLocationName: search.StorageLocationName,
		CabinetId:           int64(search.CabinetId),
		Status:              int64(search.Status),
		DepositStatus:       int64(search.DepositStatus),
	}

	return &pb.ManageOrderSearchReply{
		Order: []*pb.OrderInfo{searchr},
		Msg:   "管理员查询订单",
	}, nil
}

func (s *OrderService) ManageOrderDel(ctx context.Context, req *pb.ManageOrderDelRequest) (*pb.ManageOrderDelReply, error) {

	//查询用户是否是管理员
	var admin data.Admin
	err := s.DB.Where("status  = ? AND role = ?", req.Status, req.Role).Find(&admin).Error
	if err != nil {
		return nil, fmt.Errorf("管理员端报错")
	}

	if admin.Status != 1 && admin.Role != 1 {
		return nil, fmt.Errorf("管理员端报错")
	}

	var order data.LockerOrders
	err = s.DB.Where("id = ? ", req.Id).Delete(&order).Error

	return &pb.ManageOrderDelReply{Msg: "删除成功"}, nil
}

func (s *OrderService) ManageOrderDetail(ctx context.Context, req *pb.ManageOrderDetailRequest) (*pb.ManageOrderDetailReply, error) {
	var admin data.Admin
	err := s.DB.Where("status  = ? AND role = ?", req.Status, req.Role).Find(&admin).Error
	if err != nil {
		return nil, fmt.Errorf("管理员端报错")
	}

	if admin.Status != 1 && admin.Role != 1 {
		return nil, fmt.Errorf("管理员端报错")
	}

	var order data.LockerOrders

	fmt.Println("收到的订单ID:", req.Id) // 打印实际收到的ID
	if req.Id <= 0 {
		return &pb.ManageOrderDetailReply{
			Msg: "无效的订单ID",
		}, nil
	}

	// 执行查询
	result := s.DB.Where("id = ?", req.Id).First(&order)

	// 判断是否查询到数据
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &pb.ManageOrderDetailReply{
				Msg: "订单不存在",
			}, nil
		}
		// 数据库错误
		return &pb.ManageOrderDetailReply{
			Msg: "数据库错误：" + result.Error.Error(),
		}, nil
	}

	orderInfo := &pb.OrderInfo{
		Id:                  order.Id,
		OrderNumber:         order.OrderNumber,
		UserId:              int64((order.UserId)),
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
	return &pb.ManageOrderDetailReply{
		Msg:   "success",
		Order: []*pb.OrderInfo{orderInfo},
	}, nil
}

package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"ito-deposit/internal/conf"

	"gorm.io/gorm"
	"math/rand"
	"time"

	"ito-deposit/internal/data"

	pb "ito-deposit/api/helloworld/v1"
)

type UserService struct {
	pb.UnimplementedUserServer
	RedisDb *redis.Client
	DB      *gorm.DB
	server  *conf.Server
}

func NewUserService(datas *data.Data, server *conf.Server) *UserService {
	return &UserService{
		RedisDb: datas.Redis,
		DB:      datas.DB,
		server:  server,
	}
}

func (s *UserService) SendSms(ctx context.Context, req *pb.SendSmsRequest) (*pb.SendSmsRes, error) {
	code := rand.Intn(9000) + 1000
	fmt.Printf("[SendSms] raw req: %+v", req) // 打印整个结构体
	fmt.Printf("[SendSms] mobile=%q source=%q", req.Mobile, req.Source)
	s.RedisDb.Set(context.Background(), "sendSms"+req.Mobile+req.Source, code, time.Minute*5)
	return &pb.SendSmsRes{
		Code: 200,
		Msg:  "短信发送成功",
	}, nil
}
func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterRes, error) {
	get := s.RedisDb.Get(context.Background(), "sendSms"+req.Mobile+"register")
	if get.Val() != req.SmsCode {
		return &pb.RegisterRes{
			Code: 500,
			Msg:  "验证码错误",
		}, nil
	}
	user := data.Users{
		Username: req.Username,
		Mobile:   req.Mobile,
		Password: req.Password,
	}
	err := s.DB.Debug().Create(&user).Error
	if err != nil {
		return &pb.RegisterRes{
			Code: 500,
			Msg:  "注册失败",
		}, nil
	}
	return &pb.RegisterRes{
		Code: 200,
		Msg:  "注册成功",
	}, nil
}
func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginRes, error) {
	get := s.RedisDb.Get(context.Background(), "sendSms"+req.Mobile+"login")
	if get.Val() != req.SmsCode {
		return &pb.LoginRes{
			Code: 500,
			Msg:  "验证码错误",
		}, nil
	}
	var user data.Users
	err := s.DB.Debug().Where("mobile = ?", req.Mobile).Find(&user).Error
	if err != nil {
		return &pb.LoginRes{
			Code: 500,
			Msg:  "查询失败",
		}, nil
	}
	if user.Id == 0 {
		return &pb.LoginRes{
			Code: 500,
			Msg:  "用户不存在",
		}, nil
	}
	if req.Password != user.Password {
		return &pb.LoginRes{
			Code: 500,
			Msg:  "密码错误",
		}, nil
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// 根据您的需求设置 JWT 中的声明
		"your_custom_claim": "your_custom_value",
		"id":                "123",
	})

	signedString, err := claims.SignedString([]byte(s.server.Jwt.Authkey))
	if err != nil {
		return nil, err
	}
	return &pb.LoginRes{
		Code:  200,
		Msg:   "登陆成功",
		Id:    user.Id,
		Token: signedString,
	}, nil
}
func (s *UserService) OrderList(ctx context.Context, req *pb.OrderListRequest) (*pb.OrderListRes, error) {
	var order []data.LockerOrders
	err := s.DB.Debug().Find(&order).Error
	if err != nil {
		return &pb.OrderListRes{
			Code: 500,
			Msg:  "查询失败",
		}, nil
	}
	var lists []*pb.OrderList
	for _, o := range order {
		list := pb.OrderList{
			OrderNumber:         o.OrderNumber,
			UserId:              int64(o.UserId),
			ScheduledDuration:   int64(o.ScheduledDuration),
			ActualDuration:      int64(o.ActualDuration),
			Price:               float32(o.Price),
			Discount:            float32(o.Discount),
			AmountPaid:          float32(o.AmountPaid),
			StorageLocationName: o.StorageLocationName,
			CabinetId:           int64(o.CabinetId),
			Status:              1,
			DepositStatus:       1,
		}
		lists = append(lists, &list)
	}
	return &pb.OrderListRes{
		Code: 200,
		Msg:  "查询成功",
		List: lists,
	}, nil
}
func (s *UserService) Admin(ctx context.Context, req *pb.AdminRequest) (*pb.AdminRes, error) {

	// 1. 查询网点总数
	var pointCount int64
	if err := s.DB.Model(&data.LockerPoint{}).
		Where("admin_id = ?", req.AdminId).    // 根据管理员ID过滤
		Count(&pointCount).Error; err != nil { // 获取总数
		return &pb.AdminRes{
			Code: 500,
			Msg:  "查询网点数量失败: " + err.Error(),
		}, nil
	}

	// 2. 计算时间范围
	todayStart := time.Now().Format("2006-01-02 00:00:00")                                                                             // 今日开始时间
	yesterdayStart := time.Now().AddDate(0, 0, -1).Format("2006-01-02 00:00:00")                                                       // 昨日开始时间
	yesterdayEnd := time.Now().Format("2006-01-02 00:00:00")                                                                           // 昨日结束时间
	monthStart := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location()).Format("2006-01-02 00:00:00") // 本月开始时间

	// 3. 查询今日订单数
	var todayOrderCount int64
	if err := s.DB.Model(&data.LockerOrders{}).
		Where("locker_point_id = ?", req.LockerPointId). // 根据寄存点ID过滤
		Where("create_time > ?", todayStart).            // 创建时间大于今日开始时间
		Count(&todayOrderCount).Error; err != nil {      // 获取订单数量
		return &pb.AdminRes{
			Code: 500,
			Msg:  "查询今日订单数失败: " + err.Error(),
		}, nil
	}

	// 4. 查询昨日订单数和收益
	var yesterdayResult struct {
		OrderCount int64
		TotalPaid  float64
	}
	if err := s.DB.Debug().Model(&data.LockerOrders{}).
		Select("COUNT(1) as order_count, SUM(amount_paid) as total_paid").  // 计算订单数和总收益
		Where("locker_point_id = ?", req.LockerPointId).                    // 根据寄存点ID过滤
		Where("create_time BETWEEN ? AND ?", yesterdayStart, yesterdayEnd). // 创建时间在昨日范围内
		Scan(&yesterdayResult).Error; err != nil {                          // 扫描结果到结构体
		return &pb.AdminRes{
			Code: 500,
			Msg:  "查询昨日订单信息失败: " + err.Error(),
		}, nil
	}
	fmt.Println(yesterdayStart, yesterdayEnd, yesterdayResult)
	// 5. 查询本月订单数和收益
	var monthResult struct {
		OrderCount int64
		TotalPaid  float64
	}
	if err := s.DB.Model(&data.LockerOrders{}).
		Select("COUNT(1) as order_count, SUM(amount_paid) as total_paid"). // 计算订单数和总收益
		Where("locker_point_id = ?", req.LockerPointId).                   // 根据寄存点ID过滤
		Where("create_time > ?", monthStart).                              // 创建时间大于本月开始时间
		Scan(&monthResult).Error; err != nil {                             // 扫描结果到结构体
		return &pb.AdminRes{
			Code: 500,
			Msg:  "查询本月订单信息失败: " + err.Error(),
		}, nil
	}

	// 6. 返回结果
	return &pb.AdminRes{
		Code:              200,                                // 成功状态码
		Msg:               "查询成功",                             // 成功消息
		PointNum:          pointCount,                         // 网点总数
		LastOrderNum:      todayOrderCount,                    // 今日订单数
		YesterdayOrderNum: yesterdayResult.OrderCount,         // 昨日订单数
		LastOrderPrice:    float32(yesterdayResult.TotalPaid), // 昨日收益
		MouthPrice:        float32(monthResult.TotalPaid),     // 本月收益
		MonthNum:          monthResult.OrderCount,             // 本月订单数
	}, nil
}

package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"ito-deposit/internal/basic/pkg"
	"math/rand"
	"strconv"
	"time"

	jwt1 "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/golang-jwt/jwt/v5"
	pb "ito-deposit/api/helloworld/v1"
	"ito-deposit/internal/conf"
	"ito-deposit/internal/data"
)

type DepositService struct {
	pb.UnimplementedDepositServer
	data   *data.Data
	server *conf.Server
	pkg    *pkg.SendSms
}

func NewDepositService(data2 *data.Data, server *conf.Server, pkg *pkg.SendSms) *DepositService {
	return &DepositService{
		data:   data2,
		server: server,
		pkg:    pkg,
	}
}

func (s *DepositService) CreateDeposit(ctx context.Context, req *pb.CreateDepositRequest) (*pb.CreateDepositReply, error) {
	kratosToken, ok := jwt1.FromContext(ctx)
	if !ok {
		return &pb.CreateDepositReply{
			Code: 401,
			Msg:  "token不正确或者未传",
		}, nil
	}

	mapClaims, ok := kratosToken.(*jwt.MapClaims)

	userId := (*mapClaims)["id"].(string)
	var lockerPriceRules data.LockerPricingRules
	err := s.data.DB.Table("locker_pricing_rules").Where("network_id = ? ", req.CabinetId).Where("status = 1").
		Where("locker_type = ?", req.LockerType).Limit(1).Find(&lockerPriceRules).Error
	if err != nil {
		return nil, err
	}
	intUserId, err := strconv.ParseInt(userId, 10, 64)
	OrderNo := time.Now().Format("20060102150405") + userId
	var price float64
	if lockerPriceRules.IsDepositEnabled == 1 {
		price += lockerPriceRules.DepositAmount
	}
	price += lockerPriceRules.HourlyRate * float64(req.ScheduledDuration)
	var locker data.Lockers
	err = s.data.DB.Transaction(func(tx *gorm.DB) error {
		err = s.data.DB.Table("lockers").Where("locker_point_id = ?", req.CabinetId).Select("id").Where("status = 1").Limit(1).Find(&locker).Error
		if err != nil {
			return err
		}
		if locker.Id == 0 {
			return errors.New("寄存柜可用数量不足")
		}
		updaateres := s.data.DB.Table("lockers").Where("id = ?", locker.Id).Update("status", 2)
		if updaateres.RowsAffected == 0 {
			return errors.New("更新柜子状态失败")
		}
		if updaateres.Error != nil {
			return updaateres.Error
		}
		//cabineId := strconv.Itoa(int(locker.Id))
		addOrder := data.LockerOrders{
			OrderNumber:       OrderNo,
			UserId:            uint64(intUserId),
			StartTime:         time.Now(),
			ScheduledDuration: req.ScheduledDuration,
			Price:             price,
			Discount:          0,
			AmountPaid:        price,
			Status:            1,
			CabinetId:         locker.Id,
			CreateTime:        time.Now(),
			UpdateTime:        time.Now(),
			DepositStatus:     0,
			ActualDuration:    0,
		}
		err = s.data.DB.Table("locker_orders").Create(&addOrder).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateDepositReply{
		Code: 200,
		Msg:  "添加寄存订单成功",
		Data: &pb.DepositReplyData{
			OrderNo:  OrderNo,
			LockerId: locker.Id,
		},
	}, nil
}
func (s *DepositService) UpdateDeposit(ctx context.Context, req *pb.UpdateDepositRequest) (*pb.UpdateDepositReply, error) {
	return &pb.UpdateDepositReply{}, nil
}
func (s *DepositService) DeleteDeposit(ctx context.Context, req *pb.DeleteDepositRequest) (*pb.DeleteDepositReply, error) {
	return &pb.DeleteDepositReply{}, nil
}
func (s *DepositService) GetDeposit(ctx context.Context, req *pb.GetDepositRequest) (*pb.GetDepositReply, error) {
	return &pb.GetDepositReply{}, nil
}
func (s *DepositService) ListDeposit(ctx context.Context, req *pb.ListDepositRequest) (*pb.ListDepositReply, error) {
	return &pb.ListDepositReply{}, nil
}

func (s *DepositService) ReturnToken(ctx context.Context, req *pb.ReturnTokenReq) (*pb.ReturnTokenRes, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// 根据您的需求设置 JWT 中的声明
		"your_custom_claim": "your_custom_value",
		"id":                "123",
		"exp":               time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	signedString, err := claims.SignedString([]byte(s.server.Jwt.Authkey))
	if err != nil {
		return nil, err
	}
	return &pb.ReturnTokenRes{
		Token: signedString,
		Coe:   200,
		Msg:   "生成token成功",
	}, nil
}

func (s *DepositService) DecodeToken(ctx context.Context, req *pb.ReturnTokenReq) (*pb.ReturnTokenRes, error) {
	// 1. 从上下文获取 Kratos 包装的 Token
	kratosToken, ok := jwt1.FromContext(ctx)
	fmt.Println("kratosToken", kratosToken)
	fmt.Printf("kratosToken 类型: %T\n", kratosToken)
	if !ok {
		return &pb.ReturnTokenRes{
			Coe: 401,
			Msg: "未找到有效的 JWT Token（可能未登录或 Token 无效）",
		}, nil
	}

	mapClaims, ok := kratosToken.(*jwt.MapClaims)
	fmt.Println("mapClaims", mapClaims)

	return &pb.ReturnTokenRes{
		Token: (*mapClaims)["id"].(string),
		Coe:   200,
		Msg:   "token内容 ",
	}, nil
}

func (s *DepositService) GetDepositLocker(ctx context.Context, req *pb.GetDepositLockerReq) (*pb.GetDepositLockerRes, error) {
	pointID := req.LockerId

	// 1. 网点信息（主键查）
	var point data.LockerPoint
	err := s.data.DB.Where("id = ? ", pointID).Limit(1).Find(&point).Error
	if err != nil {
		return nil, err
	}
	if point.Id == 0 {
		return &pb.GetDepositLockerRes{}, errors.New("网点信息不存在")
	}

	// 2. 全部柜型（小表缓存）
	var types []*data.LockerType
	err = s.data.DB.Find(&types).Error
	if err != nil {
		return nil, err
	}

	// 3. 实时库存（单表聚合）
	type result struct {
		TypeID int32 `gorm:"column:type_id"`
		Num    int32 `gorm:"column:num"`
	}
	var list []result
	err = s.data.DB.Table("lockers").Where("locker_point_id = ?", pointID).Where("status = 1").
		Select("type_id,count(1) as num").Group("type_id").Find(&list).Error

	if err != nil {
		return nil, err
	}
	listMap := make(map[int32]int32)
	for _, v := range list {
		listMap[v.TypeID] = v.Num
	}

	// 4. 组装返回值
	res := &pb.GetDepositLockerRes{
		Address:   point.Address,
		Name:      point.Name,
		Longitude: float32(point.Longitude),
		Latitude:  float32(point.Latitude),
	}

	var ids []int32
	for _, t := range types {
		cnt := listMap[t.Id]
		if cnt == 0 {
			continue // 前端不展示已满柜型
		}
		ids = append(ids, t.Id)
	}
	var priceRule []*data.LockerPricingRules
	err = s.data.DB.Table("locker_pricing_rules").Where("network_id in (?)", ids).Where("status = 1").Find(&priceRule).Error
	priceRuleMap := make(map[int64]*data.LockerPricingRules)
	for _, v := range priceRule {
		priceRuleMap[v.Id] = v
	}
	for _, t := range types {
		freeDuration := float64(0)
		hourlyRate := float64(0)
		if priceRuleMap[int64(t.Id)] != nil {
			freeDuration = priceRuleMap[int64(t.Id)].FreeDuration
			hourlyRate = priceRuleMap[int64(t.Id)].HourlyRate
		}
		res.Locker = append(res.Locker, &pb.Locker{
			Name:         t.Name,
			Description:  t.Description,
			Size:         t.Size,
			Num:          listMap[t.Id],
			HourlyRate:   float32(hourlyRate),
			FreeDuration: float32(freeDuration),
			LockerType:   t.Id,
		})
	}

	return res, nil
}

func (s *DepositService) UpdateDepositLockerId(ctx context.Context, req *pb.UpdateDepositLockerIdReq) (*pb.UpdateDepositLockerIdRes, error) {
	var order data.LockerOrders
	err := s.data.DB.Where("order_number = ?", req.OrderId).Limit(1).Find(&order).Error
	if err != nil {
		return nil, err
	}
	var locker data.Lockers
	err = s.data.DB.Table("lockers").Where("id = ?", order.CabinetId).Limit(1).Find(&locker).Error
	if err != nil {
		return nil, err
	}
	if locker.Id == 0 {
		return nil, errors.New("订单关联寄存柜不存在")
	}
	var newLocker data.Lockers
	err = s.data.DB.Table("lockers").Where("locker_point_id = ?", locker.LockerPointId).Where("type_id = ?", locker.TypeId).Where("status = 1").Limit(1).Find(&newLocker).Error
	if err != nil {
		return nil, err
	}
	if newLocker.Id == 0 {
		return nil, errors.New("无可用寄存柜")
	}

	updateLocker := s.data.DB.Table("lockers").Where("id = ?", newLocker.Id).Update("status", 2).Debug()

	if updateLocker.Error != nil {
		return nil, updateLocker.Error
	}
	res := s.data.DB.Table("locker_orders").Where("order_number = ?", req.OrderId).Update("cabinet_id", newLocker.Id)
	if res.Error != nil {
		return nil, res.Error
	}
	updateStatus := 1
	redisKey := fmt.Sprintf("locker:%d", locker.Id)
	lockerNum, err := s.data.Redis.Get(ctx, redisKey).Result()
	if err != nil && err != redis.Nil {
		s.data.Redis.Set(ctx, redisKey, 0, 86400*time.Second)
	} else {
		intLockerNum, _ := strconv.Atoi(lockerNum)
		if intLockerNum >= 3 {
			updateStatus = 4
		}
	}
	s.data.DB.Table("lockers").Where("id = ?", order.CabinetId).Update("status", updateStatus)
	s.data.Redis.Incr(ctx, redisKey)
	return &pb.UpdateDepositLockerIdRes{
		Code:     200,
		Msg:      "修改成功",
		LockerId: newLocker.Id,
	}, nil
}

func (s *DepositService) SendCodeByOrder(ctx context.Context, req *pb.SendCodeByOrderReq) (*pb.SendCodeByOrderRes, error) {
	var lockerOrder data.LockerOrders
	err := s.data.DB.Table("locker_orders").Where("order_number = ?", req.OrderNo).Limit(1).Find(&lockerOrder).Error
	if err != nil {
		return nil, err
	}
	var phone string
	err = s.data.DB.Table("users").Where("id = ?", lockerOrder.UserId).Pluck("mobile", &phone).Error
	if err != nil {
		return nil, err
	}
	code := rand.Intn(900000) + 100000
	codeRes := s.pkg.SendSms(phone, code)
	if !codeRes {
		return nil, errors.New("验证码发送失败")
	}
	redisKey := fmt.Sprintf("point:%d", lockerOrder.LockerPointId)
	err = s.data.Redis.Set(ctx, redisKey, code, 300*time.Second).Err()
	if err != nil {
		return nil, err
	}

	return &pb.SendCodeByOrderRes{
		Msg:  "寄存短信发送成功",
		Code: 200,
		Data: "",
	}, nil
}

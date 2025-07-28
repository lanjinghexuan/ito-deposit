package service

import (
	"context"
	"errors"
	"fmt"
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
}

func NewDepositService(data2 *data.Data, server *conf.Server) *DepositService {
	return &DepositService{
		data:   data2,
		server: server,
	}
}

func (s *DepositService) CreateDeposit(ctx context.Context, req *pb.CreateDepositRequest) (*pb.CreateDepositReply, error) {
	return &pb.CreateDepositReply{}, nil
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
	var point data.Location
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
	err = s.data.DB.Table("locker").Where("locker_point_id = ?", pointID).Where("status = 1").
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

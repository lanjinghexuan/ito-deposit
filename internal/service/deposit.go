package service

import (
	"context"
	"fmt"

	jwt1 "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/golang-jwt/jwt/v5"
	pb "ito-deposit/api/helloworld/v1"
	"ito-deposit/internal/conf"
	"ito-deposit/internal/data"
)

type DepositService struct {
	pb.UnimplementedDepositServer
	aaa    *data.Data
	server *conf.Server
}

func NewDepositService(data2 *data.Data, server *conf.Server) *DepositService {
	return &DepositService{
		aaa:    data2,
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

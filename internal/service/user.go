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

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	return &pb.CreateUserReply{}, nil
}
func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserReply, error) {
	return &pb.UpdateUserReply{}, nil
}
func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserReply, error) {
	return &pb.DeleteUserReply{}, nil
}
func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	return &pb.GetUserReply{}, nil
}
func (s *UserService) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserReply, error) {
	return &pb.ListUserReply{}, nil
}
func (s *UserService) SendSms(ctx context.Context, req *pb.SendSmsRequest) (*pb.SendSmsReply, error) {
	code := rand.Intn(9000) + 1000
	fmt.Printf("[SendSms] raw req: %+v", req) // 打印整个结构体
	fmt.Printf("[SendSms] mobile=%q source=%q", req.Mobile, req.Source)
	s.RedisDb.Set(context.Background(), "sendSms"+req.Mobile+req.Source, code, time.Minute*5)
	return &pb.SendSmsReply{
		Code: 200,
		Msg:  "短信发送成功",
	}, nil
}
func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	get := s.RedisDb.Get(context.Background(), "sendSms"+req.Mobile+"register")
	if get.Val() != req.SmsCode {
		return &pb.RegisterReply{
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
		return &pb.RegisterReply{
			Code: 500,
			Msg:  "注册失败",
		}, nil
	}
	return &pb.RegisterReply{
		Code: 200,
		Msg:  "注册成功",
	}, nil
}
func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	get := s.RedisDb.Get(context.Background(), "sendSms"+req.Mobile+"login")
	if get.Val() != req.SmsCode {
		return &pb.LoginReply{
			Code: 500,
			Msg:  "验证码错误",
		}, nil
	}
	var user data.Users
	err := s.DB.Debug().Where("mobile = ?", req.Mobile).Find(&user).Error
	if err != nil {
		return &pb.LoginReply{
			Code: 500,
			Msg:  "查询失败",
		}, nil
	}
	if user.Id == 0 {
		return &pb.LoginReply{
			Code: 500,
			Msg:  "用户不存在",
		}, nil
	}
	if req.Password != user.Password {
		return &pb.LoginReply{
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
	return &pb.LoginReply{
		Code:  200,
		Msg:   "登陆成功",
		Id:    user.Id,
		Token: signedString,
	}, nil
}

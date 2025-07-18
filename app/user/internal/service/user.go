package service

import (
	"context"
	pb "ito/api/user/v1"
	model "ito/app/user/internal/data"
	"ito/internal/data"
)

type UserService struct {
	pb.UnimplementedUserServer
	DB *data.Data
}

func NewUserService(data *data.Data) *UserService {
	return &UserService{
		UnimplementedUserServer: pb.UnimplementedUserServer{},
		DB:                      data,
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
	var userData []*model.Users
	s.DB.DB.Table("users").Find(&userData)
	var dataList []*pb.UserList
	for _, v := range userData {
		dataList = append(dataList, &pb.UserList{
			Id:     v.Id,
			Name:   v.Name,
			Mobile: v.Mobile,
		})
	}
	return &pb.ListUserReply{
		Code: 200,
		Data: dataList,
		Msg:  "查询成功",
	}, nil
}

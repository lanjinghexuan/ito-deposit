package service

import (
	"context"

	pb "ito-deposit/api/helloworld/v1"
)

type HomeService struct {
	pb.UnimplementedHomeServer
}

func NewHomeService() *HomeService {
	return &HomeService{}
}

func (s *HomeService) CreateHome(ctx context.Context, req *pb.CreateHomeRequest) (*pb.CreateHomeReply, error) {
	return &pb.CreateHomeReply{}, nil
}
func (s *HomeService) UpdateHome(ctx context.Context, req *pb.UpdateHomeRequest) (*pb.UpdateHomeReply, error) {
	return &pb.UpdateHomeReply{}, nil
}
func (s *HomeService) DeleteHome(ctx context.Context, req *pb.DeleteHomeRequest) (*pb.DeleteHomeReply, error) {
	return &pb.DeleteHomeReply{}, nil
}
func (s *HomeService) GetHome(ctx context.Context, req *pb.GetHomeRequest) (*pb.GetHomeReply, error) {
	return &pb.GetHomeReply{}, nil
}
func (s *HomeService) ListHome(ctx context.Context, req *pb.ListHomeRequest) (*pb.ListHomeReply, error) {
	return &pb.ListHomeReply{}, nil
}

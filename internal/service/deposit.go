package service

import (
	"context"

	pb "ito-deposit/api/helloworld/v1"
)

type DepositService struct {
	pb.UnimplementedDepositServer
}

func NewDepositService() *DepositService {
	return &DepositService{}
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

package service

import (
	"context"

	pb "ito-deposit/api/helloworld/v1"
)

type OrderService struct {
	pb.UnimplementedOrderServer
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderReply, error) {
	return &pb.CreateOrderReply{}, nil
}
func (s *OrderService) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.UpdateOrderReply, error) {
	return &pb.UpdateOrderReply{}, nil
}
func (s *OrderService) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderReply, error) {
	return &pb.DeleteOrderReply{}, nil
}
func (s *OrderService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderReply, error) {
	return &pb.GetOrderReply{}, nil
}
func (s *OrderService) ListOrder(ctx context.Context, req *pb.ListOrderRequest) (*pb.ListOrderReply, error) {
	return &pb.ListOrderReply{}, nil
}

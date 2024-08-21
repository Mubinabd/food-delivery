package service

import (
	"context"

	pb "gitlab.com/bahodirova/product/genproto/product"
	"gitlab.com/bahodirova/product/storage"
)

type OrderService struct {
	storage storage.StorageI
	pb.UnimplementedOrderServiceServer
}

func NewOrderService(storage storage.StorageI) *OrderService {
	return &OrderService{
		storage: storage,
	}
}

func (a *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderReq) (*pb.Empty, error) {
	return a.storage.Order().CreateOrder(req)
}
func (a *OrderService) GetOrder(ctx context.Context, req *pb.GetById) (*pb.Order, error) {
	return a.storage.Order().GetOrder(req)
}
func (a *OrderService) GetAllOrders(ctx context.Context, req *pb.GetAllOrdersReq) (*pb.GetAllOrderRes, error) {
	return a.storage.Order().GetAllOrders(req)
}
func (a *OrderService) DeleteOrder(ctx context.Context, req *pb.GetById) (*pb.DeleteOrderRes, error) {
	return a.storage.Order().DeleteOrder(req)
}
func (a *OrderService) UpdateOrder(ctx context.Context, req *pb.UpdateOrderReq) (*pb.UpdateOrderRes, error) {
	return a.storage.Order().UpdateOrder(req)
}
func(a *OrderService)PaidOrder(ctx context.Context, req *pb.PaidReq) (*pb.PaidRes, error){
	return a.storage.Order().PaidOrder(req)
}
func(a *OrderService)HistoryOrder(ctx context.Context, req *pb.GetCourierOrderHistoryRequest) (*pb.GetCourierOrderHistoryResponse, error){
	return a.storage.Order().HistoryOrder(req)
}
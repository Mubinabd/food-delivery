package service

import (
	"context"

	pb "gitlab.com/bahodirova/product/genproto/product"
	"gitlab.com/bahodirova/product/storage"
)

type OrderItemService struct {
	storage storage.StorageI
	pb.UnimplementedOrderItemServiceServer
}

func NewOrderItemService(storage storage.StorageI) *OrderItemService {
	return &OrderItemService{
		storage: storage,
	}
}

func (a *OrderItemService) CreateOrderItem(ctx context.Context, req *pb.CreateOrderItemRequest) (*pb.Empty, error) {
	return a.storage.OrderItem().CreateOrderItem(req)
}
func (a *OrderItemService) GetOrderItem(ctx context.Context, req *pb.GetById) (*pb.OrderItem, error) {
	return a.storage.OrderItem().GetOrderItem(req)
}
func (a *OrderItemService) GetAllOrderItems(ctx context.Context, req *pb.GetAllOrderItemsReq) (*pb.GetAllOrderItemsRes, error) {
	return a.storage.OrderItem().GetAllOrderItems(req)
}
func (a *OrderItemService) UpdateOrderItem(ctx context.Context, req *pb.UpdateOrderItemRequest) (*pb.UpdateOrderItemResponse, error) {
	return a.storage.OrderItem().UpdateOrderItem(req)
}
func (a *OrderItemService) GetOrderItemsByOrder(ctx context.Context, req *pb.GetByOrderReq) (*pb.GetAllOrderItemsRes, error) {
	return a.storage.OrderItem().GetOrderItemsByOrder(req)
}
func (a *OrderItemService) GetOrderItemsByProduct(ctx context.Context, req *pb.GetByProductReq) (*pb.GetAllOrderItemsRes, error) {
	return a.storage.OrderItem().GetOrderItemsByProduct(req)
}



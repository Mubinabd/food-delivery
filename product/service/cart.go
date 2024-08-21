package service

import (
	"context"

	pb "gitlab.com/bahodirova/product/genproto/product"
	"gitlab.com/bahodirova/product/storage"
)

type CartService struct {
	storage storage.StorageI
	pb.UnimplementedCartServiceServer
}

func NewCartService(storage storage.StorageI) *CartService {
	return &CartService{
		storage: storage,
	}
}

func (a *CartService) CreateCart(ctx context.Context, req *pb.CreateCartReq) (*pb.Empty, error) {
	return a.storage.Cart().CreateCart(req)
}
func (a *CartService) GetCart(ctx context.Context, req *pb.GetById) (*pb.Cart, error) {
	return a.storage.Cart().GetCart(req)
}
func (a *CartService) GetAllCarts(ctx context.Context, req *pb.GetAllCartsReq) (*pb.GetAllCartsRes, error) {
	return a.storage.Cart().GetAllCarts(req)
}
func (a *CartService) DeleteCart(ctx context.Context, req *pb.GetById) (*pb.DeleteCartResp, error) {
	return a.storage.Cart().DeleteCart(req)
}
func (a *CartService) UpdateCart(ctx context.Context, req *pb.UpdateCartReq) (*pb.UpdateCartRes, error) {
	return a.storage.Cart().UpdateCart(req)
}


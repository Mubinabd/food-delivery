package service

import (
	"context"

	pb "gitlab.com/bahodirova/product/genproto/product"
	"gitlab.com/bahodirova/product/storage"
)

type ProductService struct {
	storage storage.StorageI
	pb.UnimplementedProductServiceServer
}

func NewProductService(storage storage.StorageI) *ProductService {
	return &ProductService{
		storage: storage,
	}
}

func (a *ProductService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Empty, error) {
	return a.storage.Product().CreateProduct(req)
}
func (a *ProductService) GetProduct(ctx context.Context, req *pb.GetById) (*pb.Product, error) {
	return a.storage.Product().GetProduct(req)
}
func (a *ProductService) GetAllProducts(ctx context.Context, req *pb.GetAllProductsReq) (*pb.GetAllProductsRes, error) {
	return a.storage.Product().GetAllProducts(req)
}
func (a *ProductService) DeleteProduct(ctx context.Context, req *pb.GetById) (*pb.DeleteProductResponse, error) {
	return a.storage.Product().DeleteProduct(req)
}
func (a *ProductService) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	return a.storage.Product().UpdateProduct(req)
}
func (a *ProductService) SearchProducts(ctx context.Context, req *pb.SearchProductsReq) (*pb.GetAllProductsRes, error) {
	return a.storage.Product().SearchProducts(req)
}


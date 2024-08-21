package service

import (
	"context"

	pb "gitlab.com/bahodirova/product/genproto/product"
	"gitlab.com/bahodirova/product/storage"
)

type CourierLocationService struct {
	storage storage.StorageI
	pb.UnimplementedCourierLocationServiceServer
}

func NewCourierLocationService(storage storage.StorageI) *CourierLocationService {
	return &CourierLocationService{
		storage: storage,
	}
}

func (a *CourierLocationService) CreateCourierLocation(ctx context.Context, req *pb.CreateCourierLocationRequest) (*pb.Empty, error) {
	return a.storage.CourierLocation().CreateCourierLocation(req)
}
func (a *CourierLocationService) GetCourierLocation(ctx context.Context, req *pb.GetById) (*pb.CourierLocation, error) {
	return a.storage.CourierLocation().GetCourierLocation(req)
}
func (a *CourierLocationService) GetAllCourierLocations(ctx context.Context, req *pb.GetAllCourierLocationsReq) (*pb.GetAllCourierLocationsRes, error) {
	return a.storage.CourierLocation().GetAllCourierLocations(req)
}
func (a *CourierLocationService) UpdateCourierLocation(ctx context.Context, req *pb.UpdateCourierLocationRequest) (*pb.UpdateCourierLocationResponse, error) {
	return a.storage.CourierLocation().UpdateCourierLocation(req)
}
func (a *CourierLocationService) GetCourierLocationsByTimeRange(ctx context.Context, req *pb.GetCourierLocationsByTimeRangeReq) (*pb.GetCourierLocationsByTimeRangeRes, error) {
	return a.storage.CourierLocation().GetCourierLocationsByTimeRange(req)
}
func (a *CourierLocationService) UpdateCourierLocationStatus(ctx context.Context, req *pb.UpdateCourierLocationStatusReq) (*pb.UpdateCourierLocationStatusRes, error) {
	return a.storage.CourierLocation().UpdateCourierLocationStatus(req)
}


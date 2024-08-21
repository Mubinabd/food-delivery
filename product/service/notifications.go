package service

import (
	"context"

	pb "gitlab.com/bahodirova/product/genproto/product"
	"gitlab.com/bahodirova/product/storage"
)

type NotificationService struct {
	storage storage.StorageI
	pb.UnimplementedNotificationServiceServer
}

func NewNotificationService(storage storage.StorageI) *NotificationService {
	return &NotificationService{
		storage: storage,
	}
}

func (a *NotificationService) CreateNotification(ctx context.Context, req *pb.CreateNotificationReq) (*pb.Empty, error) {
	return a.storage.Notification().CreateNotification(req)
}
func (a *NotificationService) GetNotification(ctx context.Context, req *pb.GetById) (*pb.Notification, error) {
	return a.storage.Notification().GetNotification(req)
}
func (a *NotificationService) GetAllNotifications(ctx context.Context, req *pb.GetAllNotificationsReq) (*pb.GetAllNotificationsRes, error) {
	return a.storage.Notification().GetAllNotifications(req)
}
func (a *NotificationService) MarkNotificationAsRead(ctx context.Context, req *pb.MarkNotificationAsReadReq) (*pb.MarkNotificationAsReadResp, error) {
	return a.storage.Notification().MarkNotificationAsRead(req)
}

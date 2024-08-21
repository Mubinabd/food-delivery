// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v3.12.4
// source: protos/notification.proto

package product

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	NotificationService_CreateNotification_FullMethodName     = "/product.NotificationService/CreateNotification"
	NotificationService_GetAllNotifications_FullMethodName    = "/product.NotificationService/GetAllNotifications"
	NotificationService_GetNotification_FullMethodName        = "/product.NotificationService/GetNotification"
	NotificationService_MarkNotificationAsRead_FullMethodName = "/product.NotificationService/MarkNotificationAsRead"
)

// NotificationServiceClient is the client API for NotificationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NotificationServiceClient interface {
	CreateNotification(ctx context.Context, in *CreateNotificationReq, opts ...grpc.CallOption) (*Empty, error)
	GetAllNotifications(ctx context.Context, in *GetAllNotificationsReq, opts ...grpc.CallOption) (*GetAllNotificationsRes, error)
	GetNotification(ctx context.Context, in *GetById, opts ...grpc.CallOption) (*Notification, error)
	MarkNotificationAsRead(ctx context.Context, in *MarkNotificationAsReadReq, opts ...grpc.CallOption) (*MarkNotificationAsReadResp, error)
}

type notificationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNotificationServiceClient(cc grpc.ClientConnInterface) NotificationServiceClient {
	return &notificationServiceClient{cc}
}

func (c *notificationServiceClient) CreateNotification(ctx context.Context, in *CreateNotificationReq, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, NotificationService_CreateNotification_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) GetAllNotifications(ctx context.Context, in *GetAllNotificationsReq, opts ...grpc.CallOption) (*GetAllNotificationsRes, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAllNotificationsRes)
	err := c.cc.Invoke(ctx, NotificationService_GetAllNotifications_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) GetNotification(ctx context.Context, in *GetById, opts ...grpc.CallOption) (*Notification, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Notification)
	err := c.cc.Invoke(ctx, NotificationService_GetNotification_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) MarkNotificationAsRead(ctx context.Context, in *MarkNotificationAsReadReq, opts ...grpc.CallOption) (*MarkNotificationAsReadResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MarkNotificationAsReadResp)
	err := c.cc.Invoke(ctx, NotificationService_MarkNotificationAsRead_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NotificationServiceServer is the server API for NotificationService service.
// All implementations must embed UnimplementedNotificationServiceServer
// for forward compatibility
type NotificationServiceServer interface {
	CreateNotification(context.Context, *CreateNotificationReq) (*Empty, error)
	GetAllNotifications(context.Context, *GetAllNotificationsReq) (*GetAllNotificationsRes, error)
	GetNotification(context.Context, *GetById) (*Notification, error)
	MarkNotificationAsRead(context.Context, *MarkNotificationAsReadReq) (*MarkNotificationAsReadResp, error)
	mustEmbedUnimplementedNotificationServiceServer()
}

// UnimplementedNotificationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNotificationServiceServer struct {
}

func (UnimplementedNotificationServiceServer) CreateNotification(context.Context, *CreateNotificationReq) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateNotification not implemented")
}
func (UnimplementedNotificationServiceServer) GetAllNotifications(context.Context, *GetAllNotificationsReq) (*GetAllNotificationsRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllNotifications not implemented")
}
func (UnimplementedNotificationServiceServer) GetNotification(context.Context, *GetById) (*Notification, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNotification not implemented")
}
func (UnimplementedNotificationServiceServer) MarkNotificationAsRead(context.Context, *MarkNotificationAsReadReq) (*MarkNotificationAsReadResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MarkNotificationAsRead not implemented")
}
func (UnimplementedNotificationServiceServer) mustEmbedUnimplementedNotificationServiceServer() {}

// UnsafeNotificationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NotificationServiceServer will
// result in compilation errors.
type UnsafeNotificationServiceServer interface {
	mustEmbedUnimplementedNotificationServiceServer()
}

func RegisterNotificationServiceServer(s grpc.ServiceRegistrar, srv NotificationServiceServer) {
	s.RegisterService(&NotificationService_ServiceDesc, srv)
}

func _NotificationService_CreateNotification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateNotificationReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).CreateNotification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_CreateNotification_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).CreateNotification(ctx, req.(*CreateNotificationReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_GetAllNotifications_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllNotificationsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).GetAllNotifications(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_GetAllNotifications_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).GetAllNotifications(ctx, req.(*GetAllNotificationsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_GetNotification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetById)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).GetNotification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_GetNotification_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).GetNotification(ctx, req.(*GetById))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_MarkNotificationAsRead_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MarkNotificationAsReadReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).MarkNotificationAsRead(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_MarkNotificationAsRead_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).MarkNotificationAsRead(ctx, req.(*MarkNotificationAsReadReq))
	}
	return interceptor(ctx, in, info, handler)
}

// NotificationService_ServiceDesc is the grpc.ServiceDesc for NotificationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NotificationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "product.NotificationService",
	HandlerType: (*NotificationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateNotification",
			Handler:    _NotificationService_CreateNotification_Handler,
		},
		{
			MethodName: "GetAllNotifications",
			Handler:    _NotificationService_GetAllNotifications_Handler,
		},
		{
			MethodName: "GetNotification",
			Handler:    _NotificationService_GetNotification_Handler,
		},
		{
			MethodName: "MarkNotificationAsRead",
			Handler:    _NotificationService_MarkNotificationAsRead_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/notification.proto",
}

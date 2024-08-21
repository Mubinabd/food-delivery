package storage

import (
	pb "gitlab.com/bahodirova/product/genproto/product"
)
type StorageI interface {
	Product() Product
	Cart() Cart
	Order() Order
	OrderItem() OrderItem
	Notification() Notification
	CourierLocation() CourierLocation
	Task() Task
}

type Product interface {
	CreateProduct(req *pb.CreateProductRequest) (*pb.Empty, error)
	GetProduct(req *pb.GetById) (*pb.Product, error)
	GetAllProducts(req *pb.GetAllProductsReq) (*pb.GetAllProductsRes, error)
	UpdateProduct(req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error)
	DeleteProduct(req *pb.GetById) (*pb.DeleteProductResponse, error)
	SearchProducts(req *pb.SearchProductsReq) (*pb.GetAllProductsRes, error)
}

type Cart interface {
	CreateCart(req *pb.CreateCartReq) (*pb.Empty, error)
	GetCart(req *pb.GetById) (*pb.Cart, error)
	GetAllCarts(req *pb.GetAllCartsReq) (*pb.GetAllCartsRes, error)
	UpdateCart(req *pb.UpdateCartReq) (*pb.UpdateCartRes, error)
	DeleteCart(req *pb.GetById) (*pb.DeleteCartResp, error)
}
type Order interface {
	CreateOrder(req *pb.CreateOrderReq) (*pb.Empty, error)
	GetOrder(req *pb.GetById) (*pb.Order, error)
	GetAllOrders(req *pb.GetAllOrdersReq) (*pb.GetAllOrderRes, error)
	UpdateOrder(req *pb.UpdateOrderReq) (*pb.UpdateOrderRes, error)
	DeleteOrder(req *pb.GetById) (*pb.DeleteOrderRes, error)
	PaidOrder(req *pb.PaidReq) (*pb.PaidRes, error)
	HistoryOrder(req *pb.GetCourierOrderHistoryRequest) (*pb.GetCourierOrderHistoryResponse, error)
}


type Notification interface {
	CreateNotification(req *pb.CreateNotificationReq) (*pb.Empty, error)
	GetAllNotifications(req *pb.GetAllNotificationsReq) (*pb.GetAllNotificationsRes, error)
	GetNotification(req *pb.GetById) (*pb.Notification, error)
	MarkNotificationAsRead(req *pb.MarkNotificationAsReadReq) (*pb.MarkNotificationAsReadResp, error)
}

type CourierLocation interface {			
	CreateCourierLocation(req *pb.CreateCourierLocationRequest) (*pb.Empty, error)
	GetCourierLocation(req *pb.GetById) (*pb.CourierLocation, error)
	GetAllCourierLocations(req *pb.GetAllCourierLocationsReq) (*pb.GetAllCourierLocationsRes, error)
	UpdateCourierLocation(req *pb.UpdateCourierLocationRequest) (*pb.UpdateCourierLocationResponse, error)
	GetCourierLocationsByTimeRange(req *pb.GetCourierLocationsByTimeRangeReq) (*pb.GetCourierLocationsByTimeRangeRes, error)
	UpdateCourierLocationStatus(req *pb.UpdateCourierLocationStatusReq) (*pb.UpdateCourierLocationStatusRes, error)
}

type Task interface {
	CreateTask(req *pb.CreatetaskReq) (*pb.Empty, error)
	GetTask(req *pb.GetById) (*pb.Task, error)
	GetAllTasks(req *pb.GetAllTasksReq) (*pb.GetAllTasksRes, error)
	UpdateTask(req *pb.UpdateTaskReq) (*pb.UpdateTaskRes, error)
	DeleteTask(req *pb.GetById) (*pb.DeleteTaskRes, error)
	GetTasksByUser(req *pb.GetByUserReq) (*pb.GetAllTasksRes, error)
	SearchTasks(req *pb.SearchTasksReq) (*pb.GetAllTasksRes, error)
}
type OrderItem interface {
	CreateOrderItem(req *pb.CreateOrderItemRequest) (*pb.Empty, error)
	GetOrderItem(req *pb.GetById) (*pb.OrderItem, error)
	GetAllOrderItems(req *pb.GetAllOrderItemsReq) (*pb.GetAllOrderItemsRes, error)
	UpdateOrderItem(req *pb.UpdateOrderItemRequest) (*pb.UpdateOrderItemResponse, error)
	GetOrderItemsByOrder(req *pb.GetByOrderReq) (*pb.GetAllOrderItemsRes, error)
	GetOrderItemsByProduct(req *pb.GetByProductReq) (*pb.GetAllOrderItemsRes, error)
}


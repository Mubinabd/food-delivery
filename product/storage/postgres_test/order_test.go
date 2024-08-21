package postgres

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	pb "gitlab.com/bahodirova/product/genproto/product"
	"gitlab.com/bahodirova/product/storage/postgres"
)

func TestCreateOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	manager := postgres.NewOrderManager(db)

	req := &pb.CreateOrderReq{
		UserId:          "user-123",
		CourierId:       "courier-123",
		Status:          "pending",
		TotalAmount:     100,
		DeliveryAddress: "123 Main St",
	}

	mock.ExpectExec("INSERT INTO orders").
		WithArgs(sqlmock.AnyArg(), req.UserId, req.CourierId, req.Status, req.TotalAmount, req.DeliveryAddress).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = manager.CreateOrder(req)
	assert.NoError(t, err)
}

func TestGetOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	manager := postgres.NewOrderManager(db)

	id := uuid.NewString()
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "user_id", "courier_id", "status", "total_amount", "delivery_address", "created_at"}).
		AddRow(id, "user-123", "courier-123", "pending", 100, "123 Main St", now)

	mock.ExpectQuery("SELECT id, user_id, courier_id, status, total_amount, delivery_address, created_at FROM orders WHERE id = \\$1").
		WithArgs(id).
		WillReturnRows(rows)

	req := &pb.GetById{Id: id}
	order, err := manager.GetOrder(req)
	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, id, order.Id)
	assert.Equal(t, "user-123", order.UserId)
	assert.Equal(t, "courier-123", order.CourierId)
	assert.Equal(t, "pending", order.Status)
	assert.Equal(t, float32(100), order.TotalAmount)
	assert.Equal(t, "123 Main St", order.DeliveryAddress)
}

func TestGetAllOrders(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	manager := postgres.NewOrderManager(db)

	rows := sqlmock.NewRows([]string{"id", "user_id", "courier_id", "status", "total_amount", "delivery_address", "created_at"}).
		AddRow(uuid.NewString(), "user-123", "courier-123", "pending", 100, "123 Main St", time.Now())

	mock.ExpectQuery("SELECT id, user_id, courier_id, status, total_amount, delivery_address, created_at FROM orders").
		WillReturnRows(rows)

	req := &pb.GetAllOrdersReq{Status: "pending"}
	orders, err := manager.GetAllOrders(req)
	assert.NoError(t, err)
	assert.NotNil(t, orders)
	assert.Greater(t, len(orders.Orders), 0)
}

func TestUpdateOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	manager := postgres.NewOrderManager(db)

	req := &pb.UpdateOrderReq{
		Id:              "order-123",
		UserId:          "user-123",
		CourierId:       "courier-123",
		Status:          "completed",
		TotalAmount:     150,
		DeliveryAddress: "456 Main St",
	}

	query := `
		UPDATE orders
		SET user_id = \$1, courier_id = \$2, total_amount = \$3, status = \$4, delivery_address = \$5, updated_at = \$6
		WHERE id = \$7`

	mock.ExpectBegin()
	mock.ExpectExec(query).
		WithArgs(req.UserId, req.CourierId, req.TotalAmount, req.Status, req.DeliveryAddress, sqlmock.AnyArg(), req.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	resp, err := manager.UpdateOrder(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, "Order updated successfully", resp.Message)
}



func TestPaidOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	manager := postgres.NewOrderManager(db)

	req := &pb.PaidReq{
		NewcartId:   "cart-123",
		ProductId:   "product-123",
		TotalAmount: 50,
	}

	mock.ExpectQuery("SELECT quantity, user_id FROM cart WHERE id = \\$1").
		WithArgs(req.NewcartId).
		WillReturnRows(sqlmock.NewRows([]string{"quantity", "user_id"}).AddRow(100.0, "user-123"))

	mock.ExpectQuery("SELECT o.user_id FROM orders o JOIN order_items oi ON o.id = oi.order_id WHERE oi.product_id = \\$1").
		WithArgs(req.ProductId).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow("user-123"))

	mock.ExpectQuery("SELECT price FROM products WHERE id = \\$1").
		WithArgs(req.ProductId).
		WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(50.0))

	mock.ExpectExec("UPDATE cart SET quantity = \\$1 WHERE id = \\$2").
		WithArgs(50.0, req.NewcartId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	resp, err := manager.PaidOrder(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, "Payment successful", resp.Message)
}

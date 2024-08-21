package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "gitlab.com/bahodirova/product/genproto/product"
	"gitlab.com/bahodirova/product/storage/postgres"
)

func TestCreateOrderItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	manager := postgres.NewOrderItemManager(db)

	req := &pb.CreateOrderItemRequest{
		OrderId:   uuid.NewString(),
		ProductId: uuid.NewString(),
		Quantity:  2,
		Price:     100.0,
	}

	mock.ExpectExec("INSERT INTO order_items").
		WithArgs(sqlmock.AnyArg(), req.OrderId, req.ProductId, req.Quantity, req.Price).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = manager.CreateOrderItem(req)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetOrderItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	manager := postgres.NewOrderItemManager(db)

	orderItemID := uuid.NewString()
	req := &pb.GetById{Id: orderItemID}

	mock.ExpectQuery("SELECT o.id, o.order_id, o.product_id, o.quantity, o.price, p.id, p.name, p.description, p.price, p.image_url, p.created_at, o2.id, o2.user_id, o2.courier_id, o2.status, o2.total_amount, o2.delivery_address, o2.created_at FROM order_items o").
		WithArgs(orderItemID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "product_id", "quantity", "price", "p.id", "p.name", "p.description", "p.price", "p.image_url", "p.created_at", "o2.id", "o2.user_id", "o2.courier_id", "o2.status", "o2.total_amount", "o2.delivery_address", "o2.created_at"}).
			AddRow(orderItemID, "orderID", "productID", int32(2), float32(100.0), "productID", "ProductName", "ProductDescription", float32(50.0), "ProductImageUrl", "ProductCreatedAt", "orderID", "userID", "courierID", "status", float32(200.0), "DeliveryAddress", "OrderCreatedAt"))

	orderItem, err := manager.GetOrderItem(req)
	require.NoError(t, err)

	assert.Equal(t, orderItemID, orderItem.Id)
	assert.Equal(t, "orderID", orderItem.OrderId)
	assert.Equal(t, "productID", orderItem.ProductId)
	assert.Equal(t, int32(2), orderItem.Quantity)
	assert.Equal(t, float32(100.0), orderItem.Price)
	assert.Equal(t, "productID", orderItem.Product.Id)
	assert.Equal(t, "ProductName", orderItem.Product.Name)
	assert.Equal(t, "ProductDescription", orderItem.Product.Description)
	assert.Equal(t, int32(0), orderItem.Product.Price)
	assert.Equal(t, "ProductImageUrl", orderItem.Product.ImageUrl)
	assert.Equal(t, "ProductCreatedAt", orderItem.Product.CreatedAt)
	assert.Equal(t, "orderID", orderItem.Order.Id)
	assert.Equal(t, "userID", orderItem.Order.UserId)
	assert.Equal(t, "courierID", orderItem.Order.CourierId)
	assert.Equal(t, "status", orderItem.Order.Status)
	assert.Equal(t, float32(200.0), orderItem.Order.TotalAmount)
	assert.Equal(t, "DeliveryAddress", orderItem.Order.DeliveryAddress)
	assert.Equal(t, "OrderCreatedAt", orderItem.Order.CreatedAt)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}




func TestUpdateOrderItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	manager := postgres.NewOrderItemManager(db)

	req := &pb.UpdateOrderItemRequest{
		Id:        "id",
		OrderId:   "orderID",
		ProductId: "productID",
		Quantity:  3,
		Price:     120.0,
	}

	query := `UPDATE order_items SET order_id = \$1, product_id = \$2, quantity = \$3, price = \$4 WHERE id = \$5`

	mock.ExpectBegin()
	mock.ExpectExec(query).
		WithArgs(req.OrderId, req.ProductId, req.Quantity, req.Price, req.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	res, err := manager.UpdateOrderItem(req)
	require.NoError(t, err)

	assert.True(t, res.Success)
	assert.Equal(t, "Order updated successfully", res.Message)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetOrderItemsByOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	manager := postgres.NewOrderItemManager(db)

	req := &pb.GetByOrderReq{OrderId: "orderID"}

	mock.ExpectQuery("SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price, o.courier_id, o.created_at, o.delivery_address, o.id as order_id, o.status, o.total_amount, o.user_id FROM order_items oi").
		WithArgs(req.OrderId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "product_id", "quantity", "price", "courier_id", "created_at", "delivery_address", "order_id", "status", "total_amount", "user_id"}).
			AddRow("id", "orderID", "productID", 2, 100.0, "courierID", "createdAt", "deliveryAddress", "orderID", "status", 200.0, "userID"))

	res, err := manager.GetOrderItemsByOrder(req)
	require.NoError(t, err)

	assert.Len(t, res.OrderItems, 1)
	assert.Equal(t, "id", res.OrderItems[0].Id)
	assert.Equal(t, "orderID", res.OrderItems[0].OrderId)
	assert.Equal(t, "productID", res.OrderItems[0].ProductId)
	assert.Equal(t, int32(2), res.OrderItems[0].Quantity)
	assert.Equal(t, float32(100), res.OrderItems[0].Price)
	assert.Equal(t, "courierID", res.OrderItems[0].Order.CourierId)
	assert.Equal(t, "createdAt", res.OrderItems[0].Order.CreatedAt)
	assert.Equal(t, "deliveryAddress", res.OrderItems[0].Order.DeliveryAddress)
	assert.Equal(t, "orderID", res.OrderItems[0].Order.Id)
	assert.Equal(t, "status", res.OrderItems[0].Order.Status)
	assert.Equal(t, float32(200), res.OrderItems[0].Order.TotalAmount)
	assert.Equal(t, "userID", res.OrderItems[0].Order.UserId)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetOrderItemsByProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	manager := postgres.NewOrderItemManager(db)

	req := &pb.GetByProductReq{ProductId: "productID"}

	mock.ExpectQuery("SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price, p.created_at as product_created_at, p.description, p.id as product_id, p.image_url, p.name, p.price as product_price FROM order_items oi").
		WithArgs(req.ProductId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "product_id", "quantity", "price", "product_created_at", "description", "product_id", "image_url", "name", "product_price"}).
			AddRow("id", "orderID", "productID", 2, 100.0, "createdAt", "description", "productID", "imageUrl", "name", 50.0))

	res, err := manager.GetOrderItemsByProduct(req)
	require.NoError(t, err)

	assert.Len(t, res.OrderItems, 1)
	assert.Equal(t, "id", res.OrderItems[0].Id)
	assert.Equal(t, "orderID", res.OrderItems[0].OrderId)
	assert.Equal(t, "productID", res.OrderItems[0].ProductId)
	assert.Equal(t, int32(2), res.OrderItems[0].Quantity)
	assert.Equal(t, float32(100), res.OrderItems[0].Price)
	assert.Equal(t, "createdAt", res.OrderItems[0].Product.CreatedAt)
	assert.Equal(t, "description", res.OrderItems[0].Product.Description)
	assert.Equal(t, "productID", res.OrderItems[0].Product.Id)
	assert.Equal(t, "imageUrl", res.OrderItems[0].Product.ImageUrl)
	assert.Equal(t, "name", res.OrderItems[0].Product.Name)
	assert.Equal(t, int32(50.0), res.OrderItems[0].Product.Price)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

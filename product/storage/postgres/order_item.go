package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	pb "gitlab.com/bahodirova/product/genproto/product"
)

type OrderItemManager struct {
	db *sql.DB
}

func NewOrderItemManager(db *sql.DB) *OrderItemManager {
	return &OrderItemManager{db: db}
}

func (o *OrderItemManager) CreateOrderItem(req *pb.CreateOrderItemRequest) (*pb.Empty, error) {

	id := uuid.NewString()

	query := `
		INSERT INTO 
			order_items
				(id, 
				order_id, 
				product_id, 
				quantity, 
				price)
		VALUES 
			($1, $2, $3, $4, $5)`

	_, err := o.db.Exec(
		query,
		id,
		req.OrderId,
		req.ProductId,
		req.Quantity,
		req.Price)

	if err != nil {
		log.Println("Error while creating order item", err)
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (o *OrderItemManager) GetOrderItem(req *pb.GetById) (*pb.OrderItem, error) {
	query := `
		SELECT
			o.id, 
			o.order_id, 
			o.product_id, 
			o.quantity, 
			o.price,
			p.id,
			p.name,
			p.description,
			p.price,
			p.image_url,
			p.created_at,
			o2.id,
			o2.user_id,
			o2.courier_id,
			o2.status,
			o2.total_amount,
			o2.delivery_address,
			o2.created_at
		FROM  
			order_items o
		JOIN   
			products p
		ON  
			o.product_id = p.id
		JOIN    
			orders o2
		ON  
			o.order_id = o2.id
		WHERE  
			o.id = $1
		`

	row := o.db.QueryRow(query, req.Id)

	var orderItem pb.OrderItem
	orderItem.Product = &pb.Product{}
	orderItem.Order = &pb.Order{}

	var productprice float64

	err := row.Scan(
		&orderItem.Id,
		&orderItem.OrderId,
		&orderItem.ProductId,
		&orderItem.Quantity,
		&orderItem.Price,
		&orderItem.Product.Id,
		&orderItem.Product.Name,
		&orderItem.Product.Description,
		&productprice,
		&orderItem.Product.ImageUrl,
		&orderItem.Product.CreatedAt,
		&orderItem.Order.Id,
		&orderItem.Order.UserId,
		&orderItem.Order.CourierId,
		&orderItem.Order.Status,
		&orderItem.Order.TotalAmount,
		&orderItem.Order.DeliveryAddress,
		&orderItem.Order.CreatedAt,
	)
	if err != nil {
		log.Println("Row scan error:", err)
		return nil, err
	}
	return &orderItem, nil
}

func (o *OrderItemManager) GetAllOrderItems(req *pb.GetAllOrderItemsReq) (*pb.GetAllOrderItemsRes, error) {
	query := `
		SELECT
			o.id, 
			o.order_id, 
			o.product_id, 
			o.quantity, 
			o.price,
			p.id as product_id,
			p.name,
			p.description,
			p.price as product_price,
			p.image_url,
			p.created_at as product_created_at,
			o2.id as order_id,
			o2.user_id,
			o2.courier_id,
			o2.status,
			o2.total_amount,
			o2.delivery_address,
			o2.created_at as order_created_at
		FROM
			order_items o
		JOIN
			products p
		ON
			o.product_id = p.id
		JOIN 
			orders o2
		ON
			o.order_id = o2.id
	`

	var args []interface{}
	argCount := 1
	filters := []string{}

	if req.OrderId != "" {
		filters = append(filters, fmt.Sprintf("o.order_id = $%d", argCount))
		args = append(args, req.OrderId)
		argCount++
	}

	if req.ProductId != "" {
		filters = append(filters, fmt.Sprintf("o.product_id = $%d", argCount))
		args = append(args, req.ProductId)
		argCount++
	}

	if len(filters) > 0 {
		query += " AND " + strings.Join(filters, " AND ")
	}

	if req.Filter != nil {
		if req.Filter.Limit > 0 {
			query += fmt.Sprintf(" LIMIT $%d", argCount)
			args = append(args, req.Filter.Limit)
			argCount++
		}
		if req.Filter.Offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, req.Filter.Offset)
			argCount++
		}
	}

	rows, err := o.db.Query(query, args...)
	if err != nil {
		log.Println("Query execution error:", err)
		return nil, err
	}
	defer rows.Close()

	var orderItems []*pb.OrderItem
	var price float64
	for rows.Next() {
		var orderItem pb.OrderItem
		orderItem.Product = &pb.Product{}
		orderItem.Order = &pb.Order{}
		err := rows.Scan(
			&orderItem.Id,
			&orderItem.OrderId,
			&orderItem.ProductId,
			&orderItem.Quantity,
			&orderItem.Price,
			&orderItem.Product.Id,
			&orderItem.Product.Name,
			&orderItem.Product.Description,
			&price,
			&orderItem.Product.ImageUrl,
			&orderItem.Product.CreatedAt,
			&orderItem.Order.Id,
			&orderItem.Order.UserId,
			&orderItem.Order.CourierId,
			&orderItem.Order.Status,
			&orderItem.Order.TotalAmount,
			&orderItem.Order.DeliveryAddress,
			&orderItem.Order.CreatedAt,
		)
		if err != nil {
			log.Println("Row scan error:", err)
			return nil, err
		}
		orderItem.Product.Price = int32(price)
		orderItems = append(orderItems, &orderItem)
	}
	return &pb.GetAllOrderItemsRes{OrderItems: orderItems}, nil
}

func (o *OrderItemManager) UpdateOrderItem(req *pb.UpdateOrderItemRequest) (*pb.UpdateOrderItemResponse, error) {
	var args []interface{}
	var conditions []string

	if req.OrderId != "" && req.OrderId != "string" {
		args = append(args, req.OrderId)
		conditions = append(conditions, fmt.Sprintf("order_id = $%d", len(args)))
	}
	if req.ProductId != "" && req.ProductId != "string" {
		args = append(args, req.ProductId)
		conditions = append(conditions, fmt.Sprintf("product_id = $%d", len(args)))
	}
	if req.Quantity != 0 {
		args = append(args, req.Quantity)
		conditions = append(conditions, fmt.Sprintf("quantity = $%d", len(args)))
	}
	if req.Price != 0 {
		args = append(args, req.Price)
		conditions = append(conditions, fmt.Sprintf("price = $%d", len(args)))
	}

	query := `UPDATE order_items SET ` + strings.Join(conditions, ", ") + ` WHERE id = $` + fmt.Sprintf("%d", len(args)+1)
	args = append(args, req.Id)

	tx, err := o.db.Begin()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &pb.UpdateOrderItemResponse{Success: true, Message: "Order updated successfully"}, nil

}


func (o *OrderItemManager) GetOrderItemsByOrder(req *pb.GetByOrderReq) (*pb.GetAllOrderItemsRes, error) {

	query := `
	SELECT
		oi.id,
		oi.order_id,
		oi.product_id,
		oi.quantity,
		oi.price,
		o.courier_id,
		o.created_at,
		o.delivery_address,
		o.id as order_id,
		o.status,
		o.total_amount,
		o.user_id
	FROM
		order_items oi
	JOIN orders o ON oi.order_id = o.id
	WHERE
		oi.order_id = $1
	`
	var args []interface{}
	args = append(args, req.OrderId)
	argCount := 2

	if req.Filter != nil {
		if req.Filter.Limit > 0 {
			query += fmt.Sprintf(" LIMIT $%d", argCount)
			args = append(args, req.Filter.Limit)
			argCount++
		}
		if req.Filter.Offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, req.Filter.Offset)
			argCount++
		}
	}

	rows, err := o.db.Query(query, args...)
	if err != nil {
		log.Println("Query execution error:", err)
		return nil, err
	}
	defer rows.Close()

	var orderItems []*pb.OrderItem
	for rows.Next() {
		var orderItem pb.OrderItem
		var order pb.Order
		err := rows.Scan(
			&orderItem.Id,
			&orderItem.OrderId,
			&orderItem.ProductId,
			&orderItem.Quantity,
			&orderItem.Price,
			&order.CourierId,
			&order.CreatedAt,
			&order.DeliveryAddress,
			&order.Id,
			&order.Status,
			&order.TotalAmount,
			&order.UserId,
		)
		if err != nil {
			log.Println("Row scan error:", err)
			return nil, err
		}
		orderItem.Order = &order
		orderItems = append(orderItems, &orderItem)
	}
	return &pb.GetAllOrderItemsRes{OrderItems: orderItems}, nil
}

func (o *OrderItemManager) GetOrderItemsByProduct(req *pb.GetByProductReq) (*pb.GetAllOrderItemsRes, error) {

	query := `
	SELECT
		oi.id,
		oi.order_id,
		oi.product_id,
		oi.quantity,
		oi.price,
		p.created_at as product_created_at,
		p.description,
		p.id as product_id,
		p.image_url,
		p.name,
		p.price as product_price
	FROM
		order_items oi
	JOIN products p ON oi.product_id = p.id
	WHERE
		oi.product_id = $1
	`
	var args []interface{}
	args = append(args, req.ProductId)
	argCount := 2

	if req.Filter != nil {
		if req.Filter.Limit > 0 {
			query += fmt.Sprintf(" LIMIT $%d", argCount)
			args = append(args, req.Filter.Limit)
			argCount++
		}
		if req.Filter.Offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, req.Filter.Offset)
			argCount++
		}
	}

	rows, err := o.db.Query(query, args...)
	if err != nil {
		log.Println("Query execution error:", err)
		return nil, err
	}
	defer rows.Close()

	var orderItems []*pb.OrderItem
	for rows.Next() {
		var orderItem pb.OrderItem
		var product pb.Product
		var productPrice float64
		err := rows.Scan(
			&orderItem.Id,
			&orderItem.OrderId,
			&orderItem.ProductId,
			&orderItem.Quantity,
			&orderItem.Price,
			&product.CreatedAt,
			&product.Description,
			&product.Id,
			&product.ImageUrl,
			&product.Name,
			&productPrice,
		)
		if err != nil {
			log.Println("Row scan error:", err)
			return nil, err
		}
		product.Price = int32(productPrice)
		orderItem.Product = &product
		orderItems = append(orderItems, &orderItem)
	}
	return &pb.GetAllOrderItemsRes{OrderItems: orderItems}, nil
}

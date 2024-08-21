package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	pb "gitlab.com/bahodirova/product/genproto/product"
)

type OrderManager struct {
	db *sql.DB
}

func NewOrderManager(db *sql.DB) *OrderManager {
	return &OrderManager{db: db}
}

func (o *OrderManager) CreateOrder(req *pb.CreateOrderReq) (*pb.Empty, error) {
	id := uuid.NewString()
	query := `
		INSERT INTO 
			orders
				(id, 
				user_id, 
				courier_id, 
				status, 
				total_amount,
				delivery_address)
		VALUES 
			($1, $2, $3, $4, $5, $6)`

	_, err := o.db.Exec(
		query,
		id,
		req.UserId,
		req.CourierId,
		req.Status,
		req.TotalAmount,
		req.DeliveryAddress)

	if err != nil {
		log.Println("Error while creating order", err)
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (o *OrderManager) GetOrder(req *pb.GetById) (*pb.Order, error) {
	query := `
	SELECT 
		id,
		user_id,
		courier_id,
		status,
		total_amount,
		delivery_address,
		created_at
	FROM
		orders
	WHERE
		id = $1`

	row := o.db.QueryRow(query, req.Id)
	var order pb.Order

	err := row.Scan(
		&order.Id,
		&order.UserId,
		&order.CourierId,
		&order.Status,
		&order.TotalAmount,
		&order.DeliveryAddress,
		&order.CreatedAt,
	)
	if err != nil {
		log.Println("no rows result set")
		return nil, err
	}
	return &order, nil
}

func (o *OrderManager) GetAllOrders(req *pb.GetAllOrdersReq) (*pb.GetAllOrderRes, error) {
	query := `
	SELECT 
		id,
		user_id,
		courier_id,
		status,
		total_amount,
		delivery_address,
		created_at
	FROM
		orders`

	var args []interface{}
	argCount := 1
	filters := []string{}

	if req.Status != "" {
		filters = append(filters, fmt.Sprintf("status = $%d", argCount))
		args = append(args, req.Status)
		argCount++
	}

	if req.TotalAmount != 0 {
		filters = append(filters, fmt.Sprintf("total_amount = $%d", argCount))
		args = append(args, req.TotalAmount)
		argCount++
	}

	if req.DeliveryAddress != "" {
		filters = append(filters, fmt.Sprintf("delivery_address = $%d", argCount))
		args = append(args, req.DeliveryAddress)
		argCount++
	}

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
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
		return nil, err
	}
	defer rows.Close()

	orders := []*pb.Order{}

	for rows.Next() {
		var order pb.Order
		err := rows.Scan(
			&order.Id,
			&order.UserId,
			&order.CourierId,
			&order.Status,
			&order.TotalAmount,
			&order.DeliveryAddress,
			&order.CreatedAt,
		)
		if err != nil {
			log.Println("no rows result set")
			return nil, err
		}
		orders = append(orders, &order)
	}
	return &pb.GetAllOrderRes{Orders: orders}, nil
}

func (o *OrderManager) UpdateOrder(req *pb.UpdateOrderReq) (*pb.UpdateOrderRes, error) {
	var args []interface{}
	var conditions []string

	if req.UserId != "" && req.UserId != "string" {
		args = append(args, req.UserId)
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", len(args)))
	}
	if req.CourierId != "" && req.CourierId != "string" {
		args = append(args, req.CourierId)
		conditions = append(conditions, fmt.Sprintf("courier_id = $%d", len(args)))
	}
	if req.TotalAmount != 0 {
		args = append(args, req.TotalAmount)
		conditions = append(conditions, fmt.Sprintf("total_amount = $%d", len(args)))
	}
	if req.Status != "" && req.Status != "string" {
		args = append(args, req.Status)
		conditions = append(conditions, fmt.Sprintf("status = $%d", len(args)))
	}
	if req.DeliveryAddress != "" && req.DeliveryAddress != "string" {
		args = append(args, req.DeliveryAddress)
		conditions = append(conditions, fmt.Sprintf("delivery_address = $%d", len(args)))
	}
	args = append(args, time.Now())
	conditions = append(conditions, fmt.Sprintf("updated_at = $%d", len(args)))

	query := `UPDATE orders SET ` + strings.Join(conditions, ", ") + ` WHERE id = $` + fmt.Sprintf("%d", len(args)+1)
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

	return &pb.UpdateOrderRes{Success: true, Message: "Order updated successfully"}, nil

}

func (o *OrderManager) DeleteOrder(req *pb.GetById) (*pb.DeleteOrderRes, error) {
	query := `
	UPDATE
		orders
	SET
		deleted_at = extract(epoch from now())
	WHERE
		id = $1
	`

	_, err := o.db.Exec(query, req.Id)
	if err != nil {
		log.Println("Error while deleting order", err)
		return nil, err
	}
	return &pb.DeleteOrderRes{Success: true, Message: "Order deleted successfully"}, nil
}

func (o *OrderManager) PaidOrder(req *pb.PaidReq) (*pb.PaidRes, error) {
    var cartSumma float64
    var cartUserID string
    var orderUserID string

    err := o.db.QueryRow("SELECT quantity, user_id FROM cart WHERE id = $1", req.NewcartId).Scan(&cartSumma, &cartUserID)
    if err != nil {
        log.Println("error while getting newcart summa", err)
        return nil, err
    }

    err = o.db.QueryRow(`
        SELECT 
            o.user_id
        FROM 
            orders o
        JOIN 
            order_items oi 
        ON 
            o.id = oi.order_id
        WHERE 
            oi.product_id = $1`, req.ProductId).Scan(&orderUserID)

    if err != nil {
        if err == sql.ErrNoRows {
            log.Println("no rows in result set for order user_id")
            return &pb.PaidRes{Success: false, Message: "Order not found for the given product"}, nil
        }
        log.Println("error while getting order user_id", err)
        return nil, err
    }

    var productPrice float64
    err = o.db.QueryRow("SELECT price FROM products WHERE id = $1", req.ProductId).Scan(&productPrice)
    if err != nil {
        log.Println("error while getting product price", err)
        return nil, err
    }

    remainingSumma := cartSumma - float64(req.TotalAmount)
    if remainingSumma < 0 {
        return &pb.PaidRes{Success: false, Message: "Insufficient funds"}, nil
    }

    _, err = o.db.Exec("UPDATE cart SET quantity = $1 WHERE id = $2", remainingSumma, req.NewcartId)
    if err != nil {
        log.Println("error while updating newcart summa", err)
        return nil, err
    }

    return &pb.PaidRes{Success: true, Message: "Payment successful"}, nil
}

func(o *OrderManager)HistoryOrder(req *pb.GetCourierOrderHistoryRequest) (*pb.GetCourierOrderHistoryResponse, error) {
	query := `
		SELECT
			id,
			user_id,
			courier_id,
			status,
			total_amount,
			delivery_address,
			created_at
		FROM
			orders
		WHERE
			courier_id = $1
		`	
	rows, err := o.db.Query(query, req.CourierId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*pb.Order

	for rows.Next() {
		var order pb.Order
		err := rows.Scan(
			&order.Id,
			&order.UserId,
			&order.CourierId,
			&order.Status,
			&order.TotalAmount,
			&order.DeliveryAddress,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}
	return &pb.GetCourierOrderHistoryResponse{Orders: orders}, nil
	}				

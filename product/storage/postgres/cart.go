package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	pb "gitlab.com/bahodirova/product/genproto/product"
)

type CartManager struct {
	db *sql.DB
}

func NewCartManager(db *sql.DB) *CartManager {
	return &CartManager{db: db}
}


func (c *CartManager) CreateCart(req *pb.CreateCartReq) (*pb.Empty, error) {
	id := uuid.NewString()
	optionsJSON, err := json.Marshal(req.Options)
	if err != nil {
		log.Println("Error while marshaling options to JSON", err)
		return nil, err
	}

	query := `
		INSERT INTO 
			cart 	
				(id, 
				user_id, 
				product_id, 
				quantity, 
				options,
				name,
				number)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7)`

	_, err = c.db.Exec(
		query,
		id,
		req.UserId,
		req.ProductId,
		req.Quantity,
		optionsJSON,
		req.Name,
		req.Number,
	)

	if err != nil {
		log.Println("Error while creating cart", err)
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (c *CartManager) GetCart(req *pb.GetById) (*pb.Cart, error) {
	query := `
		SELECT
			c.id,
			c.user_id,
			c.product_id,
			c.quantity,
			c.options,
			c.name,
			c.number,
			c.created_at,
			p.name,
			p.description,
			p.price,
			p.image_url
		FROM
			cart c
		JOIN
			products p
		ON
			c.product_id = p.id
		WHERE
			c.id = $1
		AND 
			c.deleted_at = 0`

	row := c.db.QueryRow(query, req.Id)
	var cart pb.Cart
	cart.ProductId = &pb.Product{}
	var price,number float64
	err := row.Scan(
		&cart.Id,
		&cart.UserId,
		&cart.ProductId.Id,
		&cart.Quantity,
		&cart.Options,
		&cart.Name,
		&number,
		&cart.CreatedAt,
		&cart.ProductId.Name,
		&cart.ProductId.Description,
		&price,
		&cart.ProductId.ImageUrl,
	)
	if err != nil {
		log.Println("Error while getting cart", err)
		return nil, err
	}

	cart.ProductId.Price = int32(price)
	cart.Number = int64(number)
	return &cart, nil
}

func (c *CartManager) GetAllCarts(req *pb.GetAllCartsReq) (*pb.GetAllCartsRes, error) {
	query := `SELECT
			c.id,
			c.user_id,
			c.product_id,
			c.quantity,
			c.options,
			c.name,
			c.number,
			c.created_at,
			p.name,
			p.description,
			p.price,
			p.image_url
		FROM
			cart c
		JOIN
			products p 
		ON 
			c.product_id = p.id`

	var args []interface{}
	argCount := 1
	filters := []string{}

	if req.Quantity != 0 {
		filters = append(filters, fmt.Sprintf("c.quantity = $%d", argCount))
		args = append(args, req.Quantity)
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

	rows, err := c.db.Query(query, args...)
	if err != nil {
		log.Println("Error while getting all carts", err)
		return nil, err
	}
	defer rows.Close()

	carts := []*pb.Cart{}

	for rows.Next() {
		var cart pb.Cart
		cart.ProductId = &pb.Product{}

		var price,number float64
		err = rows.Scan(
			&cart.Id,
			&cart.UserId,
			&cart.ProductId.Id,
			&cart.Quantity,
			&cart.Options,
			&cart.Name,
			&cart.Number,
			&cart.CreatedAt,
			&cart.ProductId.Name,
			&cart.ProductId.Description,
			&price,
			&cart.ProductId.ImageUrl,
		)
		if err != nil {
			log.Println("Error while scanning cart row", err)
			return nil, err
		}
		cart.ProductId.Price = int32(price)
		cart.Number = int64(number)

		carts = append(carts, &cart)
	}

	return &pb.GetAllCartsRes{Carts: carts}, nil
}

func (c *CartManager) UpdateCart(req *pb.UpdateCartReq) (*pb.UpdateCartRes, error) {
	var args []interface{}
	var conditions []string

	optionsJSON, err := json.Marshal(req.Options)
	if err != nil {
		log.Println("Error while marshaling options to JSON", err)
		return nil, err
	}

	if req.UserId != "" && req.UserId != "string" {
		args = append(args, req.UserId)
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", len(args)))
	}
	if req.ProductId != "" && req.ProductId != "string" {
		args = append(args, req.ProductId)
		conditions = append(conditions, fmt.Sprintf("product_id = $%d", len(args)))
	}
	if req.Quantity != 0 {
		args = append(args, req.Quantity)
		conditions = append(conditions, fmt.Sprintf("quantity = $%d", len(args)))
	}
	if string(optionsJSON) != "" && string(optionsJSON) != "string" {
		args = append(args, string(optionsJSON))
		conditions = append(conditions, fmt.Sprintf("options = $%d", len(args)))
	}
	if req.Name != "" && req.Name != "string" {
		args = append(args, req.Name)
		conditions = append(conditions, fmt.Sprintf("name = $%d", len(args)))
	}
	if req.Number != 0 {
		args = append(args, req.Number)
		conditions = append(conditions, fmt.Sprintf("number = $%d", len(args)))
	}

	args = append(args, time.Now())
	conditions = append(conditions, fmt.Sprintf("updated_at = $%d", len(args)))

	query := `UPDATE cart SET ` + strings.Join(conditions, ", ") + ` WHERE id = $` + fmt.Sprintf("%d", len(args)+1)
	args = append(args, req.Id)

	tx, err := c.db.Begin()
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

	return &pb.UpdateCartRes{Success: true, Message: "cart updated successfully"}, nil

}

func (c *CartManager) DeleteCart(req *pb.GetById) (*pb.DeleteCartResp, error) {
	query := `
	UPDATE
		cart
	SET
		deleted_at = extract(epoch from now())
	WHERE
		id = $1
	`

	_, err := c.db.Exec(
		query,
		req.Id)

	if err != nil {
		log.Println("Error while deleting cart", err)
		return nil, err
	}
	return &pb.DeleteCartResp{Success: true, Message: "Cart deleted successfully"}, nil
}
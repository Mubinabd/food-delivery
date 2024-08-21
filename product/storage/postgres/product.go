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

type ProductManager struct {
	db *sql.DB
}

func NewProductManager(db *sql.DB) *ProductManager {
	return &ProductManager{db: db}
}

func (p *ProductManager) CreateProduct(req *pb.CreateProductRequest) (*pb.Empty, error) {
	id := uuid.NewString()

	query := `
		INSERT INTO 
			products 
				(id, 
				name, 
				description, 
				price, 
				image_url)
		VALUES 
			($1, $2, $3, $4, $5)`

	_, err := p.db.Exec(
		query, 
		id, 
		req.Name, 
		req.Description, 
		req.Price, 
		req.ImageUrl)

	if err != nil {
		log.Println("Error while creating product", err)
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (p *ProductManager) GetProduct(req *pb.GetById) (*pb.Product, error) {
	query := `
		SELECT
			id,
			name,
			description,
			price,
			image_url,
			created_at
		FROM
			products
		WHERE
			id = $1
		AND 
			deleted_at = 0`

	row := p.db.QueryRow(query, req.Id)

	var product pb.Product
	var price float64

	err := row.Scan(
		&product.Id,
		&product.Name,
		&product.Description,
		&price,
		&product.ImageUrl,
		&product.CreatedAt)

	if err != nil {
		log.Println("Error while getting product", err)
		return nil, err
	}
	product.Price = int32(price)
	return &product, nil

}

func (p *ProductManager) GetAllProducts(req *pb.GetAllProductsReq) (*pb.GetAllProductsRes, error) {
	query := `
	SELECT 
		id,
		name,
		description,
		price,
		image_url,
		created_at
	 FROM 
	 	products`

	var args []interface{}
	argCount := 1
	filters := []string{}

	if req.Name != "" {
		filters = append(filters, fmt.Sprintf("name = $%d", argCount))
		args = append(args, req.Name)
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

	rows, err := p.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []*pb.Product{}

	for rows.Next() {
		var product pb.Product
		var price float64
		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&price,
			&product.ImageUrl,
			&product.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		product.Price = int32(price)
		products = append(products, &product)
	}
	return &pb.GetAllProductsRes{Products: products}, nil
}

func (p *ProductManager) UpdateProduct(req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	var args []interface{}
	var conditions []string

	if req.Name != "" && req.Name != "string" {
		args = append(args, req.Name)
		conditions = append(conditions, fmt.Sprintf("name = $%d", len(args)))
	}
	if req.Description != "" && req.Description != "string" {
		args = append(args, req.Description)
		conditions = append(conditions, fmt.Sprintf("description = $%d", len(args)))
	}
	if req.Price != 0  {
		args = append(args, req.Price)
		conditions = append(conditions, fmt.Sprintf("price = $%d", len(args)))
	}
	if req.ImageUrl != "" && req.ImageUrl != "string" {
		args = append(args, req.ImageUrl)
		conditions = append(conditions, fmt.Sprintf("image_url = $%d", len(args)))
	}
	args = append(args, time.Now())
	conditions = append(conditions, fmt.Sprintf("updated_at = $%d", len(args)))

	query := `UPDATE products SET ` + strings.Join(conditions, ", ") + ` WHERE id = $` + fmt.Sprintf("%d", len(args)+1)
	args = append(args, req.Id)

	tx, err := p.db.Begin()
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



	return &pb.UpdateProductResponse{Success: true, Message: "Product updated successfully"}, nil

}

func(p *ProductManager)DeleteProduct(req *pb.GetById) (*pb.DeleteProductResponse, error){
	query := `
	UPDATE
		products
	SET
		deleted_at = extract(epoch from now())
	WHERE
		id = $1
	`

	_, err := p.db.Exec(query, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteProductResponse{Success: true, Message: "Product deleted successfully"}, nil
}

func(p ProductManager)SearchProducts(req *pb.SearchProductsReq) (*pb.GetAllProductsRes, error){
	query := `
		SELECT
			id,
			name,
			description,
			price,
			image_url,
			created_at
		FROM
			products
		WHERE
			name = $1
		OR
			description = $1
		OR
			image_url = $1
		`

	rows, err := p.db.Query(query, req.Name)
	if err != nil {
		return nil, err
	}
	
	defer rows.Close()

	products := []*pb.Product{}
	for rows.Next() {
		var product pb.Product
		var price float64
		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&price,
			&product.ImageUrl,
			&product.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		product.Price = int32(price)
		products = append(products, &product)
	}
	return &pb.GetAllProductsRes{Products: products}, nil
}
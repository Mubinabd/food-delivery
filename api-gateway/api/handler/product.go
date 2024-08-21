package handler

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "gitlab.com/bahodirova/api-gateway/genproto/product"
)

// @Summary Create a new product
// @Description Create a new product with the given details
// @Tags Product
// @Accept json
// @Produce json
// @Param product body pb.CreateProductRequest true "Product"
// @Success 200 {object} string "message": "created successfully"
// @Failure 400 {object} string "error": "error message"
// @Router /api/product [POST]
func (h *HandlerStruct) CreateProduct(c *gin.Context) {
	var req pb.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	input, err := json.Marshal(req)
	err = h.Clients.KafkaProducer.ProduceMessages("create-product", input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		log.Println("cannot produce messages via kafka", err)
		return
	}

	// _, err := h.Clients.ProductClient.CreateProduct(context.Background(), &req)
	// if err != nil {
	// 	c.JSON(400, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	c.JSON(200, gin.H{
		"message": "created successfully",
	})
}

// @Summary Get a product by ID
// @Description Get a product's details by its ID
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} pb.Product
// @Failure 400 {object} string "error": "error message"
// @Router /api/product/{id} [GET]
func (h *HandlerStruct) GetProduct(c *gin.Context) {
	var req pb.GetById
	id := c.Param("id")
	req.Id = id
	res, err := h.Clients.ProductClient.GetProduct(context.Background(), &req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, res)
}

// @Summary Get all products
// @Description Get a list of all products with optional filters
// @Tags Product
// @Accept json
// @Produce json
// @Param name query string false "Product Name"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} pb.GetAllProductsRes
// @Failure 400 {object} string "error": "error message"
// @Router /api/product [GET]
func (h *HandlerStruct) GetAllProducts(c *gin.Context) {
	var req pb.GetAllProductsReq
	name := c.Query("name")
	limit := c.Query("limit")
	offset := c.Query("offset")

	req.Name = name

	if limit != "" {
		limitValue, err := strconv.Atoi(limit)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid limit value"})
			return
		}
		if req.Filter == nil {
			req.Filter = &pb.Filter{}
		}
		req.Filter.Limit = int32(limitValue)
	}

	if offset != "" {
		offsetValue, err := strconv.Atoi(offset)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid offset value"})
			return
		}
		if req.Filter == nil {
			req.Filter = &pb.Filter{}
		}
		req.Filter.Offset = int32(offsetValue)
	}
	res, err := h.Clients.ProductClient.GetAllProducts(context.Background(), &req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, res)
}

// @Summary Update a product
// @Description Update a product with the given details
// @Tags Product
// @Accept json
// @Produce json
// @Param product body pb.UpdateProductRequest true "Product"
// @Success 200 {object} string "message": "updated successfully"
// @Failure 400 {object} string "error": "error message"
// @Router /api/product [PUT]
func (h *HandlerStruct) UpdateProduct(c *gin.Context) {
	var req pb.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	input, err := json.Marshal(req)
	err = h.Clients.KafkaProducer.ProduceMessages("update-product", input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		log.Println("cannot produce messages via kafka", err)
		return
	}

	// _, err := h.Clients.ProductClient.UpdateProduct(context.Background(), &req)
	// if err != nil {
	// 	c.JSON(400, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }
	c.JSON(200, gin.H{
		"message": "updated successfully",
	})
}

// @Summary Delete a product
// @Description Delete a product by its ID
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} string "message": "deleted successfully"
// @Failure 400 {object} string "error": "error message"
// @Router /api/product/{id} [DELETE]
func (h *HandlerStruct) DeleteProduct(c *gin.Context) {
	var req pb.GetById
	id := c.Param("id")
	req.Id = id

	input, err := json.Marshal(req)
	err = h.Clients.KafkaProducer.ProduceMessages("delete-product", input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		log.Println("cannot produce messages via kafka", err)
		return
	}

	// _, err := h.Clients.ProductClient.DeleteProduct(context.Background(), &req)
	// if err != nil {
	// 	c.JSON(400, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	c.JSON(200, gin.H{
		"message": "deleted successfully",
	})
}

// @Summary Search products
// @Description Search for products with filters
// @Tags Product
// @Accept json
// @Produce json
// @Param name query string false "Product Name"
// @Param description query string false "Product Description"
// @Param price query int false "Product Price"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} pb.GetAllProductsRes
// @Failure 400 {object} string "error": "error message"
// @Router /api/product/search [GET]
func (h *HandlerStruct) SearchProducts(c *gin.Context) {
	var req pb.SearchProductsReq
	name := c.Query("name")
	description := c.Query("description")
	priceSTR := c.Query("price")
	price, err := strconv.Atoi(priceSTR)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid price value"})
		return
	}
	limit := c.Query("limit")
	offset := c.Query("offset")

	req.Name = name
	req.Description = description
	req.Price = int32(price)

	if limit != "" {
		limitValue, err := strconv.Atoi(limit)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid limit value"})
			return
		}
		if req.Filter == nil {
			req.Filter = &pb.Filter{}
		}
		req.Filter.Limit = int32(limitValue)
	}

	if offset != "" {
		offsetValue, err := strconv.Atoi(offset)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid offset value"})
			return
		}
		if req.Filter == nil {
			req.Filter = &pb.Filter{}
		}
		req.Filter.Offset = int32(offsetValue)
	}
	res, err := h.Clients.ProductClient.SearchProducts(context.Background(), &req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, res)
}

package handler

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "gitlab.com/bahodirova/api-gateway/genproto/product"
)

// @Router        /api/order-item [POST]
// @Summary       CREATE order item
// @Description   This API creates an order item
// @Tags          OrderItem
// @Accept        json
// @Produce       json
// @Security      BearerAuth
// @Param         data body pb.CreateOrderItemRequest true "Order Item"
// @Success       200 {object} string "message": "created successfully"
// @Failure       400 {object} string "error": "error message"
func (h *HandlerStruct) CreateOrderItem(c *gin.Context) {
	var req pb.CreateOrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err := h.Clients.OrderItemClient.CreateOrderItem(context.Background(), &req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "created successfully",
	})
}

// @Router        /api/order-item/{id} [GET]
// @Summary       GET order item
// @Description   This API retrieves an order item by ID
// @Tags          OrderItem
// @Accept        json
// @Produce       json
// @Param         id path string true "Order Item ID"
// @Success       200 {object} pb.OrderItem
// @Failure       400 {object} string "error": "error message"
func (h *HandlerStruct) GetOrderItem(c *gin.Context) {
	var req pb.GetById
	id := c.Param("id")
	req.Id = id
	res, err := h.Clients.OrderItemClient.GetOrderItem(context.Background(), &req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, res)
}

// @Router        /api/order-item/all [GET]
// @Summary       GET all order items
// @Description   This API retrieves all order items with optional filters
// @Tags          OrderItem
// @Accept        json
// @Produce       json
// @Param         order_id query string false "Order ID"
// @Param         product_id query string false "Product ID"
// @Param         limit query string false "Limit"
// @Param         offset query string false "Offset"
// @Success       200 {object} pb.GetAllOrderItemsRes
// @Failure       400 {object} string "error": "error message"
func (h *HandlerStruct) GetAllOrderItems(c *gin.Context) {
	var req pb.GetAllOrderItemsReq
	order := c.Query("order_id")
	product := c.Query("product_id")
	limit := c.Query("limit")
	offset := c.Query("offset")

	req.OrderId = order
	req.ProductId = product

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
	res, err := h.Clients.OrderItemClient.GetAllOrderItems(context.Background(), &req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, res)
}

// @Router        /api/order-item [PUT]
// @Summary       UPDATE order item
// @Description   This API updates an order item
// @Tags          OrderItem
// @Accept        json
// @Produce       json
// @Security      BearerAuth
// @Param         data body pb.UpdateOrderItemRequest true "Order Item"
// @Success       200 {object} string "message": "updated successfully"
// @Failure       400 {object} string "error": "error message"
func (h *HandlerStruct) UpdateOrderItem(c *gin.Context) {
	var req pb.UpdateOrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// input, err := json.Marshal(req)
	// err = h.Clients.KafkaProducer.ProduceMessages("update-order-item", input)
	// if err != nil {
	// 	c.JSON(400, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	log.Println("cannot produce messages via kafka", err)
	// 	return
	// }
	_, err := h.Clients.OrderItemClient.UpdateOrderItem(context.Background(), &req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "updated successfully",
	})
}


// @Router        /api/order-item/order [GET]
// @Summary       GET order items by order ID
// @Description   This API retrieves order items by order ID with optional filters
// @Tags          OrderItem
// @Accept        json
// @Produce       json
// @Param         order_id query string true "Order ID"
// @Param         limit query string false "Limit"
// @Param         offset query string false "Offset"
// @Success       200 {object} pb.GetAllOrderItemsRes
// @Failure       400 {object} string "error": "error message"
func (h *HandlerStruct) GetOrderItemsByOrder(c *gin.Context) {
	var req pb.GetByOrderReq
	order := c.Query("order_id")
	limit := c.Query("limit")
	offset := c.Query("offset")

	req.OrderId = order

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
	res, err := h.Clients.OrderItemClient.GetOrderItemsByOrder(context.Background(), &req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, res)
}

// @Router        /api/order-item/product [GET]
// @Summary       GET order items by product ID
// @Description   This API retrieves order items by product ID with optional filters
// @Tags          OrderItem
// @Accept        json
// @Produce       json
// @Param         product_id query string true "Product ID"
// @Param         limit query string false "Limit"
// @Param         offset query string false "Offset"
// @Success        200 {object} pb.GetAllOrderItemsRes
// @Failure        400 {object} string "error": "error message"
func (h *HandlerStruct) GetOrderItemsByProduct(c *gin.Context) {
	var req pb.GetByProductReq
	product := c.Query("product_id")
	limit := c.Query("limit")
	offset := c.Query("offset")

	req.ProductId = product

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
	res, err := h.Clients.OrderItemClient.GetOrderItemsByProduct(context.Background(), &req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, res)
}

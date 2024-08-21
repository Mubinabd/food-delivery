package handler

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "gitlab.com/bahodirova/api-gateway/genproto/product"
)

// @Router        /api/order [POST]
// @Summary       CREATE order
// @Description   This API creates an order
// @Tags          Order
// @Accept        json
// @Produce       json
// @Security      BearerAuth
// @Param         data body pb.CreateOrderReq true "Order"
// @Success       200 {object} string "message": "created successfully"
// @Failure       400 {object} string "error": "error message"
func (h *HandlerStruct) CreateOrder(c *gin.Context) {
	var req pb.CreateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	input, err := json.Marshal(req)
	err = h.Clients.KafkaProducer.ProduceMessages("create-order", input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		log.Println("cannot produce messages via kafka", err)
		return
	}

	// _, err := h.Clients.OrderClient.CreateOrder(context.Background(), &req)
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

// @Router        /api/order/{id} [GET]
// @Summary       GET order
// @Description   This API retrieves an order by ID
// @Tags          Order
// @Accept        json
// @Produce       json
// @Param         id path string true "Order ID"
// @Success       200 {object} pb.Order
// @Failure       400 {object} string "error": "error message"
func (h *HandlerStruct) GetOrder(c *gin.Context) {
	var req pb.GetById
	id := c.Param("id")
	req.Id = id
	res, err := h.Clients.OrderClient.GetOrder(context.Background(), &req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, res)
}

// @Router        /api/order/all [GET]
// @Summary       GET all orders
// @Description   This API retrieves all orders with optional filters
// @Tags          Order
// @Accept        json
// @Produce       json
// @Param         status query string false "Order Status"
// @Param         delivery_address query string false "Delivery Address"
// @Param         total_amount query string false "Total Amount"
// @Param         limit query string false "Limit"
// @Param         offset query string false "Offset"
// @Success       200 {object} pb.GetAllOrderRes
// @Failure       400 {object} string "error": "error message"
func (h *HandlerStruct) GetAllOrders(c *gin.Context) {
	var req pb.GetAllOrdersReq
	status := c.Query("status")
	delivery := c.Query("delivery_address")
	totalamountSTR := c.Query("total_amount")
	totalamount, err := strconv.Atoi(totalamountSTR)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid total amount value"})
		return
	}
	req.Status = status
	req.DeliveryAddress = delivery
	req.TotalAmount = float32(totalamount)
	limit := c.Query("limit")
	offset := c.Query("offset")

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
	res, err := h.Clients.OrderClient.GetAllOrders(context.Background(), &req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, res)
}

// @Router        /api/order [PUT]
// @Summary       UPDATE order
// @Description   This API updates an order
// @Tags          Order
// @Accept        json
// @Produce       json
// @Security      BearerAuth
// @Param         data body pb.UpdateOrderReq true "Order"
// @Success       200 {object} string "message": "updated successfully"
// @Failure       400 {object} string "error": "error message"
func (h *HandlerStruct) UpdateOrder(c *gin.Context) {
	var req pb.UpdateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	input, err := json.Marshal(req)
	err = h.Clients.KafkaProducer.ProduceMessages("update-order", input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		log.Println("cannot produce messages via kafka", err)
		return
	}

	// _, err := h.Clients.OrderClient.UpdateOrder(context.Background(), &req)
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

// @Router        /api/order/{id} [DELETE]
// @Summary       DELETE order
// @Description   This API deletes an order by ID
// @Tags          Order
// @Accept        json
// @Produce       json
// @Param         id path string true "Order ID"
// @Success       200 {object} string "message": "deleted successfully"
// @Failure       400 {object} string "error": "error message"
func (h *HandlerStruct) DeleteOrder(c *gin.Context) {
	var req pb.GetById
	id := c.Param("id")
	req.Id = id

	input, err := json.Marshal(req)
	err = h.Clients.KafkaProducer.ProduceMessages("delete-order", input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		log.Println("cannot produce messages via kafka", err)
		return
	}

	// _, err := h.Clients.OrderClient.DeleteOrder(context.Background(), &req)
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

// PaidOrder godoc
// @Summary Pay for an order
// @Description Deducts the product price from the user's cart and checks if the transaction is valid
// @Tags Orders
// @Accept json
// @Produce json
// @Param body body pb.PaidReq true "Paid Request"
// @Success 200 {object} pb.PaidRes
// @Failure 400 {object} string "message":"error while paid checkout"
// @Router /api/orders/paid [post]
func (h *HandlerStruct) PaidOrder(c *gin.Context) {
    var req pb.PaidReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, pb.PaidRes{Success: false, Message: "Invalid request"})
        return
    }

    res, err := h.Clients.OrderClient.PaidOrder(c.Request.Context(), &req)
    if err != nil {
        c.JSON(500, pb.PaidRes{Success: false, Message: err.Error()})
		log.Println("error while paid checkout", err)
        return
    }

    c.JSON(200, res)
}

// @Router        /api/order/history [GET]
// @Summary       GET order history
// @Description   This API retrieves an order history
// @Tags          Order
// @Accept        json
// @Produce       json
// @Param         id path string true "Courier ID"
// @Success       200 {object} pb.GetCourierOrderHistoryResponse
// @Failure       400 {object} string "error": "error message"
func(h *HandlerStruct)HistoryOrder(c *gin.Context){
	var req pb.GetCourierOrderHistoryRequest

	courier_id := c.Query("courier_id")

	req.CourierId = courier_id

	res, err := h.Clients.OrderClient.HistoryOrder(c.Request.Context(), &req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, res)
}
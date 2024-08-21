package handler

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "gitlab.com/bahodirova/api-gateway/genproto/product"
)


// @Router 				/api/cart [POST]
// @Summary 			CREATE cart
// @Description 		This API creates a new cart.
// @Tags 				Cart
// @Accept 				json
// @Produce 			json
// @Security            BearerAuth
// @Param 				data body pb.CreateCartReq true "Cart data"
// @Success 			200 {object} string "message": "created successfully"
// @Failure 			400 {object} string "error": "error message"
func (h *HandlerStruct) CreateCart(c *gin.Context) {
	var req pb.CreateCartReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	input, err := json.Marshal(req)
	err = h.Clients.KafkaProducer.ProduceMessages("create-cart", input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		log.Println("cannot produce messages via kafka", err)
		return
	}

	// _, err := h.Clients.CartClient.CreateCart(context.Background(), &req)
	// if err != nil {
	// 	c.JSON(400, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	c.JSON(200, gin.H{"message": "created successfully"})
}

// @Router 				/api/cart/{id} [GET]
// @Summary 			GET cart by ID
// @Description 		Retrieve a specific cart by its ID.
// @Tags 				Cart
// @Accept 				json
// @Produce 			json
// @Param 				id path string true "Cart ID"
// @Success 			200 {object} pb.Cart
// @Failure 			400 {object} string "error": "error message"
func (h *HandlerStruct) GetCart(c *gin.Context) {
	var req pb.GetById
	id := c.Param("id")
	req.Id = id

	res, err := h.Clients.CartClient.GetCart(context.Background(), &req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, res)
}

// @Router 				/api/cart/all [GET]
// @Summary 			GET all carts
// @Description 		Retrieve all carts with optional filters.
// @Tags 				Cart
// @Accept 				json
// @Produce 			json
// @Param 				quantity query int false "Quantity"
// @Param 				limit query int false "Limit"
// @Param 				offset query int false "Offset"
// @Success 			200 {object} pb.GetAllCartsRes
// @Failure 			400 {object} string "error": "error message"
func (h *HandlerStruct) GetAllCarts(c *gin.Context) {
	var req pb.GetAllCartsReq
	quantitySTR := c.Query("quantity")
	quantity, err := strconv.Atoi(quantitySTR)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid quantity value",
		})
		return
	}
	limit := c.Query("limit")
	offset := c.Query("offset")

	req.Quantity = int32(quantity)

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
	res, err := h.Clients.CartClient.GetAllCarts(context.Background(), &req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, res)
}

// @Router 				/api/cart [PUT]
// @Summary 			UPDATE cart
// @Description 		This API updates an existing cart.
// @Tags 				Cart
// @Accept 				json
// @Produce 			json
// @Security            BearerAuth
// @Param 				data body pb.UpdateCartReq true "Cart data"
// @Success 			200 {object} string "message": "updated successfully"
// @Failure 			400 {object} string "error": "error message"
func (h *HandlerStruct) UpdateCart(c *gin.Context) {
	var req pb.UpdateCartReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	input, err := json.Marshal(req)
	err = h.Clients.KafkaProducer.ProduceMessages("update-cart", input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		log.Println("cannot produce messages via kafka", err)
		return
	}

	// _, err := h.Clients.CartClient.UpdateCart(context.Background(), &req)
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

// @Router 				/api/cart/{id} [DELETE]
// @Summary 			DELETE cart
// @Description 		This API deletes a cart by its ID.
// @Tags 				Cart
// @Accept 				json
// @Produce 			json
// @Param 				id path string true "Cart ID"
// @Success 			200 {object} string "message": "deleted successfully"
// @Failure 			400 {object} string "error": "error message"
func (h *HandlerStruct) DeleteCart(c *gin.Context) {
	var req pb.GetById
	id := c.Param("id")
	req.Id = id

	input, err := json.Marshal(req)
	err = h.Clients.KafkaProducer.ProduceMessages("delete-cart", input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		log.Println("cannot produce messages via kafka", err)
		return
	}

	// _, err := h.Clients.CartClient.DeleteCart(context.Background(), &req)
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

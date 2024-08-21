package handler

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "gitlab.com/bahodirova/api-gateway/genproto/product"
)

// @Router        /api/notification [POST]
// @Summary       CREATE notification
// @Description   This API creates a notification
// @Tags          Notification
// @Accept        json
// @Produce       json
// @Security      BearerAuth
// @Param         data body pb.CreateNotificationReq true "Notification"
// @Success       200 {object} string "message": "created successfully"
// @Failure       400 {object} string "error": "error message"
func (h *HandlerStruct) CreateNotification(c *gin.Context) {
	var req pb.CreateNotificationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	input, err := json.Marshal(req)
	err = h.Clients.KafkaProducer.ProduceMessages("create-notification", input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		log.Println("cannot produce messages via kafka", err)
		return
	}

	// _, err := h.Clients.NotificationClient.CreateNotification(context.Background(), &req)
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

// @Router        /api/notification/{id} [GET]
// @Summary       GET notification
// @Description   This API retrieves a notification by ID
// @Tags          Notification
// @Accept        json
// @Produce       json
// @Param         id path string true "Notification ID"
// @Success       200 {object} pb.Notification
// @Failure       400 {object} string "error": "error message"
func (h *HandlerStruct) GetNotification(c *gin.Context) {
	var req pb.GetById
	id := c.Param("id")
	req.Id = id
	res, err := h.Clients.NotificationClient.GetNotification(context.Background(), &req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, res)
}

// @Router        /api/notification/all [GET]
// @Summary       GET all notifications
// @Description   This API retrieves all notifications with optional filters
// @Tags          Notification
// @Accept        json
// @Produce       json
// @Param         user_id query string false "User ID"
// @Param         limit query string false "Limit"
// @Param         offset query string false "Offset"
// @Success       200 {object} pb.GetAllNotificationsRes
// @Failure       400 {object} string "error": "error message"
func (h *HandlerStruct) GetAllNotifications(c *gin.Context) {
	var req pb.GetAllNotificationsReq
	user := c.Query("user_id")
	limit := c.Query("limit")
	offset := c.Query("offset")

	req.UserId = user

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
	res, err := h.Clients.NotificationClient.GetAllNotifications(context.Background(), &req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, res)
}

// @Router        /api/notification/read [PUT]
// @Summary       MARK notification as read
// @Description   This API marks a notification as read
// @Tags          Notification
// @Accept        json
// @Produce       json
// @Security      BearerAuth
// @Param         data body pb.MarkNotificationAsReadReq true "Notification"
// @Success       200 {object} string "message": "updated successfully"
// @Failure       400 {object} string "error": "error message"
func (h *HandlerStruct) MarkNotificationAsRead(c *gin.Context) {
	var req pb.MarkNotificationAsReadReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	_, err := h.Clients.NotificationClient.MarkNotificationAsRead(context.Background(), &req)
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

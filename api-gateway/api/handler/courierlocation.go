package handler

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "gitlab.com/bahodirova/api-gateway/genproto/product"
)

// @Router        /api/courier-location [POST]
// @Summary       CREATE courier location
// @Description   This API creates a courier location
// @Tags          CourierLocation
// @Accept        json
// @Produce       json
// @Security      BearerAuth
// @Param         data body pb.CreateCourierLocationRequest true "Courier Location"
// @Success       200 {object} string "message": "created successfully"
// @Failure       400 {object} string "error": "error message"
func (h *HandlerStruct) CreateCourierLocation(c *gin.Context) {
	var req pb.CreateCourierLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err := h.Clients.CourierLocationClient.CreateCourierLocation(context.Background(), &req)
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

// @Router        /api/courier-location/{id} [GET]
// @Summary       GET courier location
// @Description   This API retrieves a courier location by ID
// @Tags          CourierLocation
// @Accept        json
// @Produce       json
// @Param         id path string true "Courier Location ID"
// @Success       200 {object} pb.CourierLocation
// @Failure       400 {object} string "error": "error message"
func (h *HandlerStruct) GetCourierLocation(c *gin.Context) {
	var req pb.GetById
	id := c.Param("id")
	req.Id = id
	res, err := h.Clients.CourierLocationClient.GetCourierLocation(context.Background(), &req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, res)
}

// @Router        /api/courier-location/all [GET]
// @Summary       GET all courier locations
// @Description   This API retrieves all courier locations with optional filters
// @Tags          CourierLocation
// @Accept        json
// @Produce       json
// @Param         courier_id query string false "Courier ID"
// @Param         limit query string false "Limit"
// @Param         offset query string false "Offset"
// @Success       200 {object} pb.GetAllCourierLocationsRes
// @Failure       400 {object} string "error": "error message"
func (h *HandlerStruct) GetAllCourierLocation(c *gin.Context) {
	var req pb.GetAllCourierLocationsReq
	courier := c.Query("courier_id")
	limit := c.Query("limit")
	offset := c.Query("offset")

	req.CourierId = courier

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
	res, err := h.Clients.CourierLocationClient.GetAllCourierLocations(context.Background(), &req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, res)
}

// @Router        /api/courier-location [PUT]
// @Summary       UPDATE courier location
// @Description   This API updates a courier location
// @Tags          CourierLocation
// @Accept        json
// @Produce       json
// @Security      BearerAuth
// @Param         data body pb.UpdateCourierLocationRequest true "Courier Location"
// @Success       200 {object} string "message": "updated successfully"
// @Failure       400 {object} string "error": "error message"
func (h *HandlerStruct) UpdateCourierLocation(c *gin.Context) {
	var req pb.UpdateCourierLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	_, err := h.Clients.CourierLocationClient.UpdateCourierLocation(context.Background(), &req)
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

// @Router        /api/courier-location/by-time-range [GET]
// @Summary       GET courier locations by time range
// @Description   This API retrieves courier locations within a specific time range
// @Tags          CourierLocation
// @Accept        json
// @Produce       json
// @Param         courier_id query string true "Courier ID"
// @Param         start_time query string true "Start Time"
// @Param         end_time query string true "End Time"
// @Param         limit query string false "Limit"
// @Param         offset query string false "Offset"
// @Success       200 {object} pb.GetCourierLocationsByTimeRangeRes
// @Failure       400 {object} string "error": "error message"
func (h *HandlerStruct) GetCourierLocationsByTimeRange(c *gin.Context) {
	var req pb.GetCourierLocationsByTimeRangeReq
	courier := c.Query("courier_id")
	start_time := c.Query("start_time")
	end_time := c.Query("end_time")
	limit := c.Query("limit")
	offset := c.Query("offset")

	req.CourierId = courier
	req.StartTime = start_time
	req.EndTime = end_time

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
	res, err := h.Clients.CourierLocationClient.GetCourierLocationsByTimeRange(context.Background(), &req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, res)
}

// @Router        /api/courier-location/status [PUT]
// @Summary       UPDATE courier location status
// @Description   This API updates the status of a courier location
// @Tags          CourierLocation
// @Accept        json
// @Produce       json
// @Security      BearerAuth
// @Param         data body pb.UpdateCourierLocationStatusReq true "Courier Location Status"
// @Success       200 {object} string "message": "updated successfully"
// @Failure       400 {object} string "error": "error message"
func (h *HandlerStruct) UpdateCourierLocationStatus(c *gin.Context) {
	var req pb.UpdateCourierLocationStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	_, err := h.Clients.CourierLocationClient.UpdateCourierLocationStatus(context.Background(), &req)
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

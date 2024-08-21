package handler

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "gitlab.com/bahodirova/api-gateway/genproto/product"
)

// @Router        /api/task [POST]
// @Summary       Create a new task
// @Description   This API creates a new task
// @Tags          Task
// @Accept        json
// @Produce       json
// @Param         data body pb.CreatetaskReq true "Task"
// @Success       200 {object} string "message": "Created successfully"
// @Failure       400 {object} string "error": "error message"
// @Failure       500 {object} string "error": "error message"
func (h *HandlerStruct) CreateTask(c *gin.Context) {
	var req pb.CreatetaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	input, err := json.Marshal(req)
	err = h.Clients.KafkaProducer.ProduceMessages("create-task", input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		log.Println("cannot produce messages via kafka", err)
		return
	}

	// _, err := h.Clients.TaskClient.CreateTask(context.Background(), &req)
	// if err != nil {
	// 	c.JSON(400, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	c.JSON(200, gin.H{
		"message": "Created successfully",
	})
}

// @Router        /api/task/{id} [GET]
// @Summary       Get a task by ID
// @Description   This API retrieves a task by its ID
// @Tags          Task
// @Accept        json
// @Produce       json
// @Param         id path string true "Task ID"
// @Success       200 {object} pb.Task
// @Failure       400 {object} string "error": "error message"
func (h *HandlerStruct) GetTask(c *gin.Context) {
	var req pb.GetById
	id := c.Param("id")
	req.Id = id
	res, err := h.Clients.TaskClient.GetTask(context.Background(), &req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, res)
}

// @Router        /api/task/all [GET]
// @Summary       Get all tasks
// @Description   This API retrieves all tasks with optional filters
// @Tags          Task
// @Accept        json
// @Produce       json
// @Param         status query string false "Task Status"
// @Param         date query string false "Task Date"
// @Param         limit query int false "Limit"
// @Param         offset query int false "Offset"
// @Success       200 {object} pb.GetAllTasksRes
// @Failure       400 {object} string "error": "error message"
func (h *HandlerStruct) GetAllTasks(c *gin.Context) {
	var req pb.GetAllTasksReq
	status := c.Query("status")
	date := c.Query("date")
	limit := c.Query("limit")
	offset := c.Query("offset")

	req.Status = status
	req.Date = date

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
	res, err := h.Clients.TaskClient.GetAllTasks(context.Background(), &req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, res)
}

// @Router        /api/task [PUT]
// @Summary       Update a task
// @Description   This API updates a task
// @Tags          Task
// @Accept        json
// @Produce       json
// @Param         data body pb.UpdateTaskReq true "Task"
// @Success       200 {object} string "message": "Updated successfully"
// @Failure       400 {object} string "error": "error message"
// @Failure       500 {object} string "error": "error message"
func (h *HandlerStruct) UpdateTask(c *gin.Context) {
	var req pb.UpdateTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	input, err := json.Marshal(req)
	err = h.Clients.KafkaProducer.ProduceMessages("update-task", input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		log.Println("cannot produce messages via kafka", err)
		return
	}

	// _, err := h.Clients.TaskClient.UpdateTask(context.Background(), &req)
	// if err != nil {
	// 	c.JSON(400, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	c.JSON(200, gin.H{
		"message": "Updated successfully",
	})
}

// @Router        /api/task/{id} [DELETE]
// @Summary       Delete a task by ID
// @Description   This API deletes a task by its ID
// @Tags          Task
// @Accept        json
// @Produce       json
// @Param         id path string true "Task ID"
// @Success       200 {object} string "message": "Deleted successfully"
// @Failure       400 {object} string "error": "error message"
// @Failure       500 {object} string "error": "error message"
func (h *HandlerStruct) DeleteTask(c *gin.Context) {
	var req pb.GetById
	id := c.Param("id")
	req.Id = id

	input, err := json.Marshal(req)
	err = h.Clients.KafkaProducer.ProduceMessages("delete-task", input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		log.Println("cannot produce messages via kafka", err)
		return
	}

	// _, err := h.Clients.TaskClient.DeleteTask(context.Background(), &req)
	// if err != nil {
	// 	c.JSON(400, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	c.JSON(200, gin.H{
		"message": "Deleted successfully",
	})
}

// @Router        /api/task/user [GET]
// @Summary       Get tasks by user ID
// @Description   This API retrieves tasks assigned to a specific user
// @Tags          Task
// @Accept        json
// @Produce       json
// @Param         user_id_assigned_to query int true "User ID assigned to"
// @Param         limit query int false "Limit"
// @Param         offset query int false "Offset"
// @Success       200 {object} pb.GetAllTasksRes
// @Failure       400 {object} string "error": "error message"
func (h *HandlerStruct) GetTasksByUser(c *gin.Context) {
	var req pb.GetByUserReq
	userSTR := c.Query("user_id_assigned_to")
	user, err := strconv.Atoi(userSTR)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user value"})
		return
	}

	limit := c.Query("limit")
	offset := c.Query("offset")

	req.UserIdAssignedTo = int32(user)

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
	res, err := h.Clients.TaskClient.GetTasksByUser(context.Background(), &req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, res)
}

// @Router        /api/task/search [GET]
// @Summary       Search tasks
// @Description   This API searches tasks based on filters
// @Tags          Task
// @Accept        json
// @Produce       json
// @Param         title query string false "Title"
// @Param         description query string false "Description"
// @Param         status query string false "Status"
// @Param         date query string false "Date"
// @Param         limit query int false "Limit"
// @Param         offset query int false "Offset"
// @Success       200 {object} pb.GetAllTasksRes
// @Failure       400 {object} string "error": "error message"
func (h *HandlerStruct) SearchTasks(c *gin.Context) {
	var req pb.SearchTasksReq
	title := c.Query("title")
	description := c.Query("description")
	status := c.Query("status")
	date := c.Query("date")

	limit := c.Query("limit")
	offset := c.Query("offset")

	req.Title = title
	req.Description = description
	req.Status = status
	req.Date = date

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
	res, err := h.Clients.TaskClient.SearchTasks(context.Background(), &req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, res)
}

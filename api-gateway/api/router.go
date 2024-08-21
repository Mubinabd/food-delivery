package api

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"gitlab.com/bahodirova/api-gateway/api/handler"
	"gitlab.com/bahodirova/api-gateway/api/middlerware"
	_ "gitlab.com/bahodirova/api-gateway/docs"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @BasePath  /router
// @description					Description for what is this security definition being used
func NewRouter(h *handler.HandlerStruct) *gin.Engine {
	router := gin.Default()
	enforcer, err := casbin.NewEnforcer("./config/model.conf", "./config/policy.csv")

	if err != nil {
		log.Fatal(err)
	}
	sw := router.Group("/")
	sw.Use(middlerware.NewAuth(enforcer))

	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost", "http://localhost:8090"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}
	router.Use(cors.New(corsConfig))

	handler := handler.NewHandlerStruct()
	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api")
	{
		// cart routes
		cart := v1.Group("/cart")
		{
			cart.POST("", handler.CreateCart)
			cart.GET(":id", handler.GetCart)
			cart.GET("/all", handler.GetAllCarts)
			cart.PUT("", handler.UpdateCart)
			cart.DELETE(":id", handler.DeleteCart)
		}

		// order item routes
		orderitem := v1.Group("/order-item")
		{
			orderitem.POST("", handler.CreateOrderItem)
			orderitem.GET(":id", handler.GetOrderItem)
			orderitem.GET("/all", handler.GetAllOrderItems)
			orderitem.PUT("", handler.UpdateOrderItem)
			orderitem.GET("/order", handler.GetOrderItemsByOrder)
			orderitem.GET("/product", handler.GetOrderItemsByProduct)
		}

		// Notification routes
		notification := v1.Group("/notification")
		{
			notification.POST("", handler.CreateNotification)
			notification.GET(":id", handler.GetNotification)
			notification.GET("/all", handler.GetAllNotifications)
			notification.PUT("/read", handler.MarkNotificationAsRead)
		}

		// Task routes
		task := v1.Group("/task")
		{
			task.POST("", handler.CreateTask)
			task.GET(":id", handler.GetTask)
			task.GET("/all", handler.GetAllTasks)
			task.PUT("", handler.UpdateTask)
			task.DELETE(":id", handler.DeleteTask)
			task.GET("/user", handler.GetTasksByUser)
			task.GET("/search", handler.SearchTasks)
		}

		// Order routes
		order := v1.Group("/order")
		{
			order.POST("", handler.CreateOrder)
			order.GET(":id", handler.GetOrder)
			order.GET("/all", handler.GetAllOrders)
			order.PUT("", handler.UpdateOrder)
			order.DELETE(":id", handler.DeleteOrder)
			order.GET("/courier/history", handler.HistoryOrder)
		}

		//order paid for product
		orders := v1.Group("/orders")
		{
			orders.POST("/paid", handler.PaidOrder)

		}

		// Courier location routes
		courierlocation := v1.Group("/courier-location")
		{
			courierlocation.POST("", handler.CreateCourierLocation)
			courierlocation.GET(":id", handler.GetCourierLocation)
			courierlocation.GET("/all", handler.GetAllCourierLocation)
			courierlocation.PUT("", handler.UpdateCourierLocation)
			courierlocation.GET("/by-time-range", handler.GetCourierLocationsByTimeRange)
			courierlocation.PUT("/status", handler.UpdateCourierLocationStatus)
		}

		// Product routes
		pr := v1.Group("/product")
		{
			pr.POST("", handler.CreateProduct)
			pr.GET(":id", handler.GetProduct)
			pr.GET("", handler.GetAllProducts)
			pr.PUT("", handler.UpdateProduct)
			pr.DELETE(":id", handler.DeleteProduct)
			pr.GET("/search", handler.SearchProducts)
		}
	}
	file := router.Group("/minio/upload")
	{
		file.POST("", handler.UploadFile)

	}

	return router
}

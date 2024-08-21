package main

import (
	"log"
	"net"

	"gitlab.com/bahodirova/product/config"
	pb "gitlab.com/bahodirova/product/genproto/product"
	"gitlab.com/bahodirova/product/kafka"
	"gitlab.com/bahodirova/product/service"
	"gitlab.com/bahodirova/product/storage/postgres"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	db, err := postgres.NewPostgresStorage()
	if err != nil {
		log.Fatalf("error while connecting to postgres: %v", err)
	}

	lis, err := net.Listen("tcp", cfg.HTTPPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	productService := service.NewProductService(db)
	orderService := service.NewOrderService(db)
	taskService := service.NewTaskService(db)
	notificationService := service.NewNotificationService(db)
	cartService := service.NewCartService(db)

	brokers := []string{"kafka:9092"}

	kcm := kafka.NewKafkaConsumerManager()

	if err := kcm.RegisterConsumer(brokers, "create-product", "product", kafka.ProductHandler(productService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'create-product' already exists")
		} else {
			log.Printf("Failed to register consumer for topic 'create-product': %v", err)

		}
	}
	if err := kcm.RegisterConsumer(brokers, "update-product", "product", kafka.UpdateProductHandler(productService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'update-product' already exists")
		} else {
			log.Printf("Failed to register consumer for topic 'update-product': %v", err)

		}
	}
	if err := kcm.RegisterConsumer(brokers, "delete-product", "product", kafka.DeleteProductHandler(productService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'delete-product' already exists")
		} else {
			log.Printf("Failed to register consumer for topic 'update-product': %v", err)

		}
	}

	// Cart
	if err := kcm.RegisterConsumer(brokers, "create-cart", "cart", kafka.CartHandler(cartService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'create-cart' already exists")
		} else {
			log.Printf("Failed to register consumer for topic 'create-cart': %v", err)

		}
	}
	if err := kcm.RegisterConsumer(brokers, "delete-cart", "cart", kafka.DeleteCartHandler(cartService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'delete-cart' already exists")
		} else {
			log.Printf("Failed to register consumer for topic 'delete-cart': %v", err)

		}
	}
	if err := kcm.RegisterConsumer(brokers, "update-cart", "cart", kafka.UpdateCartHandler(cartService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'update-cart' already exists")
		} else {
			log.Printf("Failed to register consumer for topic 'update-cart': %v", err)

		}
	}

	// Notification
	if err := kcm.RegisterConsumer(brokers, "create-notification", "notification", kafka.NotificationHandler(notificationService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'create-notification' already exists")
		} else {
			log.Printf("Failed to register consumer for topic 'create-notification': %v", err)

		}
	}
	

	// Order
	if err := kcm.RegisterConsumer(brokers, "create-order", "order", kafka.OrderHandler(orderService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'create-order' already exists")
		} else {
			log.Printf("Failed to register consumer for topic 'create-order': %v", err)

		}
	}
	if err := kcm.RegisterConsumer(brokers, "delete-order", "order", kafka.DeleteOrderHandler(orderService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'delete-order' already exists")
		} else {
			log.Printf("Failed to register consumer for topic 'delete-order': %v", err)

		}
	}
	if err := kcm.RegisterConsumer(brokers, "update-order", "order", kafka.UpdateOrderHandler(orderService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'update-order' already exists")
		} else {
			log.Printf("Failed to register consumer for topic 'update-order': %v", err)

		}
	}

	// Task
	if err := kcm.RegisterConsumer(brokers, "create-task", "task", kafka.TaskHandler(taskService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'create-task' already exists")
		} else {
			log.Printf("Failed to register consumer for topic 'create-task': %v", err)

		}
	}
	if err := kcm.RegisterConsumer(brokers, "delete-task", "task", kafka.DeleteTaskHandler(taskService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'delete-task' already exists")
		} else {
			log.Printf("Failed to register consumer for topic 'delete-task': %v", err)

		}
	}
	if err := kcm.RegisterConsumer(brokers, "update-task", "task", kafka.UpdateTaskHandler(taskService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'update-task' already exists")
		} else {
			log.Printf("Failed to register consumer for topic 'update-task': %v", err)

		}
	}

	s := grpc.NewServer()

	pb.RegisterCartServiceServer(s, service.NewCartService(db))
	pb.RegisterCourierLocationServiceServer(s, service.NewCourierLocationService(db))
	pb.RegisterProductServiceServer(s, service.NewProductService(db))
	pb.RegisterTaskServiceServer(s, service.NewTaskService(db))
	pb.RegisterOrderItemServiceServer(s, service.NewOrderItemService(db))
	pb.RegisterOrderServiceServer(s, service.NewOrderService(db))
	pb.RegisterNotificationServiceServer(s, service.NewNotificationService(db))

	reflection.Register(s)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

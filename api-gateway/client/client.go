package client

import (
	"log"

	"github.com/minio/minio-go"
	pbc "gitlab.com/bahodirova/api-gateway/genproto/product"
	"gitlab.com/bahodirova/api-gateway/kafka"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Clients contains all the clients including MinIOClient
type Clients struct {
	ProductClient         pbc.ProductServiceClient
	CartClient            pbc.CartServiceClient
	OrderClient           pbc.OrderServiceClient
	OrderItemClient       pbc.OrderItemServiceClient
	NotificationClient    pbc.NotificationServiceClient
	CourierLocationClient pbc.CourierLocationServiceClient
	TaskClient            pbc.TaskServiceClient
	KafkaProducer         kafka.KafkaProducer
	MinIOClient           *minio.Client
}

func InitializeMinIOClient() (*minio.Client, error) {
	minIOClient, err := minio.New("localhost:9000", "root", "mubina0804", false)
	if err != nil {
		return nil, err
	}
	return minIOClient, nil
}

func NewClients() *Clients {
	conn, err := grpc.NewClient("delivery_service:8085", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server: %v", err)
	}

	productS := pbc.NewProductServiceClient(conn)
	cartS := pbc.NewCartServiceClient(conn)
	orderS := pbc.NewOrderServiceClient(conn)
	orderItemS := pbc.NewOrderItemServiceClient(conn)
	notificationS := pbc.NewNotificationServiceClient(conn)
	courierLocationS := pbc.NewCourierLocationServiceClient(conn)
	taskS := pbc.NewTaskServiceClient(conn)

	kafkaProducer, err := kafka.NewKafkaProducer([]string{"kafka:9092"})
	if err != nil {
		log.Fatalf("failed to create Kafka producer: %v", err)
	}

	minIOClient, err := InitializeMinIOClient()
	if err != nil {
		log.Fatalf("failed to initialize MinIO client: %v", err)
	}

	return &Clients{
		ProductClient:         productS,
		CartClient:            cartS,
		OrderClient:           orderS,
		OrderItemClient:       orderItemS,
		NotificationClient:    notificationS,
		CourierLocationClient: courierLocationS,
		TaskClient:            taskS,
		KafkaProducer:         kafkaProducer,
		MinIOClient:           minIOClient,
	}
}
package kafka

import (
	"context"
	"encoding/json"
	"log"

	pb "gitlab.com/bahodirova/product/genproto/product"
	"gitlab.com/bahodirova/product/service"
)

func NotificationHandler(Notificationservice *service.NotificationService) func(message []byte) {
	return func(message []byte) {
		var Notification pb.CreateNotificationReq
		if err := json.Unmarshal(message, &Notification); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}

		respNotification, err := Notificationservice.CreateNotification(context.Background(), &Notification)
		if err != nil {
			log.Printf("Cannot create Notification via Kafka: %v", err)
			return
		}
		log.Printf("Created Notification: %+v",respNotification)
	}
}
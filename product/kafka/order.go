package kafka

import (
	"context"
	"encoding/json"
	"log"

	pb "gitlab.com/bahodirova/product/genproto/product"
	"gitlab.com/bahodirova/product/service"
)

func OrderHandler(Orderservice *service.OrderService) func(message []byte) {
	return func(message []byte) {
		var order pb.CreateOrderReq
		if err := json.Unmarshal(message, &order); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}

		respOrder, err := Orderservice.CreateOrder(context.Background(), &order)
		if err != nil {
			log.Printf("Cannot create order via Kafka: %v", err)
			return
		}
		log.Printf("Created order: %+v",respOrder)
	}
}
func DeleteOrderHandler(Orderservice *service.OrderService) func(message []byte) {
	return func(message []byte) {
		var Order pb.GetById
		if err := json.Unmarshal(message, &Order); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}

		respOrder, err := Orderservice.DeleteOrder(context.Background(), &Order)
		if err != nil {
			log.Printf("Cannot delete Order via Kafka: %v", err)
			return
		}
		log.Printf("Deleted Order: %+v",respOrder)
	}
}
func UpdateOrderHandler(Orderservice *service.OrderService) func(message []byte) {
	return func(message []byte) {
		var Order pb.UpdateOrderReq
		if err := json.Unmarshal(message, &Order); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}

		respOrder, err := Orderservice.UpdateOrder(context.Background(), &Order)
		if err != nil {
			log.Printf("Cannot update Order via Kafka: %v", err)
			return
		}
		log.Printf("Updated Order: %+v",respOrder)
	}
}
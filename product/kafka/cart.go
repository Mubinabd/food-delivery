package kafka

import (
	"context"
	"encoding/json"
	"log"

	pb "gitlab.com/bahodirova/product/genproto/product"
	"gitlab.com/bahodirova/product/service"
)

func CartHandler(cartservice *service.CartService) func(message []byte) {
	return func(message []byte) {
		var cart pb.CreateCartReq
		if err := json.Unmarshal(message, &cart); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}

		respCart, err := cartservice.CreateCart(context.Background(), &cart)
		if err != nil {
			log.Printf("Cannot create cart via Kafka: %v", err)
			return
		}
		log.Printf("Created cart: %+v",respCart)
	}
}
func DeleteCartHandler(cartservice *service.CartService) func(message []byte) {
	return func(message []byte) {
		var cart pb.GetById
		if err := json.Unmarshal(message, &cart); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}

		respCart, err := cartservice.DeleteCart(context.Background(), &cart)
		if err != nil {
			log.Printf("Cannot delete cart via Kafka: %v", err)
			return
		}
		log.Printf("Deleted cart: %+v",respCart)
	}
}
func UpdateCartHandler(cartservice *service.CartService) func(message []byte) {
	return func(message []byte) {
		var cart pb.UpdateCartReq
		if err := json.Unmarshal(message, &cart); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}

		respCart, err := cartservice.UpdateCart(context.Background(), &cart)
		if err != nil {
			log.Printf("Cannot update cart via Kafka: %v", err)
			return
		}
		log.Printf("Updated cart: %+v",respCart)
	}
}
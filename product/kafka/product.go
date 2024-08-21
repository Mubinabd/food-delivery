package kafka

import (
	"context"
	"encoding/json"
	"log"

	pb "gitlab.com/bahodirova/product/genproto/product"
	"gitlab.com/bahodirova/product/service"
)

func ProductHandler(productService *service.ProductService) func(message []byte) {
	return func(message []byte) {
		var product pb.CreateProductRequest
		if err := json.Unmarshal(message, &product); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}

		respProduct, err := productService.CreateProduct(context.Background(), &product)
		if err != nil {
			log.Printf("Cannot create product via Kafka: %v", err)
			return
		}
		log.Printf("Created product: %+v",respProduct)
	}
}
func DeleteProductHandler(productService *service.ProductService) func(message []byte) {
	return func(message []byte) {
		var product pb.GetById
		if err := json.Unmarshal(message, &product); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}

		respProduct, err := productService.DeleteProduct(context.Background(), &product)
		if err != nil {
			log.Printf("Cannot delete product via Kafka: %v", err)
			return
		}
		log.Printf("Deleted product: %+v",respProduct)
	}
}
func UpdateProductHandler(productService *service.ProductService) func(message []byte) {
	return func(message []byte) {
		var product pb.UpdateProductRequest
		if err := json.Unmarshal(message, &product); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}

		respProduct, err := productService.UpdateProduct(context.Background(), &product)
		if err != nil {
			log.Printf("Cannot update product via Kafka: %v", err)
			return
		}
		log.Printf("Updated product: %+v",respProduct)
	}
}
package kafka

import (
	"context"
	"encoding/json"
	"log"

	pb "gitlab.com/bahodirova/product/genproto/product"
	"gitlab.com/bahodirova/product/service"
)

func TaskHandler(Taskservice *service.TaskService) func(message []byte) {
	return func(message []byte) {
		var task pb.CreatetaskReq
		if err := json.Unmarshal(message, &task); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}

		respTask, err := Taskservice.CreateTask(context.Background(), &task)
		if err != nil {
			log.Printf("Cannot create Task via Kafka: %v", err)
			return
		}
		log.Printf("Created Task: %+v",respTask)
	}
}
func DeleteTaskHandler(Taskservice *service.TaskService) func(message []byte) {
	return func(message []byte) {
		var Task pb.GetById
		if err := json.Unmarshal(message, &Task); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}

		respTask, err := Taskservice.DeleteTask(context.Background(), &Task)
		if err != nil {
			log.Printf("Cannot delete Task via Kafka: %v", err)
			return
		}
		log.Printf("Deleted Task: %+v",respTask)
	}
}
func UpdateTaskHandler(Taskservice *service.TaskService) func(message []byte) {
	return func(message []byte) {
		var Task pb.UpdateTaskReq
		if err := json.Unmarshal(message, &Task); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}

		respTask, err := Taskservice.UpdateTask(context.Background(), &Task)
		if err != nil {
			log.Printf("Cannot update Task via Kafka: %v", err)
			return
		}
		log.Printf("Updated Task: %+v",respTask)
	}
}
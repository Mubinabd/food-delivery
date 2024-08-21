package service

import (
	"context"

	pb "gitlab.com/bahodirova/product/genproto/product"
	"gitlab.com/bahodirova/product/storage"
)

type TaskService struct {
	storage storage.StorageI
	pb.UnimplementedTaskServiceServer
}

func NewTaskService(storage storage.StorageI) *TaskService {
	return &TaskService{
		storage: storage,
	}
}

func (a *TaskService) CreateTask(ctx context.Context, req *pb.CreatetaskReq) (*pb.Empty, error) {
	return a.storage.Task().CreateTask(req)
}
func (a *TaskService) GetTask(ctx context.Context, req *pb.GetById) (*pb.Task, error) {
	return a.storage.Task().GetTask(req)
}
func (a *TaskService) GetAllTasks(ctx context.Context, req *pb.GetAllTasksReq) (*pb.GetAllTasksRes, error) {
	return a.storage.Task().GetAllTasks(req)
}
func (a *TaskService) DeleteTask(ctx context.Context, req *pb.GetById) (*pb.DeleteTaskRes, error) {
	return a.storage.Task().DeleteTask(req)
}
func (a *TaskService) UpdateTask(ctx context.Context, req *pb.UpdateTaskReq) (*pb.UpdateTaskRes, error) {
	return a.storage.Task().UpdateTask(req)
}
func (a *TaskService) GetTasksByUser(ctx context.Context, req *pb.GetByUserReq) (*pb.GetAllTasksRes, error) {
	return a.storage.Task().GetTasksByUser(req)
}
func (a *TaskService) SearchTasks(ctx context.Context, req *pb.SearchTasksReq) (*pb.GetAllTasksRes, error) {
	return a.storage.Task().SearchTasks(req)
}



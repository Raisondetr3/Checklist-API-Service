package dto

import (
	"github.com/Raisondetr3/checklist-api-service/internal/model"
	pb "github.com/Raisondetr3/checklist-api-service/pkg/pb"
)

func ProtoToTaskResponse(protoTask *pb.Task) TaskResponse {
	if protoTask == nil {
		return TaskResponse{}
	}

	return TaskResponse{
		ID:          protoTask.Id,
		Title:       protoTask.Title,
		Description: protoTask.Description,
		Completed:   protoTask.Completed,
		CreatedAt:   protoTask.CreatedAt.AsTime(),
		UpdatedAt:   protoTask.UpdatedAt.AsTime(),
	}
}

func ProtoToTaskListResponse(protoResp *pb.ListTasksResponse) TaskListResponse {
	if protoResp == nil || len(protoResp.Tasks) == 0 {
		return TaskListResponse{
			Tasks: []TaskResponse{},
		}
	}

	tasks := make([]TaskResponse, len(protoResp.Tasks))
	for i, protoTask := range protoResp.Tasks {
		tasks[i] = ProtoToTaskResponse(protoTask)
	}

	return TaskListResponse{
		Tasks: tasks,
	}
}

func CreateTaskRequestToProto(dto CreateTaskRequest) *pb.CreateTaskRequest {
	return &pb.CreateTaskRequest{
		Title:       dto.Title,
		Description: dto.Description,
	}
}

func UpdateTaskRequestToProto(id string, dto UpdateTaskRequest) *pb.UpdateTaskRequest {
	req := &pb.UpdateTaskRequest{
		Id: id,
	}

	if dto.Title != nil {
		req.Title = dto.Title
	}
	if dto.Description != nil {
		req.Description = dto.Description
	}
	if dto.Completed != nil {
		req.Completed = dto.Completed
	}

	return req
}

func GetTaskRequestToProto(id string) *pb.GetTaskRequest {
	return &pb.GetTaskRequest{
		Id: id,
	}
}

func DeleteTaskRequestToProto(id string) *pb.DeleteTaskRequest {
	return &pb.DeleteTaskRequest{
		Id: id,
	}
}

func ListTasksRequestToProto() *pb.ListTasksRequest {
	return &pb.ListTasksRequest{}
}

func TaskModelToResponse(task *model.Task) TaskResponse {
	if task == nil {
		return TaskResponse{}
	}

	return TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func TaskModelsToResponse(tasks []*model.Task) TaskListResponse {
	if tasks == nil {
		return TaskListResponse{
			Tasks: []TaskResponse{},
		}
	}

	responses := make([]TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = TaskModelToResponse(task)
	}

	return TaskListResponse{
		Tasks: responses,
	}
}

func CreateTaskRequestToModel(dto CreateTaskRequest) *model.Task {
	return model.NewTask(dto.Title, dto.Description)
}

func ProtoToModelTask(protoTask *pb.Task) *model.Task {
	if protoTask == nil {
		return nil
	}

	return &model.Task{
		ID:          protoTask.Id,
		Title:       protoTask.Title,
		Description: protoTask.Description,
		Completed:   protoTask.Completed,
		CreatedAt:   protoTask.CreatedAt.AsTime(),
		UpdatedAt:   protoTask.UpdatedAt.AsTime(),
	}
}

func ProtoToModelTasks(protoTasks []*pb.Task) []*model.Task {
	if protoTasks == nil {
		return nil
	}

	tasks := make([]*model.Task, len(protoTasks))
	for i, protoTask := range protoTasks {
		tasks[i] = ProtoToModelTask(protoTask)
	}

	return tasks
}
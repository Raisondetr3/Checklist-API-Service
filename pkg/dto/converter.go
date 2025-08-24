package dto

import (
	pb "github.com/Raisondetr3/checklist-api-service/proto"
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

// CreateTaskRequestToProto преобразует DTO CreateTaskRequest в protobuf
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

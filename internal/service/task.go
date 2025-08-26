package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Raisondetr3/checklist-api-service/internal/client"
	"github.com/Raisondetr3/checklist-api-service/internal/model"
	"github.com/Raisondetr3/checklist-api-service/pkg/dto"
	"github.com/Raisondetr3/checklist-api-service/pkg/logger"
)

type TaskService interface {
	CreateTask(ctx context.Context, task *model.Task) (*model.Task, error)
	GetTasks(ctx context.Context, completed *bool) ([]*model.Task, int, error)
	GetTask(ctx context.Context, taskID string) (*model.Task, error)
	UpdateTask(ctx context.Context, taskID string, title, description *string, completed *bool) (*model.Task, error)
	DeleteTask(ctx context.Context, taskID string) error
}

type taskService struct {
	grpcClient client.TaskClient
}

func NewTaskService(taskClient client.TaskClient) TaskService {
	return &taskService{
		grpcClient: taskClient,
	}
}

func (t *taskService) CreateTask(ctx context.Context, task *model.Task) (*model.Task, error) {
	start := time.Now()
	operation := "CreateTask"

	if task == nil {
		err := fmt.Errorf("task cannot be nil")
		logger.LogError(ctx, err, operation)
		return nil, err
	}

	if err := task.Validate(); err != nil {
		logger.LogError(ctx, err, operation)
		return nil, fmt.Errorf("task validation failed: %w", err)
	}

	createReq := dto.CreateTaskRequest{
		Title:       task.Title,
		Description: task.Description,
	}

	protoReq := dto.CreateTaskRequestToProto(createReq)

	protoResp, err := t.grpcClient.CreateTask(ctx, protoReq)
	duration := time.Since(start)

	if err != nil {
		logger.LogError(ctx, err, operation,
			slog.Duration("duration", duration),
			slog.String("title", task.Title),
		)
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	createdTask := dto.ProtoToModelTask(protoResp.Task)

	slog.InfoContext(ctx, "Task created successfully",
		slog.String("operation", operation),
		slog.String("task_id", createdTask.ID),
		slog.String("title", createdTask.Title),
		slog.Duration("duration", duration),
	)

	return createdTask, nil
}

func (t *taskService) GetTasks(ctx context.Context, completed *bool) ([]*model.Task, int, error) {
	start := time.Now()
	operation := "GetTasks"

	protoReq := dto.ListTasksRequestToProto()

	protoResp, err := t.grpcClient.ListTasks(ctx, protoReq)
	duration := time.Since(start)

	if err != nil {
		logger.LogError(ctx, err, operation,
			slog.Duration("duration", duration),
		)
		return nil, 0, fmt.Errorf("failed to get tasks: %w", err)
	}

	allTasks := dto.ProtoToModelTasks(protoResp.Tasks)

	var filteredTasks []*model.Task
	if completed != nil {
		filteredTasks = make([]*model.Task, 0)
		for _, task := range allTasks {
			if task.Completed == *completed {
				filteredTasks = append(filteredTasks, task)
			}
		}
	} else {
		filteredTasks = allTasks
	}

	totalCount := len(filteredTasks)

	slog.InfoContext(ctx, "Tasks retrieved successfully",
		slog.String("operation", operation),
		slog.Int("total_count", len(allTasks)),
		slog.Int("filtered_count", totalCount),
		slog.Duration("duration", duration),
	)

	return filteredTasks, totalCount, nil
}

func (t *taskService) GetTask(ctx context.Context, taskID string) (*model.Task, error) {
	start := time.Now()
	operation := "GetTask"

	if taskID == "" {
		err := fmt.Errorf("task ID is required")
		logger.LogError(ctx, err, operation)
		return nil, err
	}

	protoReq := dto.GetTaskRequestToProto(taskID)

	protoResp, err := t.grpcClient.GetTask(ctx, protoReq)
	duration := time.Since(start)

	if err != nil {
		logger.LogError(ctx, err, operation,
			slog.Duration("duration", duration),
			slog.String("task_id", taskID),
		)
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	task := dto.ProtoToModelTask(protoResp.Task)

	slog.InfoContext(ctx, "Task retrieved successfully",
		slog.String("operation", operation),
		slog.String("task_id", task.ID),
		slog.Duration("duration", duration),
	)

	return task, nil
}

func (t *taskService) UpdateTask(ctx context.Context, taskID string, title, description *string, completed *bool) (*model.Task, error) {
	start := time.Now()
	operation := "UpdateTask"

	if taskID == "" {
		err := fmt.Errorf("task ID is required")
		logger.LogError(ctx, err, operation)
		return nil, err
	}

	if title == nil && description == nil && completed == nil {
		err := fmt.Errorf("at least one field must be provided for update")
		logger.LogError(ctx, err, operation,
			slog.String("task_id", taskID),
		)
		return nil, err
	}
	
	updateReq := dto.UpdateTaskRequest{
		Title:       title,
		Description: description,
		Completed:   completed,
	}

	if title != nil && *title == "" {
		err := fmt.Errorf("title cannot be empty")
		logger.LogError(ctx, err, operation, slog.String("task_id", taskID))
		return nil, err
	}

	protoReq := dto.UpdateTaskRequestToProto(taskID, updateReq)

	protoResp, err := t.grpcClient.UpdateTask(ctx, protoReq)
	duration := time.Since(start)

	if err != nil {
		logger.LogError(ctx, err, operation,
			slog.Duration("duration", duration),
			slog.String("task_id", taskID),
		)
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	updatedTask := dto.ProtoToModelTask(protoResp.Task)

	slog.InfoContext(ctx, "Task updated successfully",
		slog.String("operation", operation),
		slog.String("task_id", updatedTask.ID),
		slog.Duration("duration", duration),
	)

	return updatedTask, nil
}

func (t *taskService) DeleteTask(ctx context.Context, taskID string) error {
	start := time.Now()
	operation := "DeleteTask"

	if taskID == "" {
		err := fmt.Errorf("task ID is required")
		logger.LogError(ctx, err, operation)
		return err
	}

	protoReq := dto.DeleteTaskRequestToProto(taskID)

	protoResp, err := t.grpcClient.DeleteTask(ctx, protoReq)
	duration := time.Since(start)

	if err != nil {
		logger.LogError(ctx, err, operation,
			slog.Duration("duration", duration),
			slog.String("task_id", taskID),
		)
		return fmt.Errorf("failed to delete task: %w", err)
	}

	if !protoResp.Success {
		err := fmt.Errorf("task deletion was not successful")
		logger.LogError(ctx, err, operation,
			slog.String("task_id", taskID),
		)
		return err
	}

	slog.InfoContext(ctx, "Task deleted successfully",
		slog.String("operation", operation),
		slog.String("task_id", taskID),
		slog.Duration("duration", duration),
	)

	return nil
}
package client

import (
	"checklist-api-service/internal/model"
	"context"
)

type DBClient interface {
	CreateTask(ctx context.Context, task *model.Task) (*model.Task, error)
	GetTasks(ctx context.Context, limit, offset int, completed *bool) ([]*model.Task, int, error)
	GetTask(ctx context.Context, taskID string) (*model.Task, error)
	UpdateTask(ctx context.Context, task *model.Task) (*model.Task, error)
	DeleteTask(ctx context.Context, taskID string) error
	Health(ctx context.Context) error
}

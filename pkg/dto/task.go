package dto

import "time"

const (
	EventTaskCreated = "task.created"
	EventTaskUpdated = "task.updated"
	EventTaskDeleted = "task.deleted"
	EventTaskViewed  = "task.viewed"
	EventTasksListed = "tasks.listed"
)

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TaskResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Completed   *bool   `json:"completed,omitempty"`
}

type TaskListResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}

type DeleteTaskResponse struct {
	Success bool `json:"success"`
}

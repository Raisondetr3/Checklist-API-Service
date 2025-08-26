package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/Raisondetr3/checklist-api-service/internal/service"
	"github.com/Raisondetr3/checklist-api-service/pkg/dto"
	apiErrors "github.com/Raisondetr3/checklist-api-service/pkg/errors"

	"github.com/gorilla/mux"
)

type TaskHandlers struct {
	taskService service.TaskService
}

func NewTaskHandlers(taskService service.TaskService) *TaskHandlers {
	return &TaskHandlers{
		taskService: taskService,
	}
}

func (h *TaskHandlers) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		h.writeErrorResponse(w, "Title is required", http.StatusBadRequest)
		return
	}

	if len(req.Title) > 255 {
		h.writeErrorResponse(w, "Title is too long (max 255 characters)", http.StatusBadRequest)
		return
	}

	if len(req.Description) > 1000 {
		h.writeErrorResponse(w, "Description is too long (max 1000 characters)", http.StatusBadRequest)
		return
	}

	task := dto.CreateTaskRequestToModel(req)

	createdTask, err := h.taskService.CreateTask(ctx, task)
	if err != nil {
		h.handleServiceError(w, err, "Failed to create task")
		return
	}

	response := dto.TaskModelToResponse(createdTask)

	h.writeJSONResponse(w, http.StatusCreated, response)
	
	slog.InfoContext(ctx, "Task created via HTTP",
		slog.String("task_id", response.ID),
		slog.String("title", response.Title),
	)
}

func (h *TaskHandlers) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var completed *bool
	if completedStr := r.URL.Query().Get("completed"); completedStr != "" {
		if completedBool, err := strconv.ParseBool(completedStr); err == nil {
			completed = &completedBool
		} else {
			h.writeErrorResponse(w, "Invalid 'completed' parameter. Use 'true' or 'false'", http.StatusBadRequest)
			return
		}
	}

	tasks, totalCount, err := h.taskService.GetTasks(ctx, completed)
	if err != nil {
		h.handleServiceError(w, err, "Failed to get tasks")
		return
	}

	response := dto.TaskModelsToResponse(tasks)

	h.writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"tasks":       response.Tasks,
		"total_count": totalCount,
	})

	slog.InfoContext(ctx, "Tasks retrieved via HTTP",
		slog.Int("count", len(tasks)),
		slog.Int("total_count", totalCount),
	)
}

func (h *TaskHandlers) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	taskID := vars["id"]

	if taskID == "" {
		h.writeErrorResponse(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	task, err := h.taskService.GetTask(ctx, taskID)
	if err != nil {
		h.handleServiceError(w, err, "Failed to get task")
		return
	}

	response := dto.TaskModelToResponse(task)

	h.writeJSONResponse(w, http.StatusOK, response)

	slog.InfoContext(ctx, "Task retrieved via HTTP",
		slog.String("task_id", response.ID),
	)
}

func (h *TaskHandlers) HandleUpdateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	taskID := vars["id"]

	if taskID == "" {
		h.writeErrorResponse(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	var req dto.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if req.Title == nil && req.Description == nil && req.Completed == nil {
		h.writeErrorResponse(w, "At least one field must be provided for update", http.StatusBadRequest)
		return
	}

	if req.Title != nil {
		if *req.Title == "" {
			h.writeErrorResponse(w, "Title cannot be empty", http.StatusBadRequest)
			return
		}
		if len(*req.Title) > 255 {
			h.writeErrorResponse(w, "Title is too long (max 255 characters)", http.StatusBadRequest)
			return
		}
	}

	if req.Description != nil && len(*req.Description) > 1000 {
		h.writeErrorResponse(w, "Description is too long (max 1000 characters)", http.StatusBadRequest)
		return
	}

	updatedTask, err := h.taskService.UpdateTask(ctx, taskID, req.Title, req.Description, req.Completed)
	if err != nil {
		h.handleServiceError(w, err, "Failed to update task")
		return
	}

	response := dto.TaskModelToResponse(updatedTask)

	h.writeJSONResponse(w, http.StatusOK, response)

	slog.InfoContext(ctx, "Task updated via HTTP",
		slog.String("task_id", response.ID),
	)
}

func (h *TaskHandlers) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	taskID := vars["id"]

	if taskID == "" {
		h.writeErrorResponse(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	err := h.taskService.DeleteTask(ctx, taskID)
	if err != nil {
		h.handleServiceError(w, err, "Failed to delete task")
		return
	}

	response := dto.DeleteTaskResponse{Success: true}
	h.writeJSONResponse(w, http.StatusOK, response)

	slog.InfoContext(ctx, "Task deleted via HTTP",
		slog.String("task_id", taskID),
	)
}

func (h *TaskHandlers) handleServiceError(w http.ResponseWriter, err error, defaultMessage string) {
	statusCode := apiErrors.HTTPStatusFromError(err)
	message := apiErrors.MessageFromError(err)

	if message == "" || strings.Contains(message, "rpc") {
		message = defaultMessage
	}

	errDTO := dto.NewErr(message)
	http.Error(w, errDTO.ToString(), statusCode)
}

func (h *TaskHandlers) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			slog.Error("Failed to encode JSON response", slog.String("error", err.Error()))
		}
	}
}

func (h *TaskHandlers) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	errDTO := dto.NewErr(message)
	http.Error(w, errDTO.ToString(), statusCode)

	slog.WarnContext(context.Background(), "HTTP error response sent",
		slog.Int("status_code", statusCode),
		slog.String("message", message),
	)
}

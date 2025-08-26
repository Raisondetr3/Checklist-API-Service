package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Raisondetr3/checklist-api-service/internal/service"
	"github.com/Raisondetr3/checklist-api-service/internal/validator"
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
		WriteErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if validationErrs := validator.ValidateCreateTaskRequest(req); validationErrs.HasErrors() {
		WriteErrorResponse(w, validationErrs.Error(), http.StatusBadRequest)
		return
	}

	task := dto.CreateTaskRequestToModel(req)

	createdTask, err := h.taskService.CreateTask(ctx, task)
	if err != nil {
		h.handleServiceError(w, err, "Failed to create task")
		return
	}

	response := dto.TaskModelToResponse(createdTask)

	WriteJSONResponse(w, http.StatusCreated, response)
	
	slog.InfoContext(ctx, "Task created via HTTP",
		slog.String("task_id", response.ID),
		slog.String("title", response.Title),
	)
}

func (h *TaskHandlers) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	completed, err := validator.ValidateCompletedParam(r.URL.Query().Get("completed"))
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks, totalCount, err := h.taskService.GetTasks(ctx, completed)
	if err != nil {
		h.handleServiceError(w, err, "Failed to get tasks")
		return
	}

	response := dto.TaskModelsToResponse(tasks)

	WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
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

	if err := validator.ValidateTaskID(taskID); err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := h.taskService.GetTask(ctx, taskID)
	if err != nil {
		h.handleServiceError(w, err, "Failed to get task")
		return
	}

	response := dto.TaskModelToResponse(task)

	WriteJSONResponse(w, http.StatusOK, response)

	slog.InfoContext(ctx, "Task retrieved via HTTP",
		slog.String("task_id", response.ID),
	)
}

func (h *TaskHandlers) HandleUpdateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	taskID := vars["id"]

	if err := validator.ValidateTaskID(taskID); err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req dto.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if validationErrs := validator.ValidateUpdateTaskRequest(req); validationErrs.HasErrors() {
		WriteErrorResponse(w, validationErrs.Error(), http.StatusBadRequest)
		return
	}

	updatedTask, err := h.taskService.UpdateTask(ctx, taskID, req.Title, req.Description, req.Completed)
	if err != nil {
		h.handleServiceError(w, err, "Failed to update task")
		return
	}

	response := dto.TaskModelToResponse(updatedTask)

	WriteJSONResponse(w, http.StatusOK, response)

	slog.InfoContext(ctx, "Task updated via HTTP",
		slog.String("task_id", response.ID),
	)
}

func (h *TaskHandlers) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	taskID := vars["id"]

	if err := validator.ValidateTaskID(taskID); err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.taskService.DeleteTask(ctx, taskID)
	if err != nil {
		h.handleServiceError(w, err, "Failed to delete task")
		return
	}

	response := dto.DeleteTaskResponse{Success: true}
	WriteJSONResponse(w, http.StatusOK, response)

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

	WriteErrorResponse(w, message, statusCode)
}
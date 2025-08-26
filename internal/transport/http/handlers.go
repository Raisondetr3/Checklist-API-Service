package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Raisondetr3/checklist-api-service/internal/config"
	"github.com/Raisondetr3/checklist-api-service/internal/service"
	"github.com/Raisondetr3/checklist-api-service/pkg/dto"

	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	config         *config.Config
	taskHandlers   *TaskHandlers
	healthHandlers *HealthHandlers
}

func NewHTTPHandlers(cfg *config.Config, taskService service.TaskService, healthService service.HealthService) *HTTPHandlers {
	return &HTTPHandlers{
		config:         cfg,
		taskHandlers:   NewTaskHandlers(taskService),
		healthHandlers: NewHealthHandlers(healthService),
	}
}

func (h *HTTPHandlers) SetupRoutes(router *mux.Router) {

	router.HandleFunc("/health", h.healthHandlers.HandleHealthCheck).Methods("GET")
	router.HandleFunc("/", h.RootHandler).Methods("GET")

	v1 := router.PathPrefix("/api/v1").Subrouter()
	
	v1.HandleFunc("/tasks", h.taskHandlers.HandleCreateTask).Methods("POST")
	v1.HandleFunc("/tasks", h.taskHandlers.HandleGetTasks).Methods("GET")
	v1.HandleFunc("/tasks/{id}", h.taskHandlers.HandleGetTask).Methods("GET")
	v1.HandleFunc("/tasks/{id}", h.taskHandlers.HandleUpdateTask).Methods("PUT", "PATCH")
	v1.HandleFunc("/tasks/{id}", h.taskHandlers.HandleDeleteTask).Methods("DELETE")
}

func (h *HTTPHandlers) RootHandler(w http.ResponseWriter, r *http.Request) {
	info := map[string]interface{}{
		"service": "checklist-api-service",
		"version": "1.0.0",
		"status":  "running",
		"endpoints": map[string][]string{
			"tasks": {
				"POST /api/v1/tasks - Create task",
				"GET /api/v1/tasks - List tasks (?completed=true/false)",
				"GET /api/v1/tasks/{id} - Get task",
				"PUT /api/v1/tasks/{id} - Update task",
				"PATCH /api/v1/tasks/{id} - Partial update task",
				"DELETE /api/v1/tasks/{id} - Delete task",
			},
			"health": {
				"GET /health - Health check",
			},
		},
	}

	WriteJSONResponse(w, http.StatusOK, info)
}

func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		jsonBytes, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			slog.Error("Failed to marshal JSON response", slog.String("error", err.Error()))
			if encodeErr := json.NewEncoder(w).Encode(map[string]string{
				"error": "Failed to encode response",
			}); encodeErr != nil {
				slog.Error("Failed to encode fallback JSON response", slog.String("error", encodeErr.Error()))
			}
			return
		}
		
		if _, writeErr := w.Write(jsonBytes); writeErr != nil {
			slog.Error("Failed to write JSON response to client", 
				slog.String("error", writeErr.Error()),
				slog.Int("status_code", statusCode),
			)
		}
	}
}

func WriteErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	errDTO := dto.NewErr(message)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	errJSON := errDTO.ToString()
	if _, writeErr := w.Write([]byte(errJSON)); writeErr != nil {
		slog.Error("Failed to write error response to client",
			slog.String("error", writeErr.Error()),
			slog.String("message", message),
			slog.Int("status_code", statusCode),
		)
	}

	slog.WarnContext(context.Background(), "HTTP error response sent",
		slog.Int("status_code", statusCode),
		slog.String("message", message),
	)
}
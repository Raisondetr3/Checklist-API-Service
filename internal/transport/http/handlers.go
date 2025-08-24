package http

import (
	"checklist-api-service/internal/config"

	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	config *config.Config
}

func NewHTTPHandlers(cfg *config.Config) *HTTPHandlers {
	return &HTTPHandlers{
		config: cfg,
	}
}

func (h *HTTPHandlers) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/tasks", h.HandleCreateTask).Methods("POST")
	router.HandleFunc("/tasks", h.HandleGetTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", h.HandleGetTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", h.HandleUpdateTask).Methods("PUT")
	router.HandleFunc("/tasks/{id}", h.HandleDeleteTask).Methods("DELETE")
}

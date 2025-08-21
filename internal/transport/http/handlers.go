package http

import (
	"Checklist/internal/config"

	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	config *config.Config
}

func NewHTTPHandlers() *HTTPHandlers {
	return &HTTPHandlers{}
}

func (h *HTTPHandlers) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/tasks", h.CreateTask).Methods("POST")
	router.HandleFunc("/tasks", h.GetTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", h.GetTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", h.UpdateTask).Methods("PUT")
	router.HandleFunc("/tasks/{id}", h.DeleteTask).Methods("DELETE")
}

package http

import (
	"net/http"

	"github.com/Raisondetr3/checklist-api-service/internal/model"
	"github.com/Raisondetr3/checklist-api-service/internal/service"
	"github.com/Raisondetr3/checklist-api-service/pkg/dto"
)

type HealthHandlers struct {
	healthService service.HealthService
}

func NewHealthHandlers(healthService service.HealthService) *HealthHandlers {
	return &HealthHandlers{
		healthService: healthService,
	}
}

func (h *HealthHandlers) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	health, _ := h.healthService.CheckHealth(ctx)
	
	healthStatus := &dto.HealthStatus{
		Status:    string(health.Status),
		Timestamp: health.Timestamp,
	}

	statusCode := h.getHTTPStatusCode(health.Status)
	
	WriteJSONResponse(w, statusCode, healthStatus)
}

func (h *HealthHandlers) getHTTPStatusCode(status model.HealthStatus) int {
	switch status {
	case model.HealthStatusHealthy:
		return http.StatusOK
	default:
		return http.StatusServiceUnavailable
	}
}
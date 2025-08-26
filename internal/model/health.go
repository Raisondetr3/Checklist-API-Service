package model

import "time"

type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
)

type Health struct {
	Status    HealthStatus
	Timestamp time.Time
}

func NewHealth(status HealthStatus) *Health {
	return &Health{
		Status:    status,
		Timestamp: time.Now(),
	}
}

func ParseDBHealthResponse(dbHealth map[string]interface{}) HealthStatus {
	status, ok := dbHealth["status"].(string)
	if !ok {
		return HealthStatusUnhealthy
	}

	switch status {
	case "healthy":
		return HealthStatusHealthy
	default:
		return HealthStatusUnhealthy
	}
}
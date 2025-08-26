package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Raisondetr3/checklist-api-service/internal/config"
	"github.com/Raisondetr3/checklist-api-service/internal/model"
	"github.com/Raisondetr3/checklist-api-service/pkg/logger"
)

type HealthService interface {
	CheckHealth(ctx context.Context) (*model.Health, error)
}

type healthService struct {
	config     *config.Config
	httpClient *http.Client
}

func NewHealthService(cfg *config.Config) HealthService {
	return &healthService{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.ExternalServices.DBService.Timeout,
		},
	}
}

func (s *healthService) CheckHealth(ctx context.Context) (*model.Health, error) {
	dbURL := fmt.Sprintf("%s/health", s.config.ExternalServices.DBService.URL)

	req, err := http.NewRequestWithContext(ctx, "GET", dbURL, nil)
	if err != nil {
		logger.LogError(ctx, err, "health_check_create_request")
		return model.NewHealth(model.HealthStatusUnhealthy), err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		logger.LogError(ctx, err, "health_check_db_request")
		return model.NewHealth(model.HealthStatusUnhealthy), err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := errors.New("db-service returned non-200 status")
		logger.LogError(ctx, err, "health_check_db_status")
		return model.NewHealth(model.HealthStatusUnhealthy), err
	}

	var dbHealth map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&dbHealth); err != nil {
		logger.LogError(ctx, err, "health_check_parse_response")
		return model.NewHealth(model.HealthStatusUnhealthy), err
	}

	status := model.ParseDBHealthResponse(dbHealth)
	return model.NewHealth(status), nil
}

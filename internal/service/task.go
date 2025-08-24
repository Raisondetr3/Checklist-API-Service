package service

import (
	"checklist-api-service/internal/client"
	"checklist-api-service/internal/config"
	// "checklist-api-service/internal/model"
	// "context"
)

type TaskServiceImpl struct {
	config        *config.Config
	dbClient      client.DBClient
}

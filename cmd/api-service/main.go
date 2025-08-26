package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Raisondetr3/checklist-api-service/internal/client"
	"github.com/Raisondetr3/checklist-api-service/internal/config"
	"github.com/Raisondetr3/checklist-api-service/internal/service"
	httpTransport "github.com/Raisondetr3/checklist-api-service/internal/transport/http"
	"github.com/Raisondetr3/checklist-api-service/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	loggerCfg := logger.Config{
		Level:    cfg.Logging.Level,
		FilePath: cfg.Logging.FilePath,
		FileName: cfg.Logging.FileName,
	}

	if err := logger.SetupLogger(loggerCfg, "api-service"); err != nil {
		panic("Failed to setup logger: " + err.Error())
	}

	slog.Info("Starting API service",
		slog.String("port", cfg.Server.Port),
		slog.String("log_level", cfg.Logging.Level),
		slog.String("db_service_url", cfg.ExternalServices.DBService.URL),
	)

	grpcClient, err := client.NewTaskClient(cfg.ExternalServices.DBService)
	if err != nil {
		slog.Error("Failed to create gRPC client", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		if err := grpcClient.Close(); err != nil {
			slog.Error("Failed to close gRPC client", slog.String("error", err.Error()))
		}
	}()

	taskService := service.NewTaskService(grpcClient)
	healthService := service.NewHealthService(cfg)


	handlers := httpTransport.NewHTTPHandlers(cfg, taskService, healthService)
	server := httpTransport.NewHTTPServer(cfg, handlers)

	go func() {
		slog.Info("Starting HTTP server", slog.String("port", cfg.Server.Port))
		if err := server.StartServer(); err != nil {
			slog.Error("Failed to start server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		slog.Error("Server forced to shutdown", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("Server exited")
}

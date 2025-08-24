package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Level    string `env:"LOG_LEVEL" envDefault:"info"`
	FilePath string `env:"LOG_FILE_PATH" envDefault:"logs"`
	FileName string `env:"LOG_FILE_NAME"`
}

func SetupLogger(cfg Config, serviceName string) error {
	if err := os.MkdirAll(cfg.FilePath, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	if cfg.FileName == "" {
		cfg.FileName = fmt.Sprintf("%s.log", serviceName)
	}

	fullPath := filepath.Join(cfg.FilePath, cfg.FileName)

	logFile, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format(time.RFC3339))
			}
			return a
		},
	}

	handler := slog.NewJSONHandler(logFile, opts)

	logger := slog.New(handler).With(
		slog.String("service", serviceName),
	)

	slog.SetDefault(logger)

	return nil
}

func LogRequest(ctx context.Context, method, path, userAgent, requestID string, duration time.Duration, statusCode int) {
	attrs := []slog.Attr{
		slog.String("type", "request"),
		slog.String("method", method),
		slog.String("path", path),
		slog.String("user_agent", userAgent),
		slog.String("request_id", requestID),
		slog.Duration("duration", duration),
		slog.Int("status_code", statusCode),
	}

	if statusCode >= 500 {
    	slog.LogAttrs(ctx, slog.LevelError, "HTTP Request", attrs...)
	} else if statusCode >= 400 {
   		slog.LogAttrs(ctx, slog.LevelWarn, "HTTP Request", attrs...)
	} else {
		slog.LogAttrs(ctx, slog.LevelInfo, "HTTP Request", attrs...)
	}
}

func LogError(ctx context.Context, err error, operation string, additionalFields ...slog.Attr) {
	attrs := []slog.Attr{
		slog.String("type", "error"),
		slog.String("operation", operation),
		slog.String("error", err.Error()),
	}
	attrs = append(attrs, additionalFields...)

	slog.LogAttrs(ctx, slog.LevelError, "Operation Error", attrs...)
}
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
	Format   string `env:"LOG_FORMAT" envDefault:"json"`
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

	writer := logFile

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

	var handler slog.Handler
	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(writer, opts)
	} else {
		handler = slog.NewTextHandler(writer, opts)
	}

	logger := slog.New(handler).With(
		slog.String("service", serviceName),
		slog.String("version", getVersion()),
	)

	slog.SetDefault(logger)

	return nil
}

func getVersion() string {
	if version := os.Getenv("APP_VERSION"); version != "" {
		return version
	}
	return "dev"
}

func LogRequest(ctx context.Context, method, path, userAgent, requestID string, duration time.Duration, statusCode int) {
	if ctx == nil {
		ctx = context.Background()
	}

	attrs := []slog.Attr{
		slog.String("type", "request"),
		slog.String("method", method),
		slog.String("path", path),
		slog.String("user_agent", userAgent),
		slog.String("request_id", requestID),
		slog.Duration("duration", duration),
		slog.Int("status_code", statusCode),
	}

	if statusCode >= 400 {
		slog.LogAttrs(ctx, slog.LevelWarn, "HTTP Request", attrs...)
	} else {
		slog.LogAttrs(ctx, slog.LevelInfo, "HTTP Request", attrs...)
	}
}

func LogError(ctx context.Context, err error, operation string, additionalFields ...slog.Attr) {
	if ctx == nil {
		ctx = context.Background()
	}

	attrs := []slog.Attr{
		slog.String("type", "error"),
		slog.String("operation", operation),
		slog.String("error", err.Error()),
	}
	attrs = append(attrs, additionalFields...)

	slog.LogAttrs(ctx, slog.LevelError, "Operation Error", attrs...)
}

func LogDatabaseQuery(ctx context.Context, query string, args []interface{}, duration time.Duration, err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	attrs := []slog.Attr{
		slog.String("type", "database"),
		slog.String("query", query),
		slog.Any("args", args),
		slog.Duration("duration", duration),
	}

	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "Database Query Failed", attrs...)
	} else {
		slog.LogAttrs(ctx, slog.LevelDebug, "Database Query", attrs...)
	}
}

func LogKafkaEvent(ctx context.Context, topic, key string, value interface{}, operation string, err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	attrs := []slog.Attr{
		slog.String("type", "kafka"),
		slog.String("topic", topic),
		slog.String("key", key),
		slog.Any("value", value),
		slog.String("operation", operation),
	}

	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "Kafka Operation Failed", attrs...)
	} else {
		slog.LogAttrs(ctx, slog.LevelInfo, "Kafka Event", attrs...)
	}
}

func WithRequestID(requestID string) *slog.Logger {
	return slog.With(slog.String("request_id", requestID))
}

func WithTaskID(taskID string) *slog.Logger {
	return slog.With(slog.String("task_id", taskID))
}

func LogRequestSimple(method, path, userAgent, requestID string, duration time.Duration, statusCode int) {
	LogRequest(context.Background(), method, path, userAgent, requestID, duration, statusCode)
}

func LogErrorSimple(err error, operation string, additionalFields ...slog.Attr) {
	LogError(context.Background(), err, operation, additionalFields...)
}

func LogDatabaseQuerySimple(query string, args []interface{}, duration time.Duration, err error) {
	LogDatabaseQuery(context.Background(), query, args, duration, err)
}

func LogKafkaEventSimple(topic, key string, value interface{}, operation string, err error) {
	LogKafkaEvent(context.Background(), topic, key, value, operation, err)
}

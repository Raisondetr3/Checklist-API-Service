package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server           ServerConfig
	Logging          LoggingConfig
	ExternalServices ExternalServicesConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type LoggingConfig struct {
	Level    string
	FilePath string
	FileName string
	Format   string
}

type ExternalServicesConfig struct {
	DBService DBServiceConfig
	Kafka     KafkaConfig
}

type DBServiceConfig struct {
	HTTPUrl          string
	GRPCAddress      string
	Timeout          time.Duration
	MaxRetries       int
	RetryDelay       time.Duration
	KeepAliveTime    time.Duration
	KeepAliveTimeout time.Duration
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
	Timeout time.Duration
}

func Load() (*Config, error) {
	cfg := &Config{}
	setDefaults(cfg)
	overrideFromEnv(cfg)

	return cfg, nil
}

func setDefaults(cfg *Config) {
	cfg.Server.Port = "8080"
	cfg.Server.ReadTimeout = 15 * time.Second
	cfg.Server.WriteTimeout = 15 * time.Second
	cfg.Server.IdleTimeout = 60 * time.Second

	cfg.Logging.Level = "info"
	cfg.Logging.FilePath = "logs"
	cfg.Logging.FileName = "api-service.log"
	cfg.Logging.Format = "json"

	cfg.ExternalServices.DBService.HTTPUrl = "http://localhost:8081"
	cfg.ExternalServices.DBService.GRPCAddress = "localhost:9090"
	cfg.ExternalServices.DBService.Timeout = 30 * time.Second
	cfg.ExternalServices.DBService.MaxRetries = 3
	cfg.ExternalServices.DBService.RetryDelay = 1 * time.Second
	cfg.ExternalServices.DBService.KeepAliveTime = 30 * time.Second
	cfg.ExternalServices.DBService.KeepAliveTimeout = 5 * time.Second

	cfg.ExternalServices.Kafka.Brokers = []string{"localhost:9092"}
	cfg.ExternalServices.Kafka.Topic = "checklist-events"
	cfg.ExternalServices.Kafka.Timeout = 10 * time.Second
}

func overrideFromEnv(cfg *Config) {
	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.Server.Port = port
	}
	if timeout := parseDurationFromEnv("SERVER_READ_TIMEOUT"); timeout > 0 {
		cfg.Server.ReadTimeout = timeout
	}
	if timeout := parseDurationFromEnv("SERVER_WRITE_TIMEOUT"); timeout > 0 {
		cfg.Server.WriteTimeout = timeout
	}
	if timeout := parseDurationFromEnv("SERVER_IDLE_TIMEOUT"); timeout > 0 {
		cfg.Server.IdleTimeout = timeout
	}

	if level := os.Getenv("LOG_LEVEL"); level != "" {
		cfg.Logging.Level = level
	}
	if path := os.Getenv("LOG_FILE_PATH"); path != "" {
		cfg.Logging.FilePath = path
	}
	if name := os.Getenv("LOG_FILE_NAME"); name != "" {
		cfg.Logging.FileName = name
	}
	if format := os.Getenv("LOG_FORMAT"); format != "" {
		cfg.Logging.Format = format
	}

	if httpUrl := os.Getenv("DB_SERVICE_HTTP_URL"); httpUrl != "" {
		cfg.ExternalServices.DBService.HTTPUrl = httpUrl
	}
	if grpcAddr := os.Getenv("DB_SERVICE_GRPC_ADDRESS"); grpcAddr != "" {
		cfg.ExternalServices.DBService.GRPCAddress = grpcAddr
	}
	if timeout := parseDurationFromEnv("DB_SERVICE_TIMEOUT"); timeout > 0 {
		cfg.ExternalServices.DBService.Timeout = timeout
	}
	if retries := parseIntFromEnv("DB_SERVICE_MAX_RETRIES"); retries > 0 {
		cfg.ExternalServices.DBService.MaxRetries = retries
	}
	if delay := parseDurationFromEnv("DB_SERVICE_RETRY_DELAY"); delay > 0 {
		cfg.ExternalServices.DBService.RetryDelay = delay
	}
	if keepAlive := parseDurationFromEnv("DB_SERVICE_KEEPALIVE_TIME"); keepAlive > 0 {
		cfg.ExternalServices.DBService.KeepAliveTime = keepAlive
	}
	if keepAliveTimeout := parseDurationFromEnv("DB_SERVICE_KEEPALIVE_TIMEOUT"); keepAliveTimeout > 0 {
		cfg.ExternalServices.DBService.KeepAliveTimeout = keepAliveTimeout
	}

	if kafkaBrokers := os.Getenv("KAFKA_BROKERS"); kafkaBrokers != "" {
		cfg.ExternalServices.Kafka.Brokers = []string{kafkaBrokers}
	}
	if kafkaTopic := os.Getenv("KAFKA_TOPIC"); kafkaTopic != "" {
		cfg.ExternalServices.Kafka.Topic = kafkaTopic
	}
	if timeout := parseDurationFromEnv("KAFKA_TIMEOUT"); timeout > 0 {
		cfg.ExternalServices.Kafka.Timeout = timeout
	}
}

func parseDurationFromEnv(key string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return 0
}

func parseIntFromEnv(key string) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return 0
}
package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server           ServerConfig           `yaml:"server"`
	Logging          LoggingConfig          `yaml:"logging"`
	ExternalServices ExternalServicesConfig `yaml:"external_services"`
}

type ServerConfig struct {
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type LoggingConfig struct {
	Level    string `yaml:"level"`
	FilePath string `yaml:"file_path"`
	FileName string `yaml:"file_name"`
	Format   string `yaml:"format"`
}

type ExternalServicesConfig struct {
	DBService DBServiceConfig `yaml:"db_service"`
	Kafka     KafkaConfig     `yaml:"kafka"`
}

type DBServiceConfig struct {
	URL     string        `yaml:"url"`
	Timeout time.Duration `yaml:"timeout"`
}

type KafkaConfig struct {
	Brokers []string      `yaml:"brokers"`
	Topic   string        `yaml:"topic"`
	Timeout time.Duration `yaml:"timeout"`
}

func Load() (*Config, error) {
	cfg, err := loadConfigFile("configs/config.yml")
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	overrideFromEnv(cfg)

	return cfg, nil
}

func loadConfigFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func overrideFromEnv(cfg *Config) {
	if port := os.Getenv("API_PORT"); port != "" {
		cfg.Server.Port = port
	}

	if level := os.Getenv("LOG_LEVEL"); level != "" {
		cfg.Logging.Level = level
	}
	if path := os.Getenv("LOG_FILE_PATH"); path != "" {
		cfg.Logging.FilePath = path
	}
	if format := os.Getenv("LOG_FORMAT"); format != "" {
		cfg.Logging.Format = format
	}

	if dbURL := os.Getenv("DB_SERVICE_URL"); dbURL != "" {
		cfg.ExternalServices.DBService.URL = dbURL
	}
	if kafkaBrokers := os.Getenv("KAFKA_BROKERS"); kafkaBrokers != "" {
		cfg.ExternalServices.Kafka.Brokers = []string{kafkaBrokers}
	}
	if kafkaTopic := os.Getenv("KAFKA_TOPIC"); kafkaTopic != "" {
		cfg.ExternalServices.Kafka.Topic = kafkaTopic
	}
}

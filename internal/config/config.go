package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	JWTSecret        string
	LogLevel         string
	ClickhouseConfig *ClickhouseConfig
	PostgresDSN      string
	ExternalAPIs     *ExternalAPIs
}

type ClickhouseConfig struct {
	Host     string
	Database string
	Username string
	Password string
	Debug    bool
}

type ExternalAPIs struct {
	SourcesAPIURL string
	SourcesAPIKey string
}

func NewConfig() (*Config, error) {
	tryLoadEnv(
		"/dc-analytics-service-backend/.env",
		"./.env",
		"../.env",
	)
	if err := validateRequiredVars(); err != nil {
		return nil, err
	}
	debug, _ := strconv.ParseBool(getEnv("CLICKHOUSE_DEBUG", "false"))

	return &Config{
		Port:      getEnv("PORT", "7002"),
		JWTSecret: getEnv("JWT_SECRET", ""),
		LogLevel:  getEnv("LOG_LEVEL", "INFO"),
		ClickhouseConfig: &ClickhouseConfig{
			Host:     getEnv("CLICKHOUSE_HOST", "localhost:9000"),
			Database: getEnv("CLICKHOUSE_DB", "default"),
			Username: getEnv("CLICKHOUSE_USERNAME", "default"),
			Password: getEnv("CLICKHOUSE_PASSWORD", ""),
			Debug:    debug,
		},
		PostgresDSN: getEnv("POSTGRES_DSN", ""),
		ExternalAPIs: &ExternalAPIs{
			SourcesAPIURL: getEnv("SOURCES_API_URL", ""),
			SourcesAPIKey: getEnv("SOURCES_API_KEY", ""),
		},
	}, nil
}

func tryLoadEnv(paths ...string) {
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			if err := godotenv.Load(path); err != nil {
				log.Printf("Warning: Could not load .env file at %s: %v", path, err)
			} else {
				log.Printf("Loaded .env file from: %s", path)
				break
			}
		}
	}
}

func validateRequiredVars() error {
	required := map[string]string{
		"JWT_SECRET":      "JWT secret key",
		"SOURCES_API_KEY": "Sources API key",
		"POSTGRES_DSN":    "PostgreSQL DSN",
	}

	for envVar, desc := range required {
		if os.Getenv(envVar) == "" {
			return fmt.Errorf("missing required environment variable: %s (%s)", envVar, desc)
		}
	}
	return nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

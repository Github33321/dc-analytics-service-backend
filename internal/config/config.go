package config

import (
	"os"
)

type Config struct {
	Port          string
	LogLevel      string
	ClickHouseDSN string
	JWTSecret     string
	PostgresDSN   string
}

func NewConfig() (*Config, error) {
	return &Config{
		Port:          getEnv("PORT", "8081"),
		LogLevel:      getEnv("LOG_LEVEL", "INFO"),
		ClickHouseDSN: getEnv("CLICKHOUSE_DSN", "tcp://localhost:9000?username=default&password=default&debug=true"),
		JWTSecret:     getEnv("JWT_SECRET", "default_super_secret"),
		PostgresDSN:   getEnv("POSTGRES_DSN", ""),
	}, nil
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

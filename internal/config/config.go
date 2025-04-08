package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type ClickhouseConfig struct {
	Host     string
	Database string
	Username string
	Password string
	Debug    bool
}

type Config struct {
	Port             string
	JWTSecret        string
	LogLevel         string
	ClickhouseConfig *ClickhouseConfig
	PostgresDSN      string
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load("/opt/dc-analytics-service/.env"); err != nil {
		fmt.Println("Не удалось загрузить файл .env, используются переменные окружения из системы")
	}
	debug := getEnv("CLICKHOUSE_DEBUG", "false")
	return &Config{
		Port:      getEnv("PORT", "8081"),
		JWTSecret: getEnv("JWT_SECRET", "default_super_secret"),
		LogLevel:  getEnv("LOG_LEVEL", "INFO"),
		ClickhouseConfig: &ClickhouseConfig{
			Host:     getEnv("CLICKHOUSE_HOST", "localhost:9000"),
			Database: getEnv("CLICKHOUSE_DB", "default"),
			Username: getEnv("CLICKHOUSE_USERNAME", "default"),
			Password: getEnv("CLICKHOUSE_PASSWORD", "default"),
			Debug:    debug == "true",
		},
		PostgresDSN: getEnv("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/analytics?sslmode=disable"),
	}, nil
}

func getEnv(key, defaultVal string) string {

	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

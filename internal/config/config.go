package config

import "os"

type ClickhouseConfig struct {
	ConnectionURL string
	Database      string
}

type Config struct {
	Port             string
	JWTSecret        string
	LogLevel         string
	ClickhouseConfig *ClickhouseConfig
	PostgresDSN      string
}

func NewConfig() (*Config, error) {
	return &Config{
		Port:      getEnv("PORT", "8081"),
		JWTSecret: getEnv("JWT_SECRET", "default_super_secret"),
		LogLevel:  getEnv("LOG_LEVEL", "INFO"),
		ClickhouseConfig: &ClickhouseConfig{
			ConnectionURL: getEnv("CLICKHOUSE_DSN", "tcp://localhost:9000?username=default&password=default&debug=true"),
			Database:      getEnv("CLICKHOUSE_DB", "default"),
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

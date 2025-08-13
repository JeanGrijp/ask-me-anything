package config

import (
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   LoggerConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type LoggerConfig struct {
	Level  string
	Format string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "ama_user"),
			Password: getEnv("DB_PASSWORD", "ama_password"),
			DBName:   getEnv("DB_NAME", "ask_me_anything"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Logger: LoggerConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		slog.Warn("Invalid integer value for environment variable", "key", key, "value", value)
	}
	return defaultValue
}

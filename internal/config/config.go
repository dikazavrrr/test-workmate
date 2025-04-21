package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port         int
	PostgresMain PostgresConfig
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func MustInit(isProd string) *Config {
	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		log.Fatalf("Invalid PORT: %v", err)
	}

	pgPort, err := strconv.Atoi(getEnv("POSTGRES_PORT", "5432"))
	if err != nil {
		log.Fatalf("Invalid POSTGRES_PORT: %v", err)
	}

	return &Config{
		Port: port,
		PostgresMain: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     pgPort,
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "password"),
			DBName:   getEnv("POSTGRES_DB", "testdb"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},
	}
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

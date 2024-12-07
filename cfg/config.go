package config

import (
	"github.com/jorgeAM/base-api/internal/platform/env"
)

type Config struct {
	Port                       string
	PostgresHost               string
	PostgresPort               int
	PostgresDatabase           string
	PostgresUser               string
	PostgresPassword           string
	PostgresMaxIdleConnections int
	PostgresMaxOpenConnections int
}

func LoadConfig() (*Config, error) {
	return &Config{
		Port:                       env.GetEnv("PORT", "8080"),
		PostgresHost:               env.GetEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:               env.GetEnv("POSTGRES_PORT", 5432),
		PostgresDatabase:           env.GetEnv("POSTGRES_DB", "coderhouse"),
		PostgresUser:               env.GetEnv("POSTGRES_USER", "admin"),
		PostgresPassword:           env.GetEnv("POSTGRES_PASSWORD", "passwd123"),
		PostgresMaxIdleConnections: env.GetEnv("POSTGRES_MAX_IDLE_CONNECTIONS", 10),
		PostgresMaxOpenConnections: env.GetEnv("POSTGRES_MAX_OPEN_CONNECTIONS", 30),
	}, nil
}

package user

import "github.com/jorgeAM/go-template/internal/shared/env"

type Config struct {
	PostgresHost               string
	PostgresPort               int
	PostgresDatabase           string
	PostgresUser               string
	PostgresPassword           string
	PostgresMaxIdleConnections int
	PostgresMaxOpenConnections int
}

func LoadConfig() *Config {
	return &Config{
		PostgresHost:               env.GetEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:               env.GetEnv("POSTGRES_PORT", 5432),
		PostgresDatabase:           env.GetEnv("POSTGRES_DB", "db"),
		PostgresUser:               env.GetEnv("POSTGRES_USER", "admin"),
		PostgresPassword:           env.GetEnv("POSTGRES_PASSWORD", "passwd123"),
		PostgresMaxIdleConnections: env.GetEnv("POSTGRES_MAX_IDLE_CONNECTIONS", 10),
		PostgresMaxOpenConnections: env.GetEnv("POSTGRES_MAX_OPEN_CONNECTIONS", 30),
	}
}

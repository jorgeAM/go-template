package config

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	userDomain "github.com/jorgeAM/go-template/internal/user/domain"
	userPersistence "github.com/jorgeAM/go-template/internal/user/infrastructure/persistence"
	_ "github.com/lib/pq"
)

type Dependencies struct {
	UserRepository userDomain.UserRepository
}

func BuildDependencies(cfg *Config) (*Dependencies, error) {
	postgresClient, err := sqlx.Connect(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.PostgresHost,
			cfg.PostgresPort,
			cfg.PostgresUser,
			cfg.PostgresPassword,
			cfg.PostgresDatabase,
		),
	)
	if err != nil {
		return nil, err
	}

	postgresUserRepository := userPersistence.NewPostgresUserRepository(postgresClient)

	return &Dependencies{
		UserRepository: postgresUserRepository,
	}, nil
}

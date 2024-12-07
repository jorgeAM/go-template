package config

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	userDomain "github.com/jorgeAM/base-api/internal/user/domain"
	"github.com/jorgeAM/base-api/internal/user/infrastructure"
	_ "github.com/lib/pq"
)

type Dependencies struct {
	PostgresClient *sqlx.DB
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

	return &Dependencies{
		PostgresClient: postgresClient,
		UserRepository: infrastructure.NewInMemUserRepository(),
	}, nil
}

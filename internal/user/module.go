package user

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jorgeAM/go-template/internal/platform/log"
	"github.com/jorgeAM/go-template/internal/shared/module"
	userdb "github.com/jorgeAM/go-template/internal/user/adapters/db"
	"github.com/jorgeAM/go-template/internal/user/domain"
	userhttp "github.com/jorgeAM/go-template/internal/user/infrastructure/http"
)

var _ module.Module = (*Module)(nil)

type Module struct {
	Pool           *pgxpool.Pool
	UserRepository domain.UserRepository
}

func NewModule() *Module {
	return &Module{}
}

func (m *Module) Name() module.Name {
	return "user"
}

func (m *Module) Init(ctx context.Context) (err error) {
	cfg := LoadConfig()

	poolCfg, err := pgxpool.ParseConfig(
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
		log.Error(ctx, "user module failed to parse postgres config", log.WithError(err))
		return err
	}

	poolCfg.MaxConns = int32(cfg.PostgresMaxOpenConnections)

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		log.Error(ctx, "user module failed to create postgres pool", log.WithError(err))
		return err
	}

	defer func() {
		if err != nil {
			pool.Close()
		}
	}()

	// pgxpool connects lazily; ping so a bad DSN fails at startup, not on the first request.
	if err = pool.Ping(ctx); err != nil {
		log.Error(ctx, "user module failed to connect to postgres", log.WithError(err))
		return err
	}

	m.Pool = pool
	m.UserRepository = userdb.NewPostgresUserRepository(pool)

	log.Info(ctx, "user module initialized")

	return nil
}

func (m *Module) RegisterHttp(_ context.Context, r chi.Router) error {
	userhttp.Register(r, m.UserRepository)

	return nil
}

func (m *Module) MigrationFS() fs.FS {
	return userdb.Migrations()
}

package identity

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	identitydb "github.com/jorgeAM/go-template/internal/identity/adapters/db"
	identityhttp "github.com/jorgeAM/go-template/internal/identity/api/http"
	"github.com/jorgeAM/go-template/internal/identity/domain"
	"github.com/jorgeAM/go-template/internal/platform/log"
	"github.com/jorgeAM/go-template/internal/shared/module"
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
	return "identity"
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
		log.Error(ctx, "identity module failed to parse postgres config", log.WithError(err))
		return err
	}

	poolCfg.MaxConns = int32(cfg.PostgresMaxOpenConnections)

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		log.Error(ctx, "identity module failed to create postgres pool", log.WithError(err))
		return err
	}

	defer func() {
		if err != nil {
			pool.Close()
		}
	}()

	if err = pool.Ping(ctx); err != nil {
		log.Error(ctx, "identity module failed to connect to postgres", log.WithError(err))
		return err
	}

	m.Pool = pool
	m.UserRepository = identitydb.NewPostgresUserRepository(pool)

	log.Info(ctx, "identity module initialized")

	return nil
}

func (m *Module) RegisterHttp(ctx context.Context, r chi.Router) error {
	return identityhttp.Register(ctx, r, m.UserRepository)
}

func (m *Module) MigrationFS() fs.FS {
	return identitydb.Migrations()
}

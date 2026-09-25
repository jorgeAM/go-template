package user

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/jorgeAM/go-template/internal/platform/log"
	"github.com/jorgeAM/go-template/internal/shared/module"
	"github.com/jorgeAM/go-template/internal/user/domain"
	userhttp "github.com/jorgeAM/go-template/internal/user/infrastructure/http"
	"github.com/jorgeAM/go-template/internal/user/infrastructure/persistence"
)

var _ module.Module = (*Module)(nil)

type Module struct {
	DB             *sqlx.DB
	UserRepository domain.UserRepository
}

func NewModule() *Module {
	return &Module{}
}

func (m *Module) Name() module.Name {
	return "user"
}

func (m *Module) Init(ctx context.Context) error {
	cfg := LoadConfig()

	db, err := sqlx.ConnectContext(
		ctx,
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
		log.Error(ctx, "user module failed to connect to postgres", log.WithError(err))
		return err
	}

	db.SetMaxIdleConns(cfg.PostgresMaxIdleConnections)
	db.SetMaxOpenConns(cfg.PostgresMaxOpenConnections)

	m.DB = db
	m.UserRepository = persistence.NewPostgresUserRepository(db)

	log.Info(ctx, "user module initialized")

	return nil
}

func (m *Module) RegisterHttp(_ context.Context, r chi.Router) error {
	userhttp.Register(r, m.UserRepository)

	return nil
}

func (m *Module) MigrationFS() fs.FS {
	return persistence.Migrations()
}

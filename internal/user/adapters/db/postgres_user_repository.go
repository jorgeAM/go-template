package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	platformdb "github.com/jorgeAM/go-template/internal/platform/db"
	"github.com/jorgeAM/go-template/internal/user/adapters/db/sqlc"
	"github.com/jorgeAM/go-template/internal/user/domain"
)

var _ domain.UserRepository = (*PostgresUserRepository)(nil)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// queries joins the transaction a platformdb.Transactor put in ctx, if any.
func (r *PostgresUserRepository) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(platformdb.TxKey("tx")).(pgx.Tx); ok {
		return sqlc.New(tx)
	}

	return sqlc.New(r.db)
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *domain.User) error {
	return r.queries(ctx).SaveUser(ctx, sqlc.SaveUserParams{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	})
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	row, err := r.queries(ctx).FindUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	return &domain.User{
		ID:        row.ID,
		Name:      row.Name,
		Email:     row.Email,
		Password:  row.Password,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		DeletedAt: row.DeletedAt,
	}, nil
}

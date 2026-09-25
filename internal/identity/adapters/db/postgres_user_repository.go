package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jorgeAM/go-template/internal/identity/adapters/db/sqlc"
	"github.com/jorgeAM/go-template/internal/identity/domain"
	platformdb "github.com/jorgeAM/go-template/internal/platform/db"
	"github.com/jorgeAM/go-template/internal/shared/errors"
	"github.com/jorgeAM/go-template/internal/shared/valueobject"
)

var _ domain.UserRepository = (*PostgresUserRepository)(nil)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(platformdb.TxKey("tx")).(pgx.Tx); ok {
		return sqlc.New(tx)
	}

	return sqlc.New(r.db)
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *domain.User) error {
	timestamps := user.Timestamps()

	err := r.queries(ctx).SaveUser(ctx, sqlc.SaveUserParams{
		ID:        user.ID().String(),
		Name:      user.Name(),
		Email:     user.Email().String(),
		Password:  user.HashedPassword(),
		CreatedAt: timestamps.CreatedAt,
		UpdatedAt: timestamps.UpdatedAt,
		DeletedAt: timestamps.DeletedAt,
	})
	if err != nil {
		return errors.Wrap(
			domain.ErrUserInternal,
			err,
			"an error occurred while saving the user",
			errors.WithMetadata("id", user.ID().String()),
		)
	}

	return nil
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id valueobject.ID) (*domain.User, error) {
	row, err := r.queries(ctx).FindUserByID(ctx, id.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(domain.ErrUserNotFound, err, "user not found", errors.WithMetadata("id", id.String()))
		}

		return nil, errors.Wrap(
			domain.ErrUserInternal,
			err,
			"an error occurred while retrieving the user",
			errors.WithMetadata("id", id.String()),
		)
	}

	return domain.UnmarshallUser(
		valueobject.ID(row.ID),
		row.Name,
		valueobject.Email(row.Email),
		row.Password,
		valueobject.Timestamps{
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
			DeletedAt: row.DeletedAt,
		},
	), nil
}

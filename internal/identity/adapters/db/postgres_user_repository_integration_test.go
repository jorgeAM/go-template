//go:build integration

package db

import (
	"context"
	"io/fs"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/jorgeAM/go-template/internal/identity/domain"
	platformdb "github.com/jorgeAM/go-template/internal/platform/db"
	"github.com/jorgeAM/go-template/internal/shared/errors"
	"github.com/jorgeAM/go-template/internal/shared/valueobject"
)

func TestPostgresUserRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()
	pool := startPostgres(ctx, t)
	repo := NewPostgresUserRepository(pool)

	t.Run("Save", func(t *testing.T) {
		t.Run("persists a new user", func(t *testing.T) {
			user := newUser(t)

			require.NoError(t, repo.Save(ctx, user))

			got, err := repo.FindByID(ctx, user.ID())
			require.NoError(t, err)
			assertSameUser(t, user, got)
		})

		t.Run("saving the same user twice upserts instead of failing", func(t *testing.T) {
			user := newUser(t)

			require.NoError(t, repo.Save(ctx, user))
			require.NoError(t, repo.Save(ctx, user))

			got, err := repo.FindByID(ctx, user.ID())
			require.NoError(t, err)
			assertSameUser(t, user, got)
		})

		t.Run("uses the transaction from context", func(t *testing.T) {
			user := newUser(t)
			transactor := platformdb.NewPgxTransactorManager(pool)

			err := transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
				require.NoError(t, repo.Save(txCtx, user))
				return assert.AnError
			})
			require.ErrorIs(t, err, assert.AnError)

			_, err = repo.FindByID(ctx, user.ID())
			assert.True(t, errors.Is(err, domain.ErrUserNotFound), "rolled back user must not be persisted, got %v", err)
		})
	})

	t.Run("FindByID", func(t *testing.T) {
		t.Run("returns not found for an unknown id", func(t *testing.T) {
			id, err := valueobject.NewUUIDv7()
			require.NoError(t, err)

			got, err := repo.FindByID(ctx, id)

			assert.Nil(t, got)
			assert.True(t, errors.Is(err, domain.ErrUserNotFound), "got %v", err)
		})
	})
}

func startPostgres(ctx context.Context, t *testing.T) *pgxpool.Pool {
	t.Helper()

	container, err := postgres.Run(ctx, "postgres:17-alpine",
		postgres.WithDatabase("identity_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.BasicWaitStrategies(),
	)
	testcontainers.CleanupContainer(t, container)
	require.NoError(t, err)

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	applyMigrations(ctx, t, pool)

	return pool
}

func applyMigrations(ctx context.Context, t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	migrations := Migrations()

	// fs.Glob returns names in lexical order, which matches the NNNNNN_ migration sequence.
	files, err := fs.Glob(migrations, "*.up.sql")
	require.NoError(t, err)
	require.NotEmpty(t, files)

	for _, file := range files {
		sql, err := fs.ReadFile(migrations, file)
		require.NoError(t, err)

		_, err = pool.Exec(ctx, string(sql))
		require.NoError(t, err, "applying %s", file)
	}
}

func newUser(t *testing.T) *domain.User {
	t.Helper()

	user, err := domain.NewUser("Jorge", "jorge@example.com", "s3cure-pass")
	require.NoError(t, err)

	return user
}

func assertSameUser(t *testing.T, want, got *domain.User) {
	t.Helper()

	assert.Equal(t, want.ID(), got.ID())
	assert.Equal(t, want.Name(), got.Name())
	assert.Equal(t, want.Email(), got.Email())
	assert.Equal(t, want.HashedPassword(), got.HashedPassword())

	// Postgres timestamptz stores microseconds; Go keeps nanoseconds.
	assert.True(t, want.Timestamps().CreatedAt.Truncate(time.Microsecond).Equal(got.Timestamps().CreatedAt))
	assert.True(t, want.Timestamps().UpdatedAt.Truncate(time.Microsecond).Equal(got.Timestamps().UpdatedAt))
	assert.Nil(t, got.Timestamps().DeletedAt)
}

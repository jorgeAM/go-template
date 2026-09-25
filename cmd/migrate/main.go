package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/golang-migrate/migrate/v4"
	pgxv5 "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"

	"github.com/jorgeAM/go-template/internal/bootstrap"
	"github.com/jorgeAM/go-template/internal/platform/log"
	"github.com/jorgeAM/go-template/internal/shared/env"
)

// Command migrate applies every module's pending up migrations, walking the
// same bootstrap.Modules() list cmd/app boots from. Each module tracks its
// version in its own schema_migrations_<name> table so modules stay independent.
func main() {
	ctx := context.Background()

	if err := log.InitDefaultLogger(); err != nil {
		panic(err)
	}

	db, err := sql.Open("pgx", buildDSN())
	if err != nil {
		log.Error(ctx, "migrate: failed to open database", log.WithError(err))
		os.Exit(1)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		log.Error(ctx, "migrate: database unreachable", log.WithError(err))
		os.Exit(1)
	}

	for _, m := range bootstrap.Modules() {
		migrations := m.MigrationFS()
		if migrations == nil {
			continue
		}

		name := string(m.Name())
		if err := up(ctx, db, name, migrations); err != nil {
			log.Error(ctx, "migrate: module failed", log.WithString("module", name), log.WithError(err))
			os.Exit(1)
		}
	}

	log.Info(ctx, "migrate: all modules up to date")
}

func up(ctx context.Context, db *sql.DB, name string, migrations fs.FS) error {
	source, err := iofs.New(migrations, ".")
	if err != nil {
		return fmt.Errorf("load migration files: %w", err)
	}

	driver, err := pgxv5.WithInstance(db, &pgxv5.Config{MigrationsTable: migrationsTable(name)})
	if err != nil {
		return fmt.Errorf("build migration driver: %w", err)
	}

	runner, err := migrate.NewWithInstance("iofs", source, "pgx5", driver)
	if err != nil {
		return fmt.Errorf("build migration runner: %w", err)
	}

	switch err := runner.Up(); {
	case errors.Is(err, migrate.ErrNoChange):
		log.Info(ctx, "migrate: no change", log.WithString("module", name))
	case err != nil:
		return err
	default:
		log.Info(ctx, "migrate: applied", log.WithString("module", name))
	}

	return nil
}

func migrationsTable(module string) string {
	return "schema_migrations_" + module
}

func buildDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		env.GetEnv("POSTGRES_HOST", "localhost"),
		env.GetEnv("POSTGRES_PORT", 5432),
		env.GetEnv("POSTGRES_USER", "admin"),
		env.GetEnv("POSTGRES_PASSWORD", "passwd123"),
		env.GetEnv("POSTGRES_DB", "db"),
	)
}

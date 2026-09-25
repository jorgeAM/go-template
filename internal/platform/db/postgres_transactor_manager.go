package db

import (
	"context"

	"github.com/jmoiron/sqlx"
)

var _ Transactor = (*PostgresTransactorManager)(nil)

type PostgresTransactorManager struct {
	db *sqlx.DB
}

func NewPostgresTransactorManager(db *sqlx.DB) *PostgresTransactorManager {
	return &PostgresTransactorManager{
		db: db,
	}
}

func (p *PostgresTransactorManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(context.WithValue(ctx, TxKey("tx"), tx)); err != nil {
		tx.Rollback()

		return err
	}

	return tx.Commit()
}

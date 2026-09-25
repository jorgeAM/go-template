package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ Transactor = (*PgxTransactorManager)(nil)

type PgxTransactorManager struct {
	pool *pgxpool.Pool
}

func NewPgxTransactorManager(pool *pgxpool.Pool) *PgxTransactorManager {
	return &PgxTransactorManager{pool: pool}
}

func (p *PgxTransactorManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}

	if err := fn(context.WithValue(ctx, TxKey("tx"), tx)); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return errors.Join(err, rbErr)
		}

		return err
	}

	return tx.Commit(ctx)
}

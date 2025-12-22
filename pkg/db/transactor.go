package db

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type TxKey string

type DBOrTx interface {
	sqlx.ExtContext
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

//go:generate mockgen -source=./transactor.go -destination=./mocks/transactor.go -package=mock -mock_names=Transactor=MockTransactor
type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

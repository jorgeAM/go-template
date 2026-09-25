package db

import (
	"context"
)

type TxKey string

//go:generate go tool mockgen -source=./transactor.go -destination=./mocks/transactor.go -package=mock -mock_names=Transactor=MockTransactor
type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

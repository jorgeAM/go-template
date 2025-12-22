package storage

import (
	"context"

	"github.com/jorgeAM/go-template/pkg/errors"
)

var (
	ErrStorageInternal = errors.Define("storage.internal_error")
)

//go:generate mockgen -source=./storage.go -destination=./mocks/storage.go -package=mock -mock_names=Signer=MockSigner
type Signer interface {
	GeneratePresignedURL(ctx context.Context, filename string, contentType ContentType) (string, error)
}

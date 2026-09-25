package module

import (
	"context"
	"io/fs"

	"github.com/go-chi/chi/v5"
)

type Name string

type Module interface {
	Name() Name
	Init(ctx context.Context) error
	RegisterHttp(ctx context.Context, r chi.Router) error
	MigrationFS() fs.FS
}

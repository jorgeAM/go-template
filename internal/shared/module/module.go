package module

import (
	"context"
	"io/fs"

	"github.com/go-chi/chi/v5"
)

type Name string

type Module interface {
	Name() Name
	// Init builds every dependency the module owns (DB pool, repositories,
	// platform adapters). It runs once, before RegisterHttp.
	Init(ctx context.Context) error
	RegisterHttp(ctx context.Context, r chi.Router) error
	// MigrationFS returns the module's SQL migrations rooted at the directory
	// holding them, or nil when the module owns no schema.
	MigrationFS() fs.FS
}

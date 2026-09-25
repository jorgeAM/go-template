package persistence

import (
	"embed"
	"io/fs"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// Migrations returns the module's migration files rooted at the migrations
// directory, so the runner sees the .sql files directly.
func Migrations() fs.FS {
	sub, err := fs.Sub(migrationFiles, "migrations")
	if err != nil {
		panic(err) // embed path is a compile-time constant; cannot fail
	}

	return sub
}

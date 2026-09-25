package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildDSN(t *testing.T) {
	t.Setenv("POSTGRES_HOST", "db.example.com")
	t.Setenv("POSTGRES_PORT", "6543")
	t.Setenv("POSTGRES_USER", "app")
	t.Setenv("POSTGRES_PASSWORD", "pw")
	t.Setenv("POSTGRES_DB", "app_db")

	got := buildDSN()

	for _, want := range []string{
		"host=db.example.com",
		"port=6543",
		"user=app",
		"password=pw",
		"dbname=app_db",
		"sslmode=disable",
	} {
		assert.True(t, strings.Contains(got, want), "buildDSN() = %q, missing %q", got, want)
	}
}

func TestMigrationsTable(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "schema_migrations_user", migrationsTable("user"))
}

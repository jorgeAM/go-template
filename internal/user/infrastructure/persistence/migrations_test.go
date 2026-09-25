package persistence

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrations(t *testing.T) {
	t.Parallel()

	files, err := fs.Glob(Migrations(), "*.up.sql")

	assert.NoError(t, err)
	assert.NotEmpty(t, files)
}

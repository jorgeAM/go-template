package valueobject

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/jorgeAM/go-template/internal/shared/errors"
)

func TestNewUUIDv7(t *testing.T) {
	t.Parallel()

	first, err := NewUUIDv7()
	assert.NoError(t, err)

	second, err := NewUUIDv7()
	assert.NoError(t, err)

	assert.Equal(t, uuid.Version(7), uuid.UUID(first).Version())
	assert.False(t, first.IsZero())
	assert.False(t, first.Equals(second))
	assert.Less(t, first.String(), second.String(), "v7 ids must sort in generation order")
}

func TestParseUUID(t *testing.T) {
	t.Parallel()

	v7, err := NewUUIDv7()
	assert.NoError(t, err)

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "accepts a v7 uuid", input: v7.String()},
		{name: "accepts a v4 uuid", input: uuid.NewString()},
		{name: "rejects garbage", input: "not-a-uuid", wantErr: true},
		{name: "rejects empty", input: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id, err := ParseUUID(tt.input)

			if tt.wantErr {
				assert.True(t, errors.Is(err, ErrInvalidUUID), "got %v, want %v", err, ErrInvalidUUID)
				assert.True(t, id.IsZero())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.input, id.String())
		})
	}
}

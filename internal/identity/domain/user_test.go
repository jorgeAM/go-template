package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jorgeAM/go-template/internal/shared/crypto"
	"github.com/jorgeAM/go-template/internal/shared/errors"
	"github.com/jorgeAM/go-template/internal/shared/valueobject"
)

func TestNewUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		userName string
		email    string
		password string
		wantErr  *errors.ErrorCode
	}{
		{name: "creates a user", userName: "Jorge", email: "Jorge@Example.com ", password: "s3cure-pass"},
		{name: "trims the name", userName: "  Jorge  ", email: "jorge@example.com", password: "s3cure-pass"},
		{name: "rejects an empty name", userName: "   ", email: "jorge@example.com", password: "s3cure-pass", wantErr: ErrInvalidUser},
		{name: "rejects an invalid email", userName: "Jorge", email: "foo@baz@gmail.com", password: "s3cure-pass", wantErr: ErrInvalidUser},
		{name: "rejects a short password", userName: "Jorge", email: "jorge@example.com", password: "1234567", wantErr: ErrInvalidUser},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			user, err := NewUser(tt.userName, tt.email, tt.password)

			if tt.wantErr != nil {
				assert.True(t, errors.Is(err, tt.wantErr), "got %v, want %v", err, tt.wantErr)
				assert.Nil(t, user)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, strings.TrimSpace(tt.userName), user.Name())
			assert.Equal(t, valueobject.Email("jorge@example.com"), user.Email())
			assert.NotEmpty(t, user.ID())
			assert.NotEqual(t, tt.password, user.HashedPassword())
			assert.True(t, crypto.ComparePassword(user.HashedPassword(), tt.password))
			assert.False(t, user.Timestamps().CreatedAt.IsZero())
		})
	}
}

package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/jorgeAM/go-template/internal/identity/domain"
	identitymock "github.com/jorgeAM/go-template/internal/identity/mocks"
	"github.com/jorgeAM/go-template/internal/shared/errors"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cmd     *CreateUserCommand
		mock    func(repo *identitymock.MockUserRepository)
		wantErr *errors.ErrorCode
	}{
		{
			name: "creates the user",
			cmd:  &CreateUserCommand{Name: "Jorge", Email: "jorge@example.com", Password: "s3cure-pass"},
			mock: func(repo *identitymock.MockUserRepository) {
				repo.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(&domain.User{})).Return(nil)
			},
		},
		{
			name: "rejects invalid user data without saving",
			cmd:  &CreateUserCommand{Name: "Jorge", Email: "foo@baz@gmail.com", Password: "s3cure-pass"},
			mock: func(repo *identitymock.MockUserRepository) {
				repo.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: domain.ErrInvalidUser,
		},
		{
			name: "maps a repository failure to internal",
			cmd:  &CreateUserCommand{Name: "Jorge", Email: "jorge@example.com", Password: "s3cure-pass"},
			mock: func(repo *identitymock.MockUserRepository) {
				repo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			wantErr: domain.ErrUserInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repo := identitymock.NewMockUserRepository(ctrl)
			tt.mock(repo)

			res, err := NewCreateUser(repo).Handle(context.Background(), tt.cmd)

			if tt.wantErr != nil {
				assert.True(t, errors.Is(err, tt.wantErr), "got %v, want %v", err, tt.wantErr)
				assert.Nil(t, res)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, res.ID)
			assert.Equal(t, "Jorge", res.Name)
			assert.Equal(t, "jorge@example.com", res.Email)
		})
	}
}

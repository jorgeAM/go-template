package query

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/jorgeAM/go-template/internal/identity/domain"
	identitymock "github.com/jorgeAM/go-template/internal/identity/mocks"
	"github.com/jorgeAM/go-template/internal/shared/errors"
)

func TestGetUser(t *testing.T) {
	t.Parallel()

	user, err := domain.NewUser("Jorge", "jorge@example.com", "s3cure-pass")
	assert.NoError(t, err)

	tests := []struct {
		name    string
		q       *GetUserQuery
		mock    func(repo *identitymock.MockUserRepository)
		wantErr *errors.ErrorCode
	}{
		{
			name: "returns the user info",
			q:    &GetUserQuery{UserID: user.ID().String()},
			mock: func(repo *identitymock.MockUserRepository) {
				repo.EXPECT().FindByID(gomock.Any(), user.ID()).Return(user, nil)
			},
		},
		{
			name: "rejects an invalid user id",
			q:    &GetUserQuery{UserID: "not-a-uuid"},
			mock: func(repo *identitymock.MockUserRepository) {
				repo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: domain.ErrInvalidUser,
		},
		{
			name: "keeps not-found from the repository",
			q:    &GetUserQuery{UserID: user.ID().String()},
			mock: func(repo *identitymock.MockUserRepository) {
				repo.EXPECT().FindByID(gomock.Any(), user.ID()).
					Return(nil, errors.New(domain.ErrUserNotFound, "user not found"))
			},
			wantErr: domain.ErrUserNotFound,
		},
		{
			name: "maps any other repository failure to internal",
			q:    &GetUserQuery{UserID: user.ID().String()},
			mock: func(repo *identitymock.MockUserRepository) {
				repo.EXPECT().FindByID(gomock.Any(), user.ID()).Return(nil, assert.AnError)
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

			res, err := NewGetUser(repo).Handle(context.Background(), tt.q)

			if tt.wantErr != nil {
				assert.True(t, errors.Is(err, tt.wantErr), "got %v, want %v", err, tt.wantErr)
				assert.Nil(t, res)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, user.ID().String(), res.ID)
			assert.Equal(t, "jorge@example.com", res.Email)
		})
	}
}

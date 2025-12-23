package query

import (
	"context"
	"testing"

	"github.com/jorgeAM/go-template/internal/user/domain"
	"github.com/jorgeAM/go-template/internal/user/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type GetUserDependencies struct {
	userRepository *mock.MockUserRepository
}

func TestGetUser(t *testing.T) {
	type test struct {
		name   string
		userID string
		mock   func(deps *GetUserDependencies)
		err    bool
	}

	tests := []test{
		{
			name:   "get user successfully",
			userID: "123",
			mock: func(deps *GetUserDependencies) {
				mockUser := &domain.User{}

				deps.userRepository.EXPECT().
					FindByID(
						gomock.Any(),
						"123",
					).
					Return(mockUser, nil)
			},
			err: false,
		},
		{
			name:   "user id is empty",
			userID: "",
			mock: func(deps *GetUserDependencies) {
				deps.userRepository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Times(0)
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			deps := &GetUserDependencies{
				userRepository: mock.NewMockUserRepository(ctrl),
			}

			tt.mock(deps)
			srv := NewGetUser(deps.userRepository)

			user, err := srv.Exec(context.Background(), tt.userID)
			if tt.err {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
			}
		})
	}
}

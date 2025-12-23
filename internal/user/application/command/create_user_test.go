package command

import (
	"context"
	"testing"

	"github.com/jorgeAM/go-template/internal/user/domain"
	"github.com/jorgeAM/go-template/internal/user/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type CreateUserDependencies struct {
	userRepository *mock.MockUserRepository
}

func TestCreateUser(t *testing.T) {
	type test struct {
		name string
		cmd  CreateUserCommand
		mock func(deps *CreateUserDependencies)
		err  bool
	}

	tests := []test{
		{
			name: "create user",
			cmd: CreateUserCommand{
				Name:     "foo",
				Email:    "foo@baz@gmail.com",
				Password: "123456",
			},
			mock: func(deps *CreateUserDependencies) {
				var mockUser domain.User

				deps.userRepository.EXPECT().
					Save(
						gomock.Any(),
						gomock.AssignableToTypeOf(&mockUser),
					).
					Return(nil)
			},
			err: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			deps := &CreateUserDependencies{
				userRepository: mock.NewMockUserRepository(ctrl),
			}

			tt.mock(deps)
			srv := NewCreateUser(deps.userRepository)

			err := srv.Exec(context.Background(), &tt.cmd)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

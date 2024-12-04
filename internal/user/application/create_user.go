package application

import (
	"context"

	"github.com/jorgeAM/base-api/internal/user"
)

type CreateUserCommand struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUser struct {
	userRepository user.UserRepository
}

func NewCreateUser(userRepository user.UserRepository) *CreateUser {
	return &CreateUser{
		userRepository: userRepository,
	}
}

func (c *CreateUser) Exec(ctx context.Context, cmd *CreateUserCommand) error {
	user := user.NewUser(cmd.Name, cmd.Email, cmd.Password)

	return c.userRepository.Save(ctx, user)
}

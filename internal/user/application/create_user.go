package application

import (
	"context"

	"github.com/jorgeAM/go-template/internal/user/domain"
)

type CreateUserCommand struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUser struct {
	userRepository domain.UserRepository
}

func NewCreateUser(userRepository domain.UserRepository) *CreateUser {
	return &CreateUser{
		userRepository: userRepository,
	}
}

func (c *CreateUser) Exec(ctx context.Context, cmd *CreateUserCommand) error {
	user := domain.NewUser(cmd.Name, cmd.Email, cmd.Password)

	return c.userRepository.Save(ctx, user)
}

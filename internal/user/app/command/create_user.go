package command

import (
	"context"

	"github.com/jorgeAM/go-template/internal/shared/errors"
	"github.com/jorgeAM/go-template/internal/user/app/models"
	"github.com/jorgeAM/go-template/internal/user/domain"
)

type CreateUserCommand struct {
	Name     string
	Email    string
	Password string
}

type CreateUser struct {
	userRepository domain.UserRepository
}

func NewCreateUser(userRepository domain.UserRepository) *CreateUser {
	return &CreateUser{
		userRepository: userRepository,
	}
}

func (c *CreateUser) Handle(ctx context.Context, cmd *CreateUserCommand) (*models.UserInfo, error) {
	// NewUser already returns a structured domain error (ErrInvalidUser/ErrUserInternal)
	// whose message is safe to show the caller, so it is returned as is.
	user, err := domain.NewUser(cmd.Name, cmd.Email, cmd.Password)
	if err != nil {
		return nil, err
	}

	if err := c.userRepository.Save(ctx, user); err != nil {
		return nil, errors.Wrap(domain.ErrUserInternal, err, "we got a problem creating the user")
	}

	info := models.NewUserInfo(user)

	return &info, nil
}

package query

import (
	"context"
	"errors"

	"github.com/jorgeAM/go-template/internal/user/domain"
)

type GetUser struct {
	userRepository domain.UserRepository
}

func NewGetUser(userRepository domain.UserRepository) *GetUser {
	return &GetUser{
		userRepository: userRepository,
	}
}

func (g *GetUser) Exec(ctx context.Context, userID string) (*domain.User, error) {
	if userID == "" {
		return nil, errors.New("user id cannot be empty")
	}

	return g.userRepository.FindByID(ctx, userID)
}

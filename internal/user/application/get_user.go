package application

import (
	"context"

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
	return g.userRepository.FindByID(ctx, userID)
}

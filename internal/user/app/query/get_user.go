package query

import (
	"context"

	"github.com/jorgeAM/go-template/internal/shared/errors"
	"github.com/jorgeAM/go-template/internal/shared/model"
	"github.com/jorgeAM/go-template/internal/user/app/models"
	"github.com/jorgeAM/go-template/internal/user/domain"
)

type GetUserQuery struct {
	UserID string
}

type GetUser struct {
	userRepository domain.UserRepository
}

func NewGetUser(userRepository domain.UserRepository) *GetUser {
	return &GetUser{
		userRepository: userRepository,
	}
}

func (g *GetUser) Handle(ctx context.Context, q *GetUserQuery) (*models.UserInfo, error) {
	userID, err := model.NewID(q.UserID)
	if err != nil {
		return nil, errors.Wrap(domain.ErrInvalidUser, err, "invalid user id", errors.WithMetadata("id", q.UserID))
	}

	user, err := g.userRepository.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, errors.Wrap(domain.ErrUserNotFound, err, "user not found")
		}

		return nil, errors.Wrap(domain.ErrUserInternal, err, "we got a problem retrieving the user")
	}

	info := models.NewUserInfo(user)

	return &info, nil
}

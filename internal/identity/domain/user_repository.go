package domain

import (
	"context"

	"github.com/jorgeAM/go-template/internal/shared/model"
)

//go:generate go tool mockgen -source=./user_repository.go -destination=../mocks/user.go -package=mock -mock_names=UserRepository=MockUserRepository
type UserRepository interface {
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id model.ID) (*User, error)
}

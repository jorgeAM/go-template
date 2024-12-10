package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source=./user.go -destination=../mock/user.go -package=mock -mock_names=Repository=MockUserRepository
type UserRepository interface {
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
}

type User struct {
	ID       string
	Name     string
	Email    string
	Password string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func NewUser(name, email, password string) *User {
	return &User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

package user

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
}

type User struct {
	ID       string
	Name     string
	Email    string
	Password string
}

func NewUser(name, email, password string) *User {
	return &User{
		ID:       uuid.New().String(),
		Name:     name,
		Email:    email,
		Password: password,
	}
}

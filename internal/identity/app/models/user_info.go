package models

import (
	"time"

	"github.com/jorgeAM/go-template/internal/identity/domain"
)

type UserInfo struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUserInfo(user *domain.User) UserInfo {
	timestamps := user.Timestamps()

	return UserInfo{
		ID:        user.ID().String(),
		Name:      user.Name(),
		Email:     user.Email().String(),
		CreatedAt: timestamps.CreatedAt,
		UpdatedAt: timestamps.UpdatedAt,
	}
}

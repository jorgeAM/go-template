package domain

import (
	"strings"

	"github.com/jorgeAM/go-template/internal/shared/crypto"
	"github.com/jorgeAM/go-template/internal/shared/errors"
	"github.com/jorgeAM/go-template/internal/shared/valueobject"
)

const minPasswordLength = 8

var (
	ErrUserInternal = errors.Define("user.internal_error")
	ErrUserNotFound = errors.Define("user.not_found")
	ErrInvalidUser  = errors.Define("user.invalid")
)

type User struct {
	id             valueobject.ID
	name           string
	email          valueobject.Email
	hashedPassword string
	timestamps     valueobject.Timestamps
}

func NewUser(name, email, password string) (*User, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New(ErrInvalidUser, "name is required")
	}

	emailVO, err := valueobject.NewEmail(email)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidUser, err, "email is invalid")
	}

	if len(password) < minPasswordLength {
		return nil, errors.New(
			ErrInvalidUser,
			"password is too short",
			errors.WithMetadata("min_length", minPasswordLength),
		)
	}

	hashed, err := crypto.HashPassword(password)
	if err != nil {
		return nil, errors.Wrap(ErrUserInternal, err, "failed to hash password")
	}

	return &User{
		id:             valueobject.GenerateUUID(),
		name:           name,
		email:          emailVO,
		hashedPassword: hashed,
		timestamps:     valueobject.NewTimestamps(),
	}, nil
}

func (u *User) ID() valueobject.ID {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() valueobject.Email {
	return u.email
}

func (u *User) HashedPassword() string {
	return u.hashedPassword
}

func (u *User) Timestamps() valueobject.Timestamps {
	return u.timestamps
}

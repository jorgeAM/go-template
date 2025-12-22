package model

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrInvalidEmail = errors.New("invalid email")
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type Email string

func NewEmail(email string) (Email, error) {
	if len(email) == 0 {
		return "", ErrInvalidEmail
	}

	normalized := strings.ToLower(strings.TrimSpace(email))

	if !emailRegex.MatchString(normalized) {
		return "", ErrInvalidEmail
	}

	return Email(normalized), nil
}

func (e Email) String() string {
	return string(e)
}

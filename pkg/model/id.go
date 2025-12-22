package model

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidID = errors.New("invalid id")
)

type ID string

func NewID(id string) (ID, error) {
	err := uuid.Validate(id)
	if err != nil {
		return "", err
	}

	return ID(id), nil
}

func GenerateUUID() ID {
	return ID(uuid.New().String())
}

func (id ID) String() string {
	return string(id)
}

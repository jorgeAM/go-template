package valueobject

import (
	"github.com/google/uuid"
	"github.com/jorgeAM/go-template/internal/shared/errors"
)

var (
	ErrInvalidUUID    = errors.Define("uuid.invalid_uuid")
	ErrUUIDGeneration = errors.Define("uuid.generation_failed")
)

type UUID [16]byte

func NewUUIDv7() (UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return UUID{}, errors.Wrap(ErrUUIDGeneration, err, "failed to generate uuid")
	}

	return UUID(id), nil
}

func ParseUUID(s string) (UUID, error) {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return UUID{}, errors.Wrap(
			ErrInvalidUUID,
			err,
			"invalid uuid",
			errors.WithMetadata("uuid", s),
		)
	}

	return UUID(parsed), nil
}

func (u UUID) IsZero() bool {
	return uuid.UUID(u) == uuid.Nil
}

func (u UUID) Equals(other UUID) bool {
	return u == other
}

func (u UUID) String() string {
	return uuid.UUID(u).String()
}

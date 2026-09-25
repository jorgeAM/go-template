package domain

import "github.com/jorgeAM/go-template/internal/shared/model"

// UnmarshallUser rebuilds a persisted User without re-running NewUser's
// invariants (the password is already hashed). Only adapters/db calls it.
func UnmarshallUser(
	id model.ID,
	name string,
	email model.Email,
	hashedPassword string,
	timestamps model.Timestamps,
) *User {
	return &User{
		id:         id,
		name:       name,
		email:      email,
		password:   hashedPassword,
		timestamps: timestamps,
	}
}

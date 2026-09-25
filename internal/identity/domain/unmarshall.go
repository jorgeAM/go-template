package domain

import "github.com/jorgeAM/go-template/internal/shared/model"

func UnmarshallUser(
	id model.ID,
	name string,
	email model.Email,
	hashedPassword string,
	timestamps model.Timestamps,
) *User {
	return &User{
		id:             id,
		name:           name,
		email:          email,
		hashedPassword: hashedPassword,
		timestamps:     timestamps,
	}
}

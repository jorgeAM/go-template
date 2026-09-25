package domain

import "github.com/jorgeAM/go-template/internal/shared/valueobject"

func UnmarshallUser(
	id valueobject.ID,
	name string,
	email valueobject.Email,
	hashedPassword string,
	timestamps valueobject.Timestamps,
) *User {
	return &User{
		id:             id,
		name:           name,
		email:          email,
		hashedPassword: hashedPassword,
		timestamps:     timestamps,
	}
}

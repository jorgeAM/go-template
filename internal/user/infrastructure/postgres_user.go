package infrastructure

import (
	"time"

	"github.com/jorgeAM/go-template/internal/user/domain"
)

type postgresUser struct {
	ID        string     `db:"id"`
	Name      string     `db:"name"`
	Email     string     `db:"email"`
	Password  string     `db:"password"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" goqu:"omitnil"`
}

func (dto postgresUser) toDomain() (*domain.User, error) {
	return &domain.User{
		ID:       dto.ID,
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
	}, nil
}

func fromDomain(entity *domain.User) (*postgresUser, error) {
	return &postgresUser{
		ID:        entity.ID,
		Name:      entity.Name,
		Email:     entity.Email,
		Password:  entity.Password,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
	}, nil
}

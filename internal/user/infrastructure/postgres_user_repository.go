package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"

	"github.com/jorgeAM/go-template/internal/user/domain"
)

var _ domain.UserRepository = (*PostgresUserRepository)(nil)

type PostgresUserRepository struct {
	db     *sqlx.DB
	schema string
	table  string
}

func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db:     db,
		schema: "my_schema",
		table:  "users",
	}
}

func (p *PostgresUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	dialect := goqu.Dialect("postgres")
	ds := dialect.
		From(fmt.Sprintf("%s.%s", p.schema, p.table)).
		Where(goqu.Ex{
			"id": id,
		})

	sqlQuery, args, err := ds.Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}

	var dto postgresUser
	if err := p.db.GetContext(
		ctx,
		&dto,
		sqlQuery,
		args...,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("certificate not found")
		}

		return nil, err
	}

	return dto.toDomain()

}

func (p *PostgresUserRepository) Save(ctx context.Context, user *domain.User) error {
	dto, err := fromDomain(user)
	if err != nil {
		return err
	}

	ds := goqu.
		Insert("coder_v3_schema.users").
		Rows(dto).
		OnConflict(goqu.DoUpdate("id", goqu.Record{
			"name": dto.Name,
		}))

	sql, _, err := ds.ToSQL()
	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(ctx, sql)
	if err != nil {
		return err
	}

	return nil
}

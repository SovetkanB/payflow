package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SovetkanB/payflow/user-service/internal/model"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type postgresRepo struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &postgresRepo{
		db: db,
	}
}

func (pr *postgresRepo) CreateUser(ctx context.Context, user *model.User) error {
	err := pr.db.QueryRow(
		`INSERT INTO users (email, username, password)
		 VALUES ($1, $2, $3)
		 RETURNING id, created_at, updated_at`,
		user.Email, user.Username, user.Password,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (pr *postgresRepo) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User

	err := pr.db.Get(
		&user,
		`SELECT id, email, username, password, created_at, updated_at
		 FROM users
		 WHERE id = $1`,
		id,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, model.ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (pr *postgresRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	err := pr.db.Get(
		&user,
		`SELECT id, email, username, password, created_at, updated_at
		FROM users
		WHERE email = $1`,
		email,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, model.ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

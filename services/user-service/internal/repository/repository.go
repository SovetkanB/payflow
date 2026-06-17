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
	CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*model.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
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

func (pr *postgresRepo) CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error {
	err := pr.db.QueryRow(
		`INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`,
		token.UserID, token.Token, token.ExpiresAt,
	).Scan(&token.ID, &token.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (pr *postgresRepo) GetRefreshToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	var rt model.RefreshToken

	err := pr.db.Get(
		&rt,
		`SELECT id, user_id, token, expires_at, created_at
		FROM refresh_tokens
		WHERE token = $1`,
		token,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, model.ErrNotFound
		default:
			return nil, err
		}
	}

	return &rt, nil
}

func (pr *postgresRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	result, err := pr.db.ExecContext(
		ctx,
		`DELETE FROM refresh_tokens WHERE token = $1`,
		token,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return model.ErrNotFound
	}

	return nil
}

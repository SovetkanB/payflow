package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SovetkanB/payflow/user-service/internal/model"
	"github.com/SovetkanB/payflow/user-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, req model.CreateUserRequest) (*model.UserResponse, error) {
	if req.Email == "" || req.Username == "" || req.Password == "" {
		return nil, errors.New("email, username and password are required")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashed),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user.ToUserResponse(), nil
}

func (s *Service) GetUser(ctx context.Context, id string) (*model.UserResponse, error) {
	if id == "" {
		return nil, errors.New("invalid id")
	}

	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, model.ErrNotFound
		default:
			return nil, err
		}
	}

	return user.ToUserResponse(), nil
}

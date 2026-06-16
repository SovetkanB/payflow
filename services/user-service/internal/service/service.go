package service

import (
	"context"
	"errors"

	"github.com/SovetkanB/payflow/user-service/internal/auth"
	"github.com/SovetkanB/payflow/user-service/internal/config"
	"github.com/SovetkanB/payflow/user-service/internal/model"
	"github.com/SovetkanB/payflow/user-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo repository.Repository
	cfg  *config.Config
}

func NewService(repo repository.Repository, cfg *config.Config) *Service {
	return &Service{repo: repo, cfg: cfg}
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
		return nil, err
	}

	return user.ToUserResponse(), nil
}

func (s *Service) Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil || !(bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) == nil) {
		return nil, model.ErrInvalidPassword
	}

	token, err := auth.GenerateToken(user.ID, user.Email, s.cfg.JWTSecret, s.cfg.JWTExpiretion)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token: token,
		User:  user.ToUserResponse(),
	}, nil
}

package user

import (
	"context"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(ctx context.Context, username, password string) (User, error) {
	if username == "" {
		return User{}, ErrInvalidUsername
	}
	if password == "" {
		return User{}, ErrInvalidPassword
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	return s.repo.Create(ctx, username, string(passwordHash))
}

func (s *Service) Login(ctx context.Context, username, password string) (User, error) {
	if username == "" {
		return User{}, ErrInvalidUsername
	}
	if password == "" {
		return User{}, ErrInvalidPassword
	}
	user, exists, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return User{}, err
	}
	if !exists {
		return User{}, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return User{}, ErrInvalidCredentials
	}
	return user, nil
}

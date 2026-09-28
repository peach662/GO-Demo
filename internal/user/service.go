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

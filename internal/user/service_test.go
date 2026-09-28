package user

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type fakeRepository struct {
	users []User
}

func (f *fakeRepository) Create(ctx context.Context, username, passwordHash string) (User, error) {
	for _, item := range f.users {
		if item.Username == username {
			return User{}, ErrUsernameTaken
		}

	}
	user := User{
		ID:           len(f.users) + 1,
		Username:     username,
		PasswordHash: passwordHash,
	}
	f.users = append(f.users, user)

	return user, nil
}

func (f *fakeRepository) GetByUsername(ctx context.Context, username string) (User, bool, error) {
	for _, item := range f.users {
		if item.Username == username {
			return item, true, nil
		}
	}
	return User{}, false, nil
}

func TestServiceRegister(t *testing.T) {
	repo := &fakeRepository{}
	service := NewService(repo)

	user, err := service.Register(context.Background(), "testuser", "password")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.ID == 0 {
		t.Fatalf("expected user ID to be non-zero")
	}
	if user.Username != "testuser" {
		t.Fatalf("expected username to be testuser, got %s", user.Username)
	}

	if user.PasswordHash == "password" {
		t.Fatalf("expected password hash to be password, got %s", user.PasswordHash)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("password")) != nil {
		t.Fatalf("expected password hash to be password, got %s", user.PasswordHash)
	}
	_, err = service.Register(context.Background(), "", "password")
	if !errors.Is(err, ErrInvalidUsername) {
		t.Fatalf("expected ErrInvalidUsername, got %v", err)
	}
	_, err = service.Register(context.Background(), "testuser", "")
	if !errors.Is(err, ErrInvalidPassword) {
		t.Fatalf("expected ErrInvalidPassword, got %v", err)
	}
	_, err = service.Register(context.Background(), "testuser", "password")
	if !errors.Is(err, ErrUsernameTaken) {
		t.Fatalf("expected ErrUsernameTaken, got %v", err)
	}
	_, err = service.Register(context.Background(), "testuser", "")
	if !errors.Is(err, ErrInvalidPassword) {
		t.Fatalf("expected ErrInvalidPassword, got %v", err)
	}

}

func TestServiceLogin(t *testing.T) {
	repo := &fakeRepository{}
	service := NewService(repo)

	result,err:=service.Register(context.Background(), "testuser", "password")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	result, err = service.Login(context.Background(), "testuser", "password")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result.ID == 0 {
			t.Fatalf("expected user ID to be non-zero")
		}
		if result.Username != "testuser" {
			t.Fatalf("expected username to be testuser, got %s", result.Username)
		}
	
		_, err = service.Login(context.Background(), "testuser", "wrongpassword")
	
		if !errors.Is(err, ErrInvalidCredentials) {
			t.Fatalf("expected ErrInvalidCredentials, got %v", err)
		}

		_, err = service.Login(context.Background(), "abcs", "password")

		if !errors.Is(err, ErrInvalidCredentials) {
			t.Fatalf("expected ErrInvalidCredentials, got %v", err)
		}
		
}

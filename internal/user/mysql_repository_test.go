package user

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/database"
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	if err := godotenv.Load("../../.env"); err != nil && !os.IsNotExist(err) {
		t.Fatalf("load .env: %v", err)
	}
	if os.Getenv("MYSQL_DSN") == "" {
		if err := godotenv.Load("../../.env.example"); err != nil && !os.IsNotExist(err) {
			t.Fatalf("load .env.example: %v", err)
		}
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	db, err := database.OpenMySQL(cfg.MySQLDSN)
	if err != nil {
		t.Fatalf("open mysql: %v", err)

	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}

func TestMySQLRepositoryGetByUsername(t *testing.T) {
	db := openTestDB(t)
	repo := NewMySQLRepository(db)

	ctx := context.Background()
	username := fmt.Sprintf("user-%d", time.Now().UnixNano())
	passwordHash := fmt.Sprintf("password-%d", time.Now().UnixNano())

	_, err := db.ExecContext(ctx, "INSERT INTO users (username, password_hash) VALUES (?, ?)", username, passwordHash)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}

	t.Cleanup(func() {
		_, err := db.ExecContext(
			context.Background(),
			`DELETE FROM users WHERE username = ?`,
			username,
		)
		if err != nil {
			t.Errorf("delete test user: %v", err)
		}
	})

	user, found, err := repo.GetByUsername(ctx, username)
	if err != nil {
		t.Fatalf("get user: %v", err)
	}
	if !found {
		t.Fatalf("user not found")
	}
	if user.Username != username {
		t.Fatalf("user username mismatch: %s != %s", user.Username, username)
	}

	if user.PasswordHash != passwordHash {
		t.Fatalf("user password hash mismatch: got %s", user.PasswordHash)
	}
	_, found, err = repo.GetByUsername(ctx, username+"-missing")
	if err != nil {
		t.Fatalf("get missing user: %v", err)
	}
	if found {
		t.Fatalf("missing user found")
	}
}

package todo

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

func TestMySQLRepositoryList(t *testing.T) {
	db := openTestDB(t)

	ctx := context.Background()
	title := fmt.Sprintf("repository-list-test-%d", time.Now().UnixNano())

	result, err := db.ExecContext(
		ctx,
		`INSERT INTO todos (title, done) VALUES (?, ?)`,
		title,
		false,
	)
	if err != nil {
		t.Fatalf("insert test todo: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("get inserted todo ID: %v", err)
	}

	t.Cleanup(func() {
		_, err := db.ExecContext(
			context.Background(),
			`DELETE FROM todos WHERE id = ?`,
			id,
		)
		if err != nil {
			t.Errorf("delete test todo: %v", err)
		}
	})

	repo := NewMySQLRepository(db)

	todos, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("list todos: %v", err)
	}

	found := false
	for _, item := range todos {
		if item.ID == int(id) && item.Title == title {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected to find test todo %q", title)
	}
}

func TestMySQLRepositoryGetByID(t *testing.T) {
	db := openTestDB(t)
	var err error

	ctx := context.Background()
	title := fmt.Sprintf("repository-list-test-%d", time.Now().UnixNano())

	result, err := db.ExecContext(
		ctx,
		`INSERT INTO todos (title, done) VALUES (?, ?)`,
		title,
		false,
	)
	if err != nil {
		t.Fatalf("insert test todo: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("get inserted todo ID: %v", err)
	}

	t.Cleanup(func() {
		_, err := db.ExecContext(
			context.Background(),
			`DELETE FROM todos WHERE id = ?`,
			id,
		)
		if err != nil {
			t.Errorf("delete test todo: %v", err)
		}
	})

	repo := NewMySQLRepository(db)

	item, found, err := repo.GetByID(ctx, int(id))
	if err != nil {
		t.Fatalf("get todo by ID: %v", err)
	}

	if !found {
		t.Fatalf("expected todo to be found")
	}
	if item.ID != int(id) {
		t.Errorf("expected ID %d, got %d", id, item.ID)
	}
	if item.Title != title {
		t.Errorf("expected title %q, got %q", title, item.Title)
	}
	_, found, err = repo.GetByID(ctx, 999999999)

	if err != nil {
		t.Fatalf("get missing todo: %v", err)
	}
	if found {
		t.Errorf("expected todo to be not found")
	}

}

func TestMySQLRepositoryGetByIDDatabaseError(t *testing.T) {
	db := openTestDB(t)
	_ = db.Close()

	repo := NewMySQLRepository(db)

	_, found, err := repo.GetByID(
		context.Background(),
		1,
	)

	if err == nil {
		t.Fatal("expected database error")
	}
	if found {
		t.Error("expected found to be false")
	}
}

func TestMySQLRepositoryCreate(t *testing.T) {
	db := openTestDB(t)

	ctx := context.Background()
	title := fmt.Sprintf(
		"repository-create-test-%d",
		time.Now().UnixNano(),
	)

	repo := NewMySQLRepository(db)

	created, err := repo.Create(context.Background(), title)
	if err != nil {
		t.Fatalf("create todo: %v", err)
	}

	t.Cleanup(func() {
		_, err := db.ExecContext(
			context.Background(),
			`DELETE FROM todos WHERE id = ?`,
			created.ID,
		)
		if err != nil {
			t.Errorf("delete test todo: %v", err)
		}
	})

	if created.ID <= 0 {
		t.Errorf("expected ID to be positive, got %d", created.ID)
	}
	if created.Title != title {
		t.Errorf("expected title %q, got %q", title, created.Title)
	}
	if created.Done {
		t.Errorf("expected done to be false")
	}

	item, found, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("get todo by ID: %v", err)
	}
	if !found {
		t.Fatalf("expected todo to be found")
	}
	if item.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, item.ID)
	}
	if item.Title != title {
		t.Errorf("expected title %q, got %q", title, item.Title)
	}
	if item.Done {
		t.Errorf("expected stored done to be false")
	}
}

func TestMySQLRepositoryUpdateStatus(t *testing.T) {
	db := openTestDB(t)

	ctx := context.Background()
	title := fmt.Sprintf(
		"repository-create-test-%d",
		time.Now().UnixNano(),
	)

	repo := NewMySQLRepository(db)

	created, err := repo.Create(ctx, title)
	if err != nil {
		t.Fatalf("create todo: %v", err)
	}

	t.Cleanup(func() {
		_, err := db.ExecContext(
			context.Background(),
			`DELETE FROM todos WHERE id = ?`,
			created.ID,
		)
		if err != nil {
			t.Errorf("delete test todo: %v", err)
		}
	})

	updated, found, err := repo.UpdateStatus(ctx, created.ID, true)

	if err != nil {
		t.Fatalf("update todo status: %v", err)
	}
	if !found {
		t.Fatalf("expected todo to be found")
	}
	if !updated.Done {
		t.Errorf("expected done to be true")
	}
	if updated.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, updated.ID)
	}
	if updated.Title != title {
		t.Errorf("expected title %q, got %q", title, updated.Title)
	}

	item, found, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("get todo by ID: %v", err)
	}
	if !found {
		t.Fatalf("expected todo to be found")
	}
	if !item.Done {
		t.Errorf("expected stored done to be true")
	}

	_, found, err = repo.UpdateStatus(ctx, 999999999, true)
	if err != nil {
		t.Fatalf("update missing todo: %v", err)
	}
	if found {
		t.Errorf("expected missing todo to be not found")
	}
}

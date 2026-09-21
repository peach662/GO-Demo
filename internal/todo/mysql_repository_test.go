package todo

import (
	"awesomeProject/internal/database"
	"context"
	"fmt"
	"testing"
	"time"
)

func TestMySQLRepositoryList(t *testing.T) {
	db, err := database.OpenMySQL()
	if err != nil {
		t.Fatalf("open mysql: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

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

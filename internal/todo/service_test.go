package todo

import (
	"context"
	"testing"
)

type fakeRepository struct {
	todos []Todo
}

func (f *fakeRepository) List(ctx context.Context) ([]Todo, error) {
	result := make([]Todo, len(f.todos))
	copy(result, f.todos)
	return result, nil
}

func (f *fakeRepository) GetByID(ctx context.Context, id int) (Todo, bool, error) {
	for _, item := range f.todos {
		if item.ID == id {
			return item, true, nil
		}
	}
	return Todo{}, false, nil
}

func (f *fakeRepository) Create(ctx context.Context, title string) (Todo, error) {
	item := Todo{
		ID:     len(f.todos) + 1,
		Title:  title,
		Status: StatusPending,
	}
	f.todos = append(f.todos, item)
	return item, nil
}

func (f *fakeRepository) UpdateStatus(
	ctx context.Context,
	id int,
	status Status,
) (Todo, bool, error) {
	for i := range f.todos {
		if f.todos[i].ID == id {
			f.todos[i].Status = status
			return f.todos[i], true, nil
		}
	}
	return Todo{}, false, nil
}

func TestServiceCreate(t *testing.T) {
	repo := &fakeRepository{
		todos: []Todo{
			{ID: 1, Title: "Learn Go", Status: StatusPending},
			{ID: 2, Title: "Build a web app", Status: StatusPending},
		},
	}

	service := NewService(repo)
	created, err := service.Create(context.Background(), "Write tests")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if created.ID != 3 {
		t.Errorf("Expected ID 3, got %d", created.ID)
	}
	if created.Title != "Write tests" {
		t.Errorf("Expected title 'Write tests', got '%s'", created.Title)
	}
	if created.Status != StatusPending {
		t.Errorf("Expected status %q, got %q", StatusPending, created.Status)
	}
	list, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("list todos: %v", err)
	}

	if len(list) != 3 {
		t.Errorf("expected 3 todos, got %d", len(list))
	}

}

func TestServiceUpdateStatus(t *testing.T) {
	repo := &fakeRepository{
		todos: []Todo{
			{ID: 1, Title: "Learn Go", Status: StatusPending},
			{ID: 2, Title: "Build a web app", Status: StatusPending},
		},
	}

	service := NewService(repo)

	updated, found, err := service.UpdateStatus(context.Background(), 1, StatusProcessing)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !found {
		t.Fatalf("Expected to find todo with ID 1")
	}
	if updated.Status != StatusProcessing {
		t.Errorf("Expected status %q, got %q", StatusProcessing, updated.Status)
	}
	stored, found, err := service.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !found {
		t.Fatalf("expected to find todo with ID 1")
	}
	if stored.Status != StatusProcessing {
		t.Errorf("expected stored todo status %q", StatusProcessing)
	}
	_, found, err = service.UpdateStatus(context.Background(), 999, StatusProcessing)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if found {
		t.Errorf("Expected not to find todo with ID 999")

	}
}

func TestServiceRejectsInvalidStatusTransition(t *testing.T) {
	repo := &fakeRepository{
		todos: []Todo{{ID: 1, Title: "Learn Go", Status: StatusPending}},
	}

	service := NewService(repo)
	_, found, err := service.UpdateStatus(context.Background(), 1, StatusCompleted)
	if err != ErrInvalidTransition {
		t.Fatalf("expected invalid transition error, got %v", err)
	}
	if found {
		t.Fatal("expected invalid transition not to report found")
	}

	item, found, err := service.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("get todo: %v", err)
	}
	if !found {
		t.Fatal("expected todo to exist")
	}
	if item.Status != StatusPending {
		t.Errorf("expected status %q, got %q", StatusPending, item.Status)
	}
}

func TestServiceGetByID(t *testing.T) {
	repo := &fakeRepository{
		todos: []Todo{
			{ID: 1, Title: "Learn Go", Status: StatusPending},
			{ID: 2, Title: "Build a web app", Status: StatusPending},
		},
	}

	service := NewService(repo)
	item, found, err := service.GetByID(context.Background(), 2)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !found {
		t.Fatalf("Expected to find item with ID 2")
	}
	if item.ID != 2 {
		t.Errorf("Expected ID 2, got %d", item.ID)
	}

	_, found, err = service.GetByID(
		context.Background(),
		999,
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if found {
		t.Errorf("Expected not to find item with ID 999")
	}
}

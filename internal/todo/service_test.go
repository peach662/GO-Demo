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
		ID:    len(f.todos) + 1,
		Title: title,
		Done:  false,
	}
	f.todos = append(f.todos, item)
	return item, nil
}

func (f *fakeRepository) UpdateStatus(
	ctx context.Context,
	id int,
	done bool,
) (Todo, bool, error) {
	for i := range f.todos {
		if f.todos[i].ID == id {
			f.todos[i].Done = done
			return f.todos[i], true, nil
		}
	}
	return Todo{}, false, nil
}

func TestServiceCreate(t *testing.T) {
	repo := &fakeRepository{
		todos: []Todo{
			{ID: 1, Title: "Learn Go", Done: false},
			{ID: 2, Title: "Build a web app", Done: false},
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
	if created.Done != false {
		t.Errorf("Expected Done false, got %v", created.Done)
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
			{ID: 1, Title: "Learn Go", Done: false},
			{ID: 2, Title: "Build a web app", Done: false},
		},
	}

	service := NewService(repo)

	updated, found, err := service.UpdateStatus(context.Background(), 1, true)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !found {
		t.Fatalf("Expected to find todo with ID 1")
	}
	if updated.Done != true {
		t.Errorf("Expected Done true, got %v", updated.Done)
	}
	stored, found, err := service.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !found {
		t.Fatalf("expected to find todo with ID 1")
	}
	if !stored.Done {
		t.Errorf("expected stored todo Done to be true")
	}
	_, found, err = service.UpdateStatus(context.Background(), 999, true)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if found {
		t.Errorf("Expected not to find todo with ID 999")

	}
}

func TestServiceGetByID(t *testing.T) {
	repo := &fakeRepository{
		todos: []Todo{
			{ID: 1, Title: "Learn Go", Done: false},
			{ID: 2, Title: "Build a web app", Done: false},
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

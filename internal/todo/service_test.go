package todo

import (
	"sync"
	"testing"
)

func TestServiceCreate(t *testing.T) {
	service := NewService([]Todo{
		{ID: 1, Title: "Learn Go", Done: false},
		{ID: 2, Title: "Build a web app", Done: false},
	})
	created := service.Create("Write tests")

	if created.ID != 3 {
		t.Errorf("Expected ID 3, got %d", created.ID)
	}
	if created.Title != "Write tests" {
		t.Errorf("Expected title 'Write tests', got '%s'", created.Title)
	}
	if created.Done != false {
		t.Errorf("Expected Done false, got %v", created.Done)
	}
	if got := len(service.List()); got != 3 {
		t.Errorf("expected 3 todos, got %d", got)
	}

}

func TestServiceUpdateStatus(t *testing.T) {
	service := NewService([]Todo{
		{ID: 1, Title: "Learn Go", Done: false},
		{ID: 2, Title: "Build a web app", Done: false},
	})

	updated, found := service.UpdateStatus(1, true)
	if !found {
		t.Fatalf("Expected to find todo with ID 1")
	}
	if updated.Done != true {
		t.Errorf("Expected Done true, got %v", updated.Done)
	}
	stored, found := service.GetByID(1)
	if !found {
		t.Fatalf("expected to find todo with ID 1")
	}
	if !stored.Done {
		t.Errorf("expected stored todo Done to be true")
	}
	_, found = service.UpdateStatus(999, true)
	if found {
		t.Errorf("Expected not to find todo with ID 999")

	}
}

func TestServiceGetByID(t *testing.T) {
	service := NewService([]Todo{
		{ID: 1, Title: "Learn Go", Done: false},
		{ID: 2, Title: "Build a web app", Done: false},
	})

	item, found := service.GetByID(2)

	if !found {
		t.Fatalf("Expected to find item with ID 2")
	}
	if item.ID != 2 {
		t.Errorf("Expected ID 2, got %d", item.ID)
	}
	if item.Title != "Build a web app" {
		t.Errorf("Expected title 'Build a web app', got '%s'", item.Title)
	}
	_, found = service.GetByID(999)
	if found {
		t.Errorf("Expected not to find item with ID 999")
	}
}

func TestServiceListReturnsCopy(t *testing.T) {
	service := NewService([]Todo{
		{ID: 1, Title: "Learn Go", Done: false},
	})

	list := service.List()
	list[0].Title = "Changed outside"

	stored, found := service.GetByID(1)
	if !found {
		t.Fatal("expected todo with ID 1 to exist")
	}
	if stored.Title != "Learn Go" {
		t.Errorf("expected internal title %q, got %q", "Learn Go", stored.Title)
	}
}

func TestServiceConcurrentCreate(t *testing.T) {
	service := NewService(nil)

	const workers = 100
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			service.Create("concurrent todo")
		}()
	}

	wg.Wait()

	if got := len(service.List()); got != workers {
		t.Errorf("expected %d todos, got %d", workers, got)
	}

	ids := make(map[int]bool)

	for _, item := range service.List() {
		if ids[item.ID] {
			t.Errorf("duplicate todo ID: %d", item.ID)
		}
		ids[item.ID] = true
	}
}

func TestServiceCreateWithChannel(t *testing.T) {
	service := NewService(nil)
	resultCh := make(chan Todo)

	go func() {
		resultCh <- service.Create("created by goroutine")
	}()

	created := <-resultCh

	if created.Title != "created by goroutine" {
		t.Errorf("expected title %q, got %q", "created by goroutine", created.Title)
	}
}

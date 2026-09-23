package main

import (
	"awesomeProject/internal/todo"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeRepository struct {
	todos []todo.Todo
}

func newTestService(initialTodos []todo.Todo) *todo.Service {
	return todo.NewService(&fakeRepository{
		todos: initialTodos,
	})
}

func (r *fakeRepository) List(ctx context.Context) ([]todo.Todo, error) {
	result := make([]todo.Todo, len(r.todos))
	copy(result, r.todos)
	return result, nil
}

func (r *fakeRepository) GetByID(ctx context.Context, id int) (todo.Todo, bool, error) {
	for _, t := range r.todos {
		if t.ID == id {
			return t, true, nil
		}
	}
	return todo.Todo{}, false, nil
}

func (r *fakeRepository) Create(ctx context.Context, title string) (todo.Todo, error) {
	id := len(r.todos) + 1
	t := todo.Todo{ID: id, Title: title, Status: todo.StatusPending}
	r.todos = append(r.todos, t)
	return t, nil
}

func (r *fakeRepository) UpdateStatus(ctx context.Context, id int, status todo.Status) (todo.Todo, bool, error) {
	for i := range r.todos {
		if r.todos[i].ID == id {
			r.todos[i].Status = status
			return r.todos[i], true, nil
		}
	}

	return todo.Todo{}, false, nil
}
func TestHealthRoute(t *testing.T) {

	service := newTestService([]todo.Todo{
		{ID: 1, Title: "Learn Go", Status: todo.StatusPending},
		{ID: 2, Title: "Build a web app", Status: todo.StatusPending},
		{ID: 3, Title: "Deploy to production", Status: todo.StatusPending},
	})

	router := newRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if body := recorder.Body.String(); !strings.Contains(body, "pong") {
		t.Errorf("handler returned unexpected body: got %v want %v",
			body, "pong")
	}
}

func TestGetTodosRoute(t *testing.T) {
	service := newTestService([]todo.Todo{
		{ID: 1, Title: "Learn Go", Status: todo.StatusPending},
		{ID: 2, Title: "Build a web app", Status: todo.StatusPending},
	})

	router := newRouter(service)
	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response struct {
		Code int         `json:"code"`
		Data []todo.Todo `json:"data"`
	}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if len(response.Data) != 2 {
		t.Fatalf("Expected 2 todos, got %d", len(response.Data))
	}
	if response.Data[0].ID != 1 {
		t.Errorf("expected first todo ID 1, got %d", response.Data[0].ID)
	}

	if response.Code != 0 {
		t.Errorf("expected code 0, got %d", response.Code)
	}
	if response.Data[0].Title != "Learn Go" {
		t.Errorf("Expected first todo title 'Learn Go', got '%s'", response.Data[0].Title)
	}
}

func TestGetTodoNotFoundRoute(t *testing.T) {
	service := newTestService([]todo.Todo{
		{ID: 1, Title: "Learn Go", Status: todo.StatusPending},
		{ID: 2, Title: "Build a web app", Status: todo.StatusPending},
	})
	router := newRouter(service)
	req := httptest.NewRequest(http.MethodGet, "/todos/999", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
	var response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if response.Code != 404 {
		t.Errorf("expected code 404, got %d", response.Code)
	}
	if response.Message != "todo 不存在" {
		t.Errorf("expected message %q, got %q", "todo 不存在", response.Message)
	}
}

func TestCreateTodoRoute(t *testing.T) {
	service := newTestService([]todo.Todo{
		{ID: 1, Title: "Learn Go", Status: todo.StatusPending},
		{ID: 2, Title: "Build a web app", Status: todo.StatusPending},
	})
	router := newRouter(service)
	reqBody := `{"title":"Write HTTP tests"}`

	req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response struct {
		Code int       `json:"code"`
		Data todo.Todo `json:"data"`
	}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Data.ID != 3 {
		t.Errorf("Expected new todo ID 3, got %d", response.Data.ID)
	}
	if response.Data.Title != "Write HTTP tests" {
		t.Errorf("Expected new todo title 'Write HTTP tests', got '%s'", response.Data.Title)
	}
	if response.Code != 0 {
		t.Errorf("expected code 0, got %d", response.Code)
	}
	if response.Data.Status != todo.StatusPending {
		t.Errorf("Expected new todo status %q, got %q", todo.StatusPending, response.Data.Status)
	}
	items, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("list todos: %v", err)
	}
	if len(items) != 3 {
		t.Errorf("expected 3 todos, got %d", len(items))
	}
}

func TestCreateTodoValidationError(t *testing.T) {
	service := newTestService([]todo.Todo{
		{ID: 1, Title: "Learn Go", Status: todo.StatusPending},
		{ID: 2, Title: "Build a web app", Status: todo.StatusPending},
	})
	router := newRouter(service)
	reqBody := `{}`
	req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	var response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if response.Code != 400 {
		t.Errorf("expected code 400, got %d", response.Code)
	}
	if response.Message != "请求参数错误" {
		t.Errorf("expected message %q, got %q", "请求参数错误", response.Message)
	}

	items, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("list todos: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 todos, got %d", len(items))
	}
}

func TestUpdateTodoStatusRoute(t *testing.T) {

	service := newTestService([]todo.Todo{
		{ID: 1, Title: "Learn Go", Status: todo.StatusPending},
		{ID: 2, Title: "Build a web app", Status: todo.StatusPending},
	})
	router := newRouter(service)
	reqBody := `{"status":"PROCESSING"}`
	req := httptest.NewRequest(http.MethodPatch, "/todos/1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response struct {
		Code int       `json:"code"`
		Data todo.Todo `json:"data"`
	}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Data.ID != 1 {
		t.Errorf("Expected updated todo ID 1, got %d", response.Data.ID)
	}
	if response.Data.Status != todo.StatusProcessing {
		t.Errorf("Expected updated todo status %q, got %q", todo.StatusProcessing, response.Data.Status)
	}
	if response.Code != 0 {
		t.Errorf("expected code 0, got %d", response.Code)
	}

	stored, found, err := service.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("get todo: %v", err)
	}
	if !found {
		t.Fatalf("expected to find todo with ID 1")
	}
	if stored.Status != todo.StatusProcessing {
		t.Errorf("expected stored todo status %q", todo.StatusProcessing)
	}
}

func TestUpdateTodoStatusNotFoundRoute(t *testing.T) {
	service := newTestService([]todo.Todo{
		{ID: 1, Title: "Learn Go", Status: todo.StatusPending},
		{ID: 2, Title: "Build a web app", Status: todo.StatusPending},
	})
	router := newRouter(service)
	reqBody := `{"status":"PROCESSING"}`
	req := httptest.NewRequest(http.MethodPatch, "/todos/999", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
	var response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if response.Code != 404 {
		t.Errorf("expected code 404, got %d", response.Code)
	}
	if response.Message != "todo 不存在" {
		t.Errorf("expected message %q, got %q", "todo 不存在", response.Message)
	}
	items, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("list todos: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 todos in service, got %d", len(items))
	}
}

func TestUpdateTodoInvalidTransitionRoute(t *testing.T) {
	service := newTestService([]todo.Todo{
		{ID: 1, Title: "Completed todo", Status: todo.StatusCompleted},
	})
	router := newRouter(service)
	reqBody := `{"status":"PROCESSING"}`
	req := httptest.NewRequest(http.MethodPatch, "/todos/1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	var response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if response.Code != 400 {
		t.Errorf("expected code 400, got %d", response.Code)
	}
	if response.Message != "非法状态流转" {
		t.Errorf("expected message %q, got %q", "非法状态流转", response.Message)
	}
}

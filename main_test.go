package main

import (
	"awesomeProject/internal/todo"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthRoute(t *testing.T) {

	service := todo.NewService([]todo.Todo{
		{ID: 1, Title: "Learn Go", Done: false},
		{ID: 2, Title: "Build a web app", Done: false},
		{ID: 3, Title: "Deploy to production", Done: false},
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
	service := todo.NewService([]todo.Todo{
		{ID: 1, Title: "Learn Go", Done: false},
		{ID: 2, Title: "Build a web app", Done: false},
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
	service := todo.NewService([]todo.Todo{
		{ID: 1, Title: "Learn Go", Done: false},
		{ID: 2, Title: "Build a web app", Done: false},
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
	service := todo.NewService([]todo.Todo{
		{ID: 1, Title: "Learn Go", Done: false},
		{ID: 2, Title: "Build a web app", Done: false},
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
	if response.Data.Done != false {
		t.Errorf("Expected new todo Done false, got %v", response.Data.Done)
	}
	if got := len(service.List()); got != 3 {
		t.Errorf("expected 3 todos in service, got %d", got)
	}
}

func TestCreateTodoValidationError(t *testing.T) {
	service := todo.NewService([]todo.Todo{
		{ID: 1, Title: "Learn Go", Done: false},
		{ID: 2, Title: "Build a web app", Done: false},
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

	if got := len(service.List()); got != 2 {
		t.Errorf("expected 2 todos in service, got %d", got)
	}
}

func TestUpdateTodoStatusRoute(t *testing.T) {

	service := todo.NewService([]todo.Todo{
		{ID: 1, Title: "Learn Go", Done: false},
		{ID: 2, Title: "Build a web app", Done: false},
	})
	router := newRouter(service)
	reqBody := `{"done":true}`
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
	if response.Data.Done != true {
		t.Errorf("Expected updated todo Done true, got %v", response.Data.Done)
	}
	if response.Code != 0 {
		t.Errorf("expected code 0, got %d", response.Code)
	}

	stored, found := service.GetByID(1)
	if !found {
		t.Fatalf("expected to find todo with ID 1")
	}
	if !stored.Done {
		t.Errorf("expected stored todo Done to be true")
	}
}

func TestUpdateTodoStatusNotFoundRoute(t *testing.T) {
	service := todo.NewService([]todo.Todo{
		{ID: 1, Title: "Learn Go", Done: false},
		{ID: 2, Title: "Build a web app", Done: false},
	})
	router := newRouter(service)
	reqBody := `{"done":true}`
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
	if got := len(service.List()); got != 2 {
		t.Errorf("expected 2 todos in service, got %d", got)
	}
}

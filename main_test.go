package main

import (
	"awesomeProject/internal/auth"
	"awesomeProject/internal/todo"
	"awesomeProject/internal/user"
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

type fakeUserRepository struct {
	users []user.User
}

func newTestUserService(initialUsers []user.User) *user.Service {
	return user.NewService(&fakeUserRepository{
		users: initialUsers,
	})
}

func (r *fakeUserRepository) Create(ctx context.Context, username, passwordHash string) (user.User, error) {

	for _, existing := range r.users {
		if existing.Username == username {
			return user.User{}, user.ErrUsernameTaken
		}
	}
	id := len(r.users) + 1
	item := user.User{ID: id, Username: username, PasswordHash: passwordHash}

	r.users = append(r.users, item)
	return item, nil
}

func (r *fakeUserRepository) GetByUsername(ctx context.Context, username string) (user.User, bool, error) {
	for _, item := range r.users {
		if item.Username == username {
			return item, true, nil
		}
	}
	return user.User{}, false, nil
}

const testTodoOwnerID = 1

func newTestService(initialTodos []todo.Todo) *todo.Service {
	return todo.NewService(&fakeRepository{
		todos: initialTodos,
	}, nil)
}

func bearerToken(t *testing.T, jwtService *auth.JWT, userID int) string {
	t.Helper()
	token, err := jwtService.GenerateToken(userID)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	return token
}

func (r *fakeRepository) List(ctx context.Context, userID int) ([]todo.Todo, error) {
	result := make([]todo.Todo, 0)
	for _, item := range r.todos {
		if item.UserID == userID {
			result = append(result, item)
		}
	}
	return result, nil
}

func (r *fakeRepository) GetByID(ctx context.Context, id int, userID int) (todo.Todo, bool, error) {
	for _, t := range r.todos {
		if t.ID == id && t.UserID == userID {
			return t, true, nil
		}
	}
	return todo.Todo{}, false, nil
}

func (r *fakeRepository) Create(ctx context.Context, title string, userID int) (todo.Todo, error) {
	id := len(r.todos) + 1
	t := todo.Todo{ID: id, Title: title, Status: todo.StatusPending, UserID: userID}
	r.todos = append(r.todos, t)
	return t, nil
}

func (r *fakeRepository) UpdateStatus(ctx context.Context, id int, status todo.Status, userID int) (todo.Todo, bool, error) {
	for i := range r.todos {
		if r.todos[i].ID == id && r.todos[i].UserID == userID {
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

	router := newRouter(service, user.NewService(nil), auth.NewJWT("test"), nil)

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
		{ID: 1, Title: "Learn Go", Status: todo.StatusPending, UserID: testTodoOwnerID},
		{ID: 2, Title: "Build a web app", Status: todo.StatusPending, UserID: testTodoOwnerID},
		{ID: 3, Title: "Someone else", Status: todo.StatusPending, UserID: 99},
	})
	jwtService := auth.NewJWT("test")
	router := newRouter(service, user.NewService(nil), jwtService, nil)
	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	req.Header.Set("Authorization", "Bearer "+bearerToken(t, jwtService, testTodoOwnerID))
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
		{ID: 1, Title: "Learn Go", Status: todo.StatusPending, UserID: testTodoOwnerID},
		{ID: 2, Title: "Build a web app", Status: todo.StatusPending, UserID: 99},
	})
	jwtService := auth.NewJWT("test")
	router := newRouter(service, user.NewService(nil), jwtService, nil)
	req := httptest.NewRequest(http.MethodGet, "/todos/2", nil)
	req.Header.Set("Authorization", "Bearer "+bearerToken(t, jwtService, testTodoOwnerID))
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

func TestRegisterRoute(t *testing.T) {
	todoService := newTestService(nil)
	userService := newTestUserService(nil)
	router := newRouter(todoService, userService, auth.NewJWT("test-secret"), nil)
	reqBody := `{"username":"alice","password":"secret"}`
	req := httptest.NewRequest(http.MethodPost, "/users/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	body := recorder.Body.String()
	if !strings.Contains(body, "alice") {
		t.Errorf("expected body %q, got %q", "alice", body)
	}
	if strings.Contains(body, "secret") {
		t.Errorf("expected body %q, got %q", "secret", body)
	}
	if strings.Contains(body, "password_hash") {
		t.Errorf("expected body %q, got %q", "password_hash", body)
	}
}

func TestRegisterDuplicateRoute(t *testing.T) {

	todoService := newTestService(nil)
	userService := newTestUserService([]user.User{
		{ID: 1, Username: "alice", PasswordHash: "secret"},
	})
	router := newRouter(todoService, userService, auth.NewJWT("test-secret"), nil)
	reqBody := `{"username":"alice","password":"secret"}`
	req := httptest.NewRequest(http.MethodPost, "/users/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}
func TestCreateTodoRoute(t *testing.T) {
	service := newTestService([]todo.Todo{
		{ID: 1, Title: "Learn Go", Status: todo.StatusPending, UserID: testTodoOwnerID},
		{ID: 2, Title: "Build a web app", Status: todo.StatusPending, UserID: testTodoOwnerID},
	})
	jwtService := auth.NewJWT("test")
	router := newRouter(service, user.NewService(nil), jwtService, nil)
	reqBody := `{"title":"Write HTTP tests"}`
	req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken(t, jwtService, testTodoOwnerID))
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
	if response.Data.UserID != testTodoOwnerID {
		t.Errorf("expected user ID %d, got %d", testTodoOwnerID, response.Data.UserID)
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
	items, err := service.List(context.Background(), testTodoOwnerID)
	if err != nil {
		t.Fatalf("list todos: %v", err)
	}
	if len(items) != 3 {
		t.Errorf("expected 3 todos, got %d", len(items))
	}
}

func TestCreateTodoValidationError(t *testing.T) {
	service := newTestService([]todo.Todo{
		{ID: 1, Title: "Learn Go", Status: todo.StatusPending, UserID: testTodoOwnerID},
		{ID: 2, Title: "Build a web app", Status: todo.StatusPending, UserID: testTodoOwnerID},
	})
	jwtService := auth.NewJWT("test")
	router := newRouter(service, user.NewService(nil), jwtService, nil)
	reqBody := `{}`
	req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken(t, jwtService, testTodoOwnerID))
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

	items, err := service.List(context.Background(), testTodoOwnerID)
	if err != nil {
		t.Fatalf("list todos: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 todos, got %d", len(items))
	}
}

func TestCreateTodoUnauthorized(t *testing.T) {
	service := newTestService([]todo.Todo{
		{ID: 1, Title: "Learn Go", Status: todo.StatusPending, UserID: testTodoOwnerID},
		{ID: 2, Title: "Build a web app", Status: todo.StatusPending, UserID: testTodoOwnerID},
	})
	router := newRouter(service, user.NewService(nil), auth.NewJWT("test"), nil)
	reqBody := `{"title":"Write HTTP tests"}`
	req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
	items, err := service.List(context.Background(), testTodoOwnerID)
	if err != nil {
		t.Fatalf("list todos: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 todos, got %d", len(items))
	}
}
func TestUpdateTodoStatusRoute(t *testing.T) {

	service := newTestService([]todo.Todo{
		{ID: 1, Title: "Learn Go", Status: todo.StatusPending, UserID: testTodoOwnerID},
		{ID: 2, Title: "Build a web app", Status: todo.StatusPending, UserID: testTodoOwnerID},
	})
	jwtService := auth.NewJWT("test")
	router := newRouter(service, user.NewService(nil), jwtService, nil)
	reqBody := `{"status":"PROCESSING"}`
	req := httptest.NewRequest(http.MethodPatch, "/todos/1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken(t, jwtService, testTodoOwnerID))
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

	stored, found, err := service.GetByID(context.Background(), 1, testTodoOwnerID)
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
		{ID: 1, Title: "Learn Go", Status: todo.StatusPending, UserID: testTodoOwnerID},
		{ID: 2, Title: "Build a web app", Status: todo.StatusPending, UserID: testTodoOwnerID},
	})
	jwtService := auth.NewJWT("test")
	router := newRouter(service, user.NewService(nil), jwtService, nil)
	reqBody := `{"status":"PROCESSING"}`
	req := httptest.NewRequest(http.MethodPatch, "/todos/999", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken(t, jwtService, testTodoOwnerID))
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
	items, err := service.List(context.Background(), testTodoOwnerID)
	if err != nil {
		t.Fatalf("list todos: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 todos in service, got %d", len(items))
	}
}

func TestUpdateTodoInvalidTransitionRoute(t *testing.T) {
	service := newTestService([]todo.Todo{
		{ID: 1, Title: "Completed todo", Status: todo.StatusCompleted, UserID: testTodoOwnerID},
	})
	jwtService := auth.NewJWT("test")
	router := newRouter(service, user.NewService(nil), jwtService, nil)
	reqBody := `{"status":"PROCESSING"}`
	req := httptest.NewRequest(http.MethodPatch, "/todos/1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken(t, jwtService, testTodoOwnerID))
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

func TestLoginRoute(t *testing.T) {
	todoService := newTestService(nil)
	userService := newTestUserService(nil)
	router := newRouter(todoService, userService, auth.NewJWT("test-secret"), nil)
	reqBody := `{"username":"alice","password":"secret"}`
	req := httptest.NewRequest(http.MethodPost, "/users/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	body := recorder.Body.String()
	if !strings.Contains(body, "alice") {
		t.Fatalf("expected body %q, got %q", "alice", body)
	}
	if strings.Contains(body, "secret") {
		t.Fatalf("expected body %q, got %q", "secret", body)
	}
	if strings.Contains(body, "password_hash") {
		t.Fatalf("expected body %q, got %q", "password_hash", body)
	}
	req = httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	body = recorder.Body.String()
	if !strings.Contains(body, "alice") {
		t.Fatalf("expected body %q, got %q", "alice", body)
	}
	if strings.Contains(body, "secret") {
		t.Fatalf("expected body %q, got %q", "secret", body)
	}
	if strings.Contains(body, "password_hash") {
		t.Fatalf("expected body %q, got %q", "password_hash", body)
	}
	if !strings.Contains(body, "token") {
		t.Fatalf("expected body %q, got %q", "token", body)
	}
	reqBody = `{"username":"alice","password":"123"}`
	req = httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusUnauthorized {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
	body = recorder.Body.String()
	if !strings.Contains(body, "用户名或密码错误") {
		t.Fatalf("expected body %q, got %q", "用户名或密码错误", body)
	}
}
func TestMeRoute(t *testing.T) {
	todoService := newTestService(nil)
	userService := newTestUserService(nil)
	router := newRouter(todoService, userService, auth.NewJWT("test-secret"), nil)
	req := httptest.NewRequest(http.MethodPost, "/users/register", strings.NewReader(`{"username":"alice","password":"secret"}`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	body := recorder.Body.String()
	if !strings.Contains(body, "alice") {
		t.Fatalf("expected body %q, got %q", "alice", body)
	}
	if strings.Contains(body, "secret") {
		t.Fatalf("expected body %q, got %q", "secret", body)
	}
	if strings.Contains(body, "password_hash") {
		t.Fatalf("expected body %q, got %q", "password_hash", body)
	}
	req = httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(`{"username":"alice","password":"secret"}`))
	req.Header.Set("Content-Type", "application/json")
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	body = recorder.Body.String()
	if !strings.Contains(body, "alice") {
		t.Fatalf("expected body %q, got %q", "alice", body)
	}
	if strings.Contains(body, "secret") {
		t.Fatalf("expected body %q, got %q", "secret", body)
	}
	if strings.Contains(body, "password_hash") {
		t.Fatalf("expected body %q, got %q", "password_hash", body)
	}
	if !strings.Contains(body, "token") {
		t.Fatalf("expected body %q, got %q", "token", body)
	}
	req = httptest.NewRequest(http.MethodGet, "/users/me", nil)
	var loginResp struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	err := json.Unmarshal(recorder.Body.Bytes(), &loginResp)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	token := loginResp.Data.Token
	if token == "" {
		t.Fatalf("expected token, got %q", token)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	body = recorder.Body.String()

	if !strings.Contains(body, "user_id") {
		t.Fatalf("expected body %q, got %q", "user_id", body)
	}
	req = httptest.NewRequest(http.MethodGet, "/users/me", nil)
	// 不要设置 Authorization
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			recorder.Code, http.StatusUnauthorized)
	}
	body = recorder.Body.String()
	if !strings.Contains(body, "未授权") {
		t.Fatalf("expected body %q, got %q", "未授权", body)
	}
}

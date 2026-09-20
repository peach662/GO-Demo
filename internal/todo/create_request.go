package todo

type CreateTodoRequest struct {
	Title string `json:"title" binding:"required"`
}

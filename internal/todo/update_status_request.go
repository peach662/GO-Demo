package todo

type UpdateTodoStatusRequest struct {
	Done *bool `json:"done" binding:"required"`
}

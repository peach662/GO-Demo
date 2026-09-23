package todo

type UpdateTodoStatusRequest struct {
	Status *Status `json:"status" binding:"required"`
}

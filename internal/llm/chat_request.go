package llm
type ChatRequest struct {
	Message string `json:"message" binding:"required"`
}
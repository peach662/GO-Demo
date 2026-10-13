package todo  


type EventTodo struct {

	EventType string `json:"event_type"`
	
	Todo  Todo `json:"todo"`
}

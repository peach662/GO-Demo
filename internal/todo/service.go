package todo

import "sync"

type Service struct {
	mu    sync.RWMutex
	todos []Todo
}

func NewService(initialTodos []Todo) *Service {
	return &Service{
		todos: initialTodos,
	}
}

func (s *Service) List() []Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Todo, len(s.todos))
	copy(result, s.todos)
	return result

}

func (s *Service) GetByID(id int) (Todo, bool) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, item := range s.todos {
		if item.ID == id {
			return item, true
		}
	}
	return Todo{}, false
}

func (s *Service) Create(title string) Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	newTodo := Todo{
		ID:    len(s.todos) + 1,
		Title: title,
		Done:  false,
	}
	s.todos = append(s.todos, newTodo)
	return newTodo
}

func (s *Service) UpdateStatus(id int, done bool) (Todo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, item := range s.todos {
		if item.ID == id {
			s.todos[i].Done = done
			return s.todos[i], true
		}
	}
	return Todo{}, false
}

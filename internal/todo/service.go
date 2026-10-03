package todo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Service struct {
	repo  Repository
	redis *redis.Client
}

func NewService(repo Repository, redis *redis.Client) *Service {
	return &Service{
		repo:  repo,
		redis: redis,
	}
}

func (s *Service) List(ctx context.Context, userID int) ([]Todo, error) {
	var todos []Todo
	key := fmt.Sprintf("todos:%d", userID)
	if s.redis == nil {
		return s.repo.List(ctx, userID)
	}
	cached, err := s.redis.Get(ctx, key).Result()
	if err == nil && cached != "" {
		err = json.Unmarshal([]byte(cached), &todos)
		if err == nil {
			return todos, nil

		}
	}
	todos, err = s.repo.List(ctx, userID)
	if err != nil {
		return todos, err
	}
	data, err := json.Marshal(todos)
	if err != nil {
		return todos, nil
	}

	err = s.redis.Set(ctx, key, data, 2*time.Minute).Err()
	if err != nil {
		return todos, nil
	}
	return todos, nil

}

func (s *Service) GetByID(ctx context.Context, id int, userID int) (Todo, bool, error) {

	if s.redis == nil {
		return s.repo.GetByID(ctx, id, userID)
	}
	var todo Todo
	key := fmt.Sprintf("todo:%d:%d", userID, id)
	cached, err := s.redis.Get(ctx, key).Result()
	if err == nil && cached != "" {
		err = json.Unmarshal([]byte(cached), &todo)
		if err == nil {
			return todo, true, nil
		}

	}

	todo, found, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return Todo{}, false, err
	}
	if !found {
		return Todo{}, false, nil
	}
	data, err := json.Marshal(todo)
	if err != nil {
		return todo, true, nil
	}
	err = s.redis.Set(ctx, key, data, 5*time.Minute).Err()
	if err != nil {
		return todo, true, nil
	}
	return todo, true, nil

}

func (s *Service) Create(ctx context.Context, title string, userID int) (Todo, error) {
	if s.redis == nil {
		return s.repo.Create(ctx, title, userID)
	}

	todo, err := s.repo.Create(ctx, title, userID)
	if err != nil {
		return todo, err
	}
	key := fmt.Sprintf("todos:%d", userID)
	_ = s.redis.Del(ctx, key).Err()
	return todo, nil
}

func (s *Service) UpdateStatus(ctx context.Context, id int, status Status, userID int) (Todo, bool, error) {
	current, found, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return Todo{}, false, err
	}
	if !found {
		return Todo{}, false, nil
	}
	if !CanTransition(current.Status, status) {
		return Todo{}, false, ErrInvalidTransition
	}

	updated, found, err := s.repo.UpdateStatus(ctx, id, status, userID)
	if err != nil {
		return Todo{}, false, err
	}
	if !found {
		return Todo{}, false, nil
	}
	key := fmt.Sprintf("todo:%d:%d", userID, id)
	listKey := fmt.Sprintf("todos:%d", userID)
	if s.redis == nil {
		return updated, true, nil
	}

	_ = s.redis.Del(ctx, key, listKey).Err()
	return updated, true, nil

}

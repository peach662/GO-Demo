package todo

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) List(ctx context.Context, userID int) ([]Todo, error) {

	return s.repo.List(ctx, userID)

}

func (s *Service) GetByID(ctx context.Context, id int, userID int) (Todo, bool, error) {

	return s.repo.GetByID(ctx, id, userID)
}

func (s *Service) Create(ctx context.Context, title string, userID int) (Todo, error) {
	return s.repo.Create(ctx, title, userID)
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
	return s.repo.UpdateStatus(ctx, id, status, userID)
}

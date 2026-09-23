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

func (s *Service) List(ctx context.Context) ([]Todo, error) {

	return s.repo.List(ctx)

}

func (s *Service) GetByID(ctx context.Context, id int) (Todo, bool, error) {

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, title string) (Todo, error) {

	return s.repo.Create(ctx, title)
}

func (s *Service) UpdateStatus(ctx context.Context, id int, status Status) (Todo, bool, error) {
	current, found, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Todo{}, false, err
	}
	if !found {
		return Todo{}, false, nil
	}
	if !CanTransition(current.Status, status) {
		return Todo{}, false, ErrInvalidTransition
	}
	return s.repo.UpdateStatus(ctx, id, status)
}

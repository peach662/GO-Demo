package todo

import "context"

type Repository interface {
	List(ctx context.Context) ([]Todo, error)
	GetByID(ctx context.Context, id int) (Todo, bool, error)
	Create(ctx context.Context, title string) (Todo, error)
	UpdateStatus(ctx context.Context, id int, status Status) (Todo, bool, error)
}

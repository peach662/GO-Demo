package todo

import "context"

type Repository interface {
	List(ctx context.Context,userID int) ([]Todo, error)
	GetByID(ctx context.Context, id int,userID int) (Todo, bool, error)
	Create(ctx context.Context, title string, userID int) (Todo, error)
	UpdateStatus(ctx context.Context, id int, status Status, userID int) (Todo, bool, error)
}

package todo

import (
	"context"
	"database/sql"
	"errors"
)

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{
		db: db,
	}
}
func (r *MySQLRepository) List(ctx context.Context) ([]Todo, error) {
	const query = `
		SELECT id, title, done
		FROM todos
		ORDER BY id
		`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := make([]Todo, 0)
	for rows.Next() {
		var item Todo
		if err := rows.Scan(&item.ID, &item.Title, &item.Done); err != nil {
			return nil, err
		}
		todos = append(todos, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}
func (r *MySQLRepository) GetByID(ctx context.Context, id int) (Todo, bool, error) {
	const query = `SELECT id,title,done FROM todos WHERE id = ?`

	var item Todo
	err := r.db.QueryRowContext(ctx, query, id).Scan(&item.ID, &item.Title, &item.Done)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Todo{}, false, nil
		}
		return Todo{}, false, err
	}

	return item, true, nil

}

func (r *MySQLRepository) Create(
	ctx context.Context,
	title string,
) (Todo, error) {
	const query = `INSERT INTO todos (title, done) VALUES (?, ?)`
	result, err := r.db.ExecContext(ctx, query, title, false)
	if err != nil {
		return Todo{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Todo{}, err
	}
	return Todo{ID: int(id), Title: title, Done: false}, nil
}
func (r *MySQLRepository) UpdateStatus(
	ctx context.Context,
	id int,
	done bool,
) (Todo, bool, error) {
	const query = `UPDATE todos SET done = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, done, id)
	if err != nil {
		return Todo{}, false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Todo{}, false, err
	}
	if rowsAffected == 0 {
		return Todo{}, false, nil
	}
	item, found, err := r.GetByID(ctx, id)
	if err != nil {
		return Todo{}, false, err
	}
	return item, found, nil
}

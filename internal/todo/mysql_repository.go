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
		SELECT id, title, status
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
		if err := rows.Scan(&item.ID, &item.Title, &item.Status); err != nil {
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
	const query = `SELECT id,title,status FROM todos WHERE id = ?`

	var item Todo
	err := r.db.QueryRowContext(ctx, query, id).Scan(&item.ID, &item.Title, &item.Status)
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
	const query = `INSERT INTO todos (title, status) VALUES (?, ?)`
	result, err := r.db.ExecContext(ctx, query, title, StatusPending)
	if err != nil {
		return Todo{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Todo{}, err
	}
	return Todo{
		ID:     int(id),
		Title:  title,
		Status: StatusPending,
	}, nil
}

func (r *MySQLRepository) UpdateStatus(
	ctx context.Context,
	id int,
	status Status,
) (Todo, bool, error) {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Todo{}, false, err
	}
	defer tx.Rollback()

	const selectQuery = `SELECT id,title,status FROM todos WHERE id = ? FOR UPDATE`
	var item Todo
	err = tx.QueryRowContext(ctx, selectQuery, id).Scan(&item.ID, &item.Title, &item.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Todo{}, false, nil
		}
		return Todo{}, false, err
	}
	if item.Status == status {
		return item, true, nil
	}
	const updateQuery = `UPDATE todos SET status = ? WHERE id = ?`
	if _, err = tx.ExecContext(ctx, updateQuery, status, id); err != nil {
		return Todo{}, false, err
	}

	const insertLogQuery = `INSERT INTO todo_status_logs (todo_id, from_status, to_status) VALUES (?, ?, ?)`
	if _, err = tx.ExecContext(ctx, insertLogQuery, id, item.Status, status); err != nil {
		return Todo{}, false, err
	}
	if err = tx.Commit(); err != nil {
		return Todo{}, false, err
	}
	item.Status = status

	return item, true, nil
}

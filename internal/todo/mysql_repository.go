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
func scanTodo(scanner interface {
	Scan(dest ...any) error
}) (Todo, error) {
	var item Todo
	var userID sql.NullInt64
	if err := scanner.Scan(&item.ID, &item.Title, &item.Status, &userID); err != nil {
		return Todo{}, err
	}
	if userID.Valid {
		item.UserID = int(userID.Int64)
	}
	return item, nil
}

func (r *MySQLRepository) List(ctx context.Context) ([]Todo, error) {
	const query = `
		SELECT id, title, status, user_id
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
		item, err := scanTodo(rows)
		if err != nil {
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
	const query = `SELECT id,title,status,user_id FROM todos WHERE id = ?`

	item, err := scanTodo(r.db.QueryRowContext(ctx, query, id))
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
	userID int,
) (Todo, error) {
	const query = `INSERT INTO todos (title, status, user_id) VALUES (?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, title, StatusPending, userID)
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
		UserID: userID,
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

	const selectQuery = `SELECT id,title,status,user_id FROM todos WHERE id = ? FOR UPDATE`
	item, err := scanTodo(tx.QueryRowContext(ctx, selectQuery, id))
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

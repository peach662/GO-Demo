package todo

import (
	"context"
	"database/sql"
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

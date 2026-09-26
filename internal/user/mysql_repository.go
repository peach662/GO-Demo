package user

import (
	"context"
	"database/sql"
	"errors"
)

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}
func (r *MySQLRepository) GetByUsername(ctx context.Context, username string) (User, bool, error) {
	const query = `
		SELECT id, username, password_hash
		FROM users
		WHERE username = ?
	`
	var user User
	err := r.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, false, nil
		}
		return User{}, false, err
	}
	return user, true, nil
}

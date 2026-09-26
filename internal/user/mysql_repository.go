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

func (r *MySQLRepository) Create(ctx context.Context, username, passwordHash string) (User, error) {
	const query = `
		INSERT INTO users (username, password_hash)
		VALUES (?, ?)
	`
	result, err := r.db.ExecContext(ctx, query, username, passwordHash)
	if err != nil {
		return User{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return User{}, err
	}
	return User{ID: int(id), Username: username, PasswordHash: passwordHash}, nil
}

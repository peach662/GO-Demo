package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const mysqlDSN = "app:app123@tcp(127.0.0.1:13306)/awesome_project?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"

func OpenMySQL() (*sql.DB, error) {
	db, err := sql.Open("mysql", mysqlDSN)
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping mysql: %w", err)
	}

	return db, nil
}

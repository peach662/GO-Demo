package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MySQLDSN string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}

	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		return Config{}, fmt.Errorf("MYSQL_DSN is not set")
	}

	return Config{
		MySQLDSN: dsn,
	}, nil
}

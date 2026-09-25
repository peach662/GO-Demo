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
	if err := loadLocalEnv(".env", ".env.example"); err != nil {
		return Config{}, err
	}

	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		return Config{}, fmt.Errorf("MYSQL_DSN is not set")
	}

	return Config{
		MySQLDSN: dsn,
	}, nil
}

func loadLocalEnv(paths ...string) error {
	for _, path := range paths {
		err := godotenv.Load(path)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("load %s: %w", path, err)
		}
		if os.Getenv("MYSQL_DSN") != "" {
			return nil
		}
	}
	return nil
}

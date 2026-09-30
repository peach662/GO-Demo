package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MySQLDSN  string
	JWTSecret string
	RedisAddr string
}

func Load() (Config, error) {
	if err := loadLocalEnv(".env", ".env.example"); err != nil {
		return Config{}, err
	}

	dsn := os.Getenv("MYSQL_DSN")
	jwtSecret := os.Getenv("JWT_SECRET")
	redisAddr := os.Getenv("REDIS_ADDR")
	if dsn == "" {
		return Config{}, fmt.Errorf("MYSQL_DSN is not set")
	}
	if jwtSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is not set")
	}
	if redisAddr == "" {
		return Config{}, fmt.Errorf("REDIS_ADDR is not set")
	}
	return Config{
		MySQLDSN:  dsn,
		JWTSecret: jwtSecret,
		RedisAddr: redisAddr,
	}, nil
}

func loadLocalEnv(paths ...string) error {
	for _, path := range paths {
		err := godotenv.Load(path)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("load %s: %w", path, err)
		}
		if os.Getenv("MYSQL_DSN") != "" && os.Getenv("JWT_SECRET") != "" && os.Getenv("REDIS_ADDR") != "" {
			return nil
		}
	}
	return nil
}

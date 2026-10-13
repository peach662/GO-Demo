package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MySQLDSN      string
	JWTSecret     string
	RedisAddr     string
	RabbitMQURL   string
	OpenAIAPIKey  string
	OpenAIBaseURL string
	OpenAIModel   string
}

func Load() (Config, error) {
	if err := loadLocalEnv(".env", ".env.example"); err != nil {
		return Config{}, err
	}

	dsn := os.Getenv("MYSQL_DSN")
	jwtSecret := os.Getenv("JWT_SECRET")
	redisAddr := os.Getenv("REDIS_ADDR")
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	openAIBaseURL := os.Getenv("OPENAI_BASE_URL")
	openAIModel := os.Getenv("OPENAI_MODEL")
	if dsn == "" {
		return Config{}, fmt.Errorf("MYSQL_DSN is not set")
	}
	if jwtSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is not set")
	}
	if redisAddr == "" {
		return Config{}, fmt.Errorf("REDIS_ADDR is not set")
	}
	if rabbitMQURL == "" {
		return Config{}, fmt.Errorf("RABBITMQ_URL is not set")
	}
	if openAIAPIKey == "" {
		return Config{}, fmt.Errorf("OPENAI_API_KEY is not set")
	}
	if openAIBaseURL == "" {
		return Config{}, fmt.Errorf("OPENAI_BASE_URL is not set")
	}
	if openAIModel == "" {
		return Config{}, fmt.Errorf("OPENAI_MODEL is not set")
	}

	return Config{
		MySQLDSN:      dsn,
		JWTSecret:     jwtSecret,
		RedisAddr:     redisAddr,
		RabbitMQURL:   rabbitMQURL,
		OpenAIAPIKey:  openAIAPIKey,
		OpenAIBaseURL: openAIBaseURL,
		OpenAIModel:   openAIModel,
	}, nil
}

func loadLocalEnv(paths ...string) error {
	for _, path := range paths {
		err := godotenv.Load(path)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("load %s: %w", path, err)
		}
		if os.Getenv("MYSQL_DSN") != "" &&
			os.Getenv("JWT_SECRET") != "" &&
			os.Getenv("REDIS_ADDR") != "" &&
			os.Getenv("RABBITMQ_URL") != "" &&
			os.Getenv("OPENAI_API_KEY") != "" &&
			os.Getenv("OPENAI_BASE_URL") != "" &&
			os.Getenv("OPENAI_MODEL") != "" {
			return nil
		}
	}
	return nil
}

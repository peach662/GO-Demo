package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
	"fmt"
)

type JWT struct {
	SecretKey string
}

func NewJWT(secretKey string) *JWT {
	return &JWT{SecretKey: secretKey}
}

func (j *JWT) GenerateToken(userID int) (string, error) {
	if j.SecretKey == "" {
		return "", fmt.Errorf("jwt secret is empty")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(j.SecretKey))
}
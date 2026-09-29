package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
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
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(j.SecretKey))
}

func (j *JWT) ParseToken(tokenString string) (userID int, err error) {
	if j.SecretKey == "" {
		return 0, fmt.Errorf("jwt secret is empty")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return 0, fmt.Errorf("parse token failed: %w", err)
	}
	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}
	idFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid user id type")
	}

	return int(idFloat), nil
}

package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthService struct {
	secret string
}

func NewAuthService(secret string) *AuthService {
	return &AuthService{
		secret: secret,
	}
}

func (s *AuthService) Sign(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	return token.SignedString([]byte(s.secret))
}

func (s *AuthService) Verify(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.secret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if int64(claims["exp"].(float64)) < time.Now().Unix() {
			return "", errors.New("token expired")
		}

		return claims["userId"].(string), nil
	}

	return "", errors.New("invalid token")
}

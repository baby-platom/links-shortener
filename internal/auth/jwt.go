package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/baby-platom/links-shortener/internal/config"
	"github.com/baby-platom/links-shortener/internal/logger"
	"github.com/golang-jwt/jwt/v4"
)

const UserID int = 1

type claims struct {
	jwt.RegisteredClaims
	UserID int
}

// BuildJWTString - creates and return a string representation of JWT token
func BuildJWTString() (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Config.AuthTTL)),
			},
			UserID: UserID,
		},
	)

	tokenString, err := token.SignedString([]byte(config.Config.AuthSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserId(tokenString string) (int, error) {
	claims := &claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(config.Config.AuthSecretKey), nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("Error occured while parsing JWT token: %w", err)
	}

	if !token.Valid {
		logger.Log.Warn("Token is not valid")
		return 0, errors.New("Token is not valid")
	}

	logger.Log.Info("Token is valid")
	return claims.UserID, nil
}

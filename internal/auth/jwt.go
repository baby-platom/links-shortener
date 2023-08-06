package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/baby-platom/links-shortener/internal/config"
	"github.com/baby-platom/links-shortener/internal/logger"
	"github.com/baby-platom/links-shortener/internal/shortid"
	"github.com/golang-jwt/jwt/v4"
)

type claims struct {
	jwt.RegisteredClaims
	UserID string
}

// BuildJWTString - creates and return a string representation of JWT token
func BuildJWTString() (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims{
			RegisteredClaims: jwt.RegisteredClaims{},
			UserID:           shortid.GenerateShortID(),
		},
	)

	tokenString, err := token.SignedString([]byte(config.Config.AuthSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserID(tokenString string) (string, error) {
	claims := &claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(config.Config.AuthSecretKey), nil
		},
	)

	if err != nil {
		return "", fmt.Errorf("error occured while parsing JWT token: %w", err)
	}

	if !token.Valid {
		logger.Log.Warn("Token is not valid")
		return "", errors.New("token is not valid")
	}

	logger.Log.Info("Token is valid")
	return claims.UserID, nil
}

func GetUserIDForHandler(r *http.Request, setCookieHeader string) string {
	var authToken string
	authCookie, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		data := strings.Split(setCookieHeader, ";")[0]
		authToken = strings.Split(data, "=")[1]
	} else {
		authToken = authCookie.Value
	}
	userID, _ := GetUserID(authToken)
	return userID
}

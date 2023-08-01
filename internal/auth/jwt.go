package auth

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/baby-platom/links-shortener/internal/config"
	"github.com/baby-platom/links-shortener/internal/logger"
	"github.com/golang-jwt/jwt/v4"
)

type claims struct {
	jwt.RegisteredClaims
	UserID int
}

var r = rand.New(rand.NewSource(time.Now().Unix()))
var limit = int(math.Pow(2, 20))

// BuildJWTString - creates and return a string representation of JWT token
func BuildJWTString() (string, error) {
	UserID := r.Intn(limit)

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

func GetUserID(tokenString string) (int, error) {
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
		return 0, fmt.Errorf("error occured while parsing JWT token: %w", err)
	}

	if !token.Valid {
		logger.Log.Warn("Token is not valid")
		return 0, errors.New("token is not valid")
	}

	logger.Log.Info("Token is valid")
	return claims.UserID, nil
}

func GetUserIDForHandler(w http.ResponseWriter, r *http.Request) int {
	var authToken string
	authCookie, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		data := w.Header().Get("Set-Cookie")
		data = strings.Split(data, ";")[0]
		authToken = strings.Split(data, "=")[1]
	} else {
		authToken = authCookie.Value
	}
	userID, _ := GetUserID(authToken)
	return userID
}

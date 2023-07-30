package auth

import (
	"net/http"

	"github.com/baby-platom/links-shortener/internal/config"
	"github.com/baby-platom/links-shortener/internal/logger"
)

// Middleware for compression and decompression
func Middleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var passNewAuthCookie bool
		authCookie, err := r.Cookie("auth")
		if err != nil {
			passNewAuthCookie = true
			if err == http.ErrNoCookie {
				logger.Log.Info("No auth cookie passed")
			} else {
				logger.Log.Errorw(
					"Unexpected error occured while getting auth cookie",
					"error", err,
				)
			}
		}

		if authCookie != nil {
			_, err = GetUserId(authCookie.Value)
			if err != nil {
				passNewAuthCookie = true
				logger.Log.Warnw(
					"Error occured while parsing auth cookie",
					"error", err,
				)
			}
		}

		if passNewAuthCookie {
			logger.Log.Info("Creating new auth cookie")

			newAuthToken, err := BuildJWTString()
			if err == nil {
				cookie := &http.Cookie{
					Name:   "auth",
					Value:  newAuthToken,
					MaxAge: int(config.Config.AuthTTL.Seconds()),
					Path:   "/",
				}
				http.SetCookie(w, cookie)
			} else {
				logger.Log.Errorw(
					"Error occured while building new JWT auth token",
					"error", err,
				)
			}

			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

package app

import (
	"net/http"

	"github.com/baby-platom/links-shortener/internal/config"
	"github.com/baby-platom/links-shortener/internal/logger"
	"github.com/go-chi/chi/v5"
)

// ShortenedUrlsByIDType map id:initalUrl
type ShortenedUrlsByIDType map[string]string

// ShortenedUrlsByID stores initial urls
var ShortenedUrlsByID = make(ShortenedUrlsByIDType)

// Run server
func Run() error {
	if err := logger.Initialize(config.Config.LogLevel); err != nil {
		return err
	}
	return http.ListenAndServe(config.Config.Address, logger.Middleware(Router()))
}

// Router prepares and returns chi.Router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Post("/api/shorten", shortenAPIHandler)

	r.Post("/", shortenURLHandler)
	r.Get("/{id}", restoreURLHandler)
	return r
}

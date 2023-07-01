package app

import (
	"net/http"

	"github.com/baby-platom/links-shortener/internal/config"
	"github.com/baby-platom/links-shortener/internal/logger"
	"github.com/go-chi/chi/v5"
)

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

	r.Post("/", shortenURLHandler)
	r.Get("/{id}", restoreURLHandler)
	return r
}

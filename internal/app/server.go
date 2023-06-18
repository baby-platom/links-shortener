package app

import (
	"net/http"

	"github.com/baby-platom/links-shortener/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Run server
func Run() error {
	return http.ListenAndServe(config.Config.Address, Router())
}

// Router prepares and returns chi.Router
func Router() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/", shortenURLHandler)
	r.Get("/{id}", restoreURLHandler)
	return r
}

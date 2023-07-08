package app

import (
	"net/http"

	"github.com/baby-platom/links-shortener/internal/compress"
	"github.com/baby-platom/links-shortener/internal/config"
	"github.com/baby-platom/links-shortener/internal/logger"
	"github.com/baby-platom/links-shortener/internal/shortid"
	"github.com/go-chi/chi/v5"
)

// ShortenedUrlsByID stores initial urls
var ShortenedUrlsByID = make(shortid.ShortenedUrlsByIDType)

// Run server
func Run() error {
	if err := logger.Initialize(config.Config.LogLevel); err != nil {
		return err
	}
	ShortenedUrlsByID.Load(config.Config.FileStoragePath)
	return http.ListenAndServe(
		config.Config.Address,
		Router(),
	)
}

// Router prepares and returns chi.Router
func Router() chi.Router {
	r := chi.NewRouter()
	r.Use(logger.Middleware)
	r.Use(compress.Middleware)

	r.Post("/api/shorten", shortenAPIHandler)

	r.Post("/", shortenURLHandler)
	r.Get("/{id}", restoreURLHandler)
	return r
}

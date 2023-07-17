package app

import (
	"net/http"
	"context"

	"github.com/baby-platom/links-shortener/internal/compress"
	"github.com/baby-platom/links-shortener/internal/config"
	"github.com/baby-platom/links-shortener/internal/database"
	"github.com/baby-platom/links-shortener/internal/logger"
	"github.com/baby-platom/links-shortener/internal/shortid"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// ShortenedUrlsByID stores initial urls
var ShortenedUrlsByID shortid.ShortenedUrlsByIDInterface

// Run server
func Run() error {
	if err := logger.Initialize(config.Config.LogLevel); err != nil {
		return err
	}

	switch {
	case config.Config.DatabaseDSN != "":
		ShortenedUrlsByID = shortid.NewShortenedUrlsByIDDatabase()
	case config.Config.FileStoragePath != "":
		ShortenedUrlsByID = shortid.NewShortenedUrlsByIDJson(config.Config.FileStoragePath)
	default:
		ShortenedUrlsByID = shortid.NewShortenedUrlsByID()
	}

	switch ShortenedUrlsByID.(type) {
	case *shortid.ShortenedUrlsByIDDatabaseType:
		err := database.OpenPostgres(config.Config.DatabaseDSN)
		if err != nil {
			panic(err)
		}
		defer database.Connection.Close()

		ctx := context.Background()
		exists, err := database.Connection.CheckIfShortIDsTableExists(ctx)
		if err != nil {
			panic(err)
		}

		if !exists {
			err = database.Connection.CreateShortIDsTable(ctx)
		}
		if err != nil {
			panic(err)
		}
	}

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
	r.Use(middleware.Compress(5, compress.ContentTypesToBeEncoded...))

	r.Post("/api/shorten", shortenAPIHandler)
	r.Get("/ping", pingDatabaseAPIHandler)

	r.Post("/", shortenURLHandler)
	r.Get("/{id}", restoreURLHandler)
	return r
}

package app

import (
	"context"
	"net/http"

	"github.com/baby-platom/links-shortener/internal/auth"
	"github.com/baby-platom/links-shortener/internal/compress"
	"github.com/baby-platom/links-shortener/internal/config"
	"github.com/baby-platom/links-shortener/internal/database"
	"github.com/baby-platom/links-shortener/internal/logger"
	"github.com/baby-platom/links-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// ShortenedUrlsByIDStorage  stores and provides access to shortened urls by ID.
var ShortenedUrlsByIDStorage storage.ShortenedUrlsByIDStorer

// Run server
func Run() error {
	if err := logger.Initialize(config.Config.LogLevel); err != nil {
		return err
	}

	switch {
	case config.Config.DatabaseDSN != "":
		logger.Log.Info("Use NewShortenedUrlsByIDDatabase")
		ShortenedUrlsByIDStorage = storage.CreateNewShortenedUrlsByIDDBStorer()
	case config.Config.FileStoragePath != "":
		logger.Log.Info("Use NewShortenedUrlsByIDJson")
		ShortenedUrlsByIDStorage = storage.CreateNewShortenedUrlsByIDFileStorer(config.Config.FileStoragePath)
	default:
		logger.Log.Info("Use NewShortenedUrlsByID")
		ShortenedUrlsByIDStorage = storage.CreateNewShortenedUrlsByIDMemoryStorer()
	}

	switch ShortenedUrlsByIDStorage.(type) {
	case *storage.ShortenedUrlsByIDDBStorer:
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
		logger.Log.Infof("ShortIDs table exists: %v", exists)

		if !exists {
			logger.Log.Info("Creating new ShortIDs table")
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
	r.Use(auth.Middleware)
	r.Use(middleware.Compress(5, compress.ContentTypesToBeEncoded...))

	r.Get("/api/user/urls", getUserShortenURLsAPIHandler)
	r.Post("/api/shorten/batch", shortenBatchAPIHandler)
	r.Post("/api/shorten", shortenAPIHandler)

	r.Get("/ping", pingDatabaseAPIHandler)
	r.Post("/", shortenURLHandler)
	r.Get("/{id}", restoreURLHandler)
	return r
}

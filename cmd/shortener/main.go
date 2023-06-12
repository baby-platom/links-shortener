package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/baby-platom/links-shortener/cmd/config"
	"github.com/baby-platom/links-shortener/internal/shortid"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type customHandler struct{}

var shortenedUrlsByID = make(map[string]string)

func main() {
	config.ParseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	return http.ListenAndServe(config.Config.Address, Router())
}

// Router prepares and returns chi.Router
func Router() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/", ShortenURLHandler)
	r.Get("/{id}", RestoreURLHandler)
	return r
}

// ShortenURLHandler returns a shortened version of a passed URL
func ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := shortid.GenerateShortID()
	shortenedUrlsByID[id] = string(body)
	fmt.Printf("test/n%s/n/ntest", string(body))

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s/%s", config.Config.BaseAddress, id)
}

// RestoreURLHandler restore a URL if it before was shortened
func RestoreURLHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	url, ok := shortenedUrlsByID[id]
	if !ok {
		http.Error(w, "Nonexistent Id", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

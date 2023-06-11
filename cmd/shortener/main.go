package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/baby-platom/links-shortener/internal/shortid"
)

type customHandler struct{}

var shortenedUrlsByID = make(map[string]string)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var h customHandler

	return http.ListenAndServe(`:8080`, h)
}

func (h customHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		ShortenURLHandler(w, r)
	case http.MethodGet:
		RestoreURLHandler(w, r)
	default:
		http.Error(w, "Only POST and GET methods are allowed", http.StatusBadRequest)
		return
	}
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

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "http://%s/%s", r.Host, id)
}

// RestoreURLHandler restore a URL if it before was shortened
func RestoreURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[1:]
	url, ok := shortenedUrlsByID[id]
	if !ok {
		http.Error(w, "Nonexistent Id", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

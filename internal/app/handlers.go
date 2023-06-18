package app

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/baby-platom/links-shortener/internal/config"
	"github.com/baby-platom/links-shortener/internal/shortid"
	"github.com/go-chi/chi/v5"
)

// ShortenedUrlsByIDType map id:initalUrl
type ShortenedUrlsByIDType map[string]string

// ShortenedUrlsByID stores initial urls
var ShortenedUrlsByID = make(ShortenedUrlsByIDType)

func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	bodyString := string(body)
	if strings.TrimSpace(bodyString) == "" {
		http.Error(w, "Body is empty", http.StatusBadRequest)
		return
	}

	id := shortid.GenerateShortID()
	ShortenedUrlsByID[id] = bodyString
	fmt.Printf("Shortened '%s' to '%s'", bodyString, id)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s/%s", config.Config.BaseAddress, id)
}

func restoreURLHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	url, ok := ShortenedUrlsByID[id]
	if !ok {
		http.Error(w, "Nonexistent Id", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

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
	ShortenedUrlsByID.Save(r.Context(), id, bodyString)
	fmt.Printf("Shortened '%s' to '%s'\n", bodyString, id)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s/%s", config.Config.BaseAddress, id)
}

func restoreURLHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	url, ok := ShortenedUrlsByID.Get(r.Context(), id)
	if !ok {
		http.Error(w, "Nonexistent Id", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

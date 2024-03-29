package app

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/baby-platom/links-shortener/internal/auth"
	"github.com/baby-platom/links-shortener/internal/config"
	"github.com/baby-platom/links-shortener/internal/database"
	"github.com/baby-platom/links-shortener/internal/logger"
	"github.com/baby-platom/links-shortener/internal/shortid"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	initialURL := string(body)
	if strings.TrimSpace(initialURL) == "" {
		http.Error(w, "Body is empty", http.StatusBadRequest)
		return
	}

	userID := auth.GetUserIDForHandler(r, w.Header().Get("Set-Cookie"))

	id := shortid.GenerateShortID()
	err = ShortenedUrlsByIDStorage.Save(r.Context(), id, initialURL, userID)
	if err != nil && errors.Is(err, database.ErrConflict) {
		logger.Log.Error("Cannot shorten url", zap.Error(err))
		id, err := ShortenedUrlsByIDStorage.GetIDByURL(r.Context(), initialURL)
		if err != nil {
			logger.Log.Error("Cannot get already shortened url", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "%s/%s", config.Config.BaseAddress, id)
		return
	}
	fmt.Printf("Shortened '%s' to '%s'\n", initialURL, id)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s/%s", config.Config.BaseAddress, id)
}

func restoreURLHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	url, ok, deleted := ShortenedUrlsByIDStorage.Get(r.Context(), id)

	if !ok {
		http.Error(w, "Nonexistent Id", http.StatusBadRequest)
		return
	}

	if deleted {
		w.WriteHeader(http.StatusGone)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

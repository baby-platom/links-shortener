package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/baby-platom/links-shortener/internal/config"
	"github.com/baby-platom/links-shortener/internal/database"
	"github.com/baby-platom/links-shortener/internal/logger"
	"github.com/baby-platom/links-shortener/internal/models"
	"github.com/baby-platom/links-shortener/internal/shortid"
	"go.uber.org/zap"
)

func shortenAPIHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ShortenRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Error("Cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "No 'url' passed in request", http.StatusBadRequest)
		return
	}

	id := shortid.GenerateShortID()
	ShortenedUrlsByID.Save(r.Context(), id, req.URL)
	logger.Log.Infof("Shortened '%s' to '%s'\n", req.URL, id)

	resp := models.ShortenResponse{
		Result: fmt.Sprintf("%s/%s", config.Config.BaseAddress, id),
	}
	data, err := json.Marshal(resp)
	if err != nil {
		logger.Log.Error("Error encoding response", zap.Error(err))
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func pingDatabaseAPIHandler(w http.ResponseWriter, r *http.Request) {
	result := database.Connection.HealthCheck(r.Context())
	if !result {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func shortenBatchAPIHandler(w http.ResponseWriter, r *http.Request) {
	var req []models.BatchPortionShortenRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Error("Cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(req) == 0 {
		http.Error(w, "No data passed in request", http.StatusBadRequest)
		return
	}

	var shortenedUrlsByIds []models.BatchPortionShortenResponse
	for _, portion := range req {
		id := shortid.GenerateShortID()
		b := models.BatchPortionShortenResponse{
			CorrelationID: portion.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", config.Config.BaseAddress, id),
			ID:            id,
			OriginalURL:   portion.OriginalURL,
		}
		shortenedUrlsByIds = append(shortenedUrlsByIds, b)
	}

	err := ShortenedUrlsByID.BatchSave(r.Context(), shortenedUrlsByIds)
	if err != nil {
		logger.Log.Error("Error saving shortened urls", zap.Error(err))
		http.Error(w, "Error saving shortened urls", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(shortenedUrlsByIds)
	if err != nil {
		logger.Log.Error("Error encoding response", zap.Error(err))
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

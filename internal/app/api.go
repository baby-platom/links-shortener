package app

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"reflect"

	"github.com/baby-platom/links-shortener/internal/config"
	"github.com/baby-platom/links-shortener/internal/logger"
	"github.com/baby-platom/links-shortener/internal/models"
	"github.com/baby-platom/links-shortener/internal/shortid"
)

func shortenAPIHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ShortentRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Error("Cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if reflect.ValueOf(req.URL).IsZero() {
		http.Error(w, "No 'url' passed in request", http.StatusBadRequest)
		return
	}

	id := shortid.GenerateShortID()
	ShortenedUrlsByID[id] = req.URL
	ShortenedUrlsByID.Save(config.Config.FileStoragePath)
	logger.Log.Infof("Shortened '%s' to '%s'\n", req.URL, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := models.ShortentResponse{
		Result: fmt.Sprintf("%s/%s", config.Config.BaseAddress, id),
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Error("Error encoding response", zap.Error(err))
		return
	}
}

package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/baby-platom/links-shortener/internal/auth"
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

	userID := auth.GetUserIDForHandler(r, w.Header().Get("Set-Cookie"))

	var id = shortid.GenerateShortID()
	var status = http.StatusCreated
	err := ShortenedUrlsByIDStorage.Save(r.Context(), id, req.URL, userID)
	if err != nil && errors.Is(err, database.ErrConflict) {
		logger.Log.Error("Cannot shorten url", zap.Error(err))
		id, err = ShortenedUrlsByIDStorage.GetIDByURL(r.Context(), req.URL)
		if err != nil {
			logger.Log.Error("Cannot get already shortened url", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		status = http.StatusConflict
	}

	if status == http.StatusCreated {
		logger.Log.Infof("Shortened '%s' to '%s'", req.URL, id)
	}

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
	w.WriteHeader(status)
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

	userID := auth.GetUserIDForHandler(r, w.Header().Get("Set-Cookie"))

	err := ShortenedUrlsByIDStorage.BatchSave(r.Context(), shortenedUrlsByIds, userID)
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

func getUserShortenURLsAPIHandler(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDForHandler(r, w.Header().Get("Set-Cookie"))

	userShortenURLsList, err := ShortenedUrlsByIDStorage.GetUserShortenURLsListResponse(
		r.Context(), config.Config.BaseAddress, userID,
	)
	if err != nil {
		logger.Log.Error("Error occured while getting userShortenURLsList", zap.Error(err))
		http.Error(w, "Error occured while getting userShortenURLsList", http.StatusInternalServerError)
		return
	}

	var noToken bool
	if w.Header().Get("No-Token-Passed") != "" {
		noToken = true
	}

	if noToken {
		logger.Log.Warn("Passed JWT token is invalid")
		w.WriteHeader(http.StatusUnauthorized)
	} else if len(userShortenURLsList) > 0 {
		data, err := json.Marshal(userShortenURLsList)
		if err != nil {
			logger.Log.Error("Error encoding response", zap.Error(err))
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	} else {
		logger.Log.Warnf("User '%s' has no shortened URLs", userID)
		w.WriteHeader(http.StatusNoContent)
	}
}

func deleteUsersShortenedURLsAPIHandler(w http.ResponseWriter, r *http.Request) {
	var req models.BatchDeleteRequest

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

	userID := auth.GetUserIDForHandler(r, w.Header().Get("Set-Cookie"))
	userShortenURLsList, err := ShortenedUrlsByIDStorage.GetUserShortenURLsList(r.Context(), userID)
	if err != nil {
		logger.Log.Error("Error occured while getting userShortenURLsList", zap.Error(err))
		http.Error(w, "Error occured while getting userShortenURLsList", http.StatusInternalServerError)
		return
	}

	logger.Log.Info(req)
	logger.Log.Info(userShortenURLsList)

	checkset := make(map[string]bool)
	for _, shortURL := range req {
		checkset[shortURL] = true
	}
	for _, userShortenURL := range userShortenURLsList {
		if checkset[userShortenURL] {
			delete(checkset, userShortenURL)
		}
	}

	if len(checkset) != 0 {
		logger.Log.Error("Not all pased ShortenedURLs belongs to user")
		http.Error(w, "Not all pased ShortenedURLs belongs to user", http.StatusForbidden)
		return
	}

	ShortenedUrlsByIDStorage.Delete(r.Context(), req, userID)
	w.WriteHeader(http.StatusAccepted)
}

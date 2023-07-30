package storage

import (
	"context"
	"encoding/json"
	"os"

	"github.com/baby-platom/links-shortener/internal/config"
	"github.com/baby-platom/links-shortener/internal/models"
)

// ShortenedUrlsByIDFileStorer implements ShortenedUrlsByIDStorer interface using json file.
type ShortenedUrlsByIDFileStorer struct {
	ShortenedUrlsByIDMemoryStorer
}

// CreateNewShortenedUrlsByIDJson return a new ShortenedUrlsByIDJson
func CreateNewShortenedUrlsByIDFileStorer(fname string) *ShortenedUrlsByIDFileStorer {
	NewShortenedUrlsByIDJSON := &ShortenedUrlsByIDFileStorer{*CreateNewShortenedUrlsByIDMemoryStorer()}
	err := NewShortenedUrlsByIDJSON.LoadJSON(fname)
	if err != nil {
		panic(err)
	}
	return NewShortenedUrlsByIDJSON
}

// Save data into a json file
func (s *ShortenedUrlsByIDFileStorer) SaveJSON(fname string) error {
	res, err := json.MarshalIndent(s.Data, "", "   ")
	if err != nil {
		return err
	}
	return os.WriteFile(fname, res, 0666)
}

// Load data from a json file
func (s *ShortenedUrlsByIDFileStorer) LoadJSON(fname string) error {
	res, err := os.ReadFile(fname)
	switch {
	case os.IsNotExist(err):
		return nil
	case err != nil:
		return err
	}

	if err := json.Unmarshal(res, &s.Data); err != nil {
		return err
	}
	return nil
}

// Save creates new id:url relation and saves it to the json file
func (s *ShortenedUrlsByIDFileStorer) Save(ctx context.Context, id string, url string, userID int) error {
	s.ShortenedUrlsByIDMemoryStorer.Save(ctx, id, url, userID)
	return s.SaveJSON(config.Config.FileStoragePath)
}

func (s *ShortenedUrlsByIDFileStorer) BatchSave(ctx context.Context, shortenedUrlsByIds []models.BatchPortionShortenResponse, userID int) error {
	s.ShortenedUrlsByIDMemoryStorer.BatchSave(ctx, shortenedUrlsByIds, userID)
	if err := s.SaveJSON(config.Config.FileStoragePath); err != nil {
		return err
	}
	return nil
}

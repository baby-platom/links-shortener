package shortid

import (
	"context"
	"encoding/json"
	"os"

	"github.com/baby-platom/links-shortener/internal/config"
	"github.com/baby-platom/links-shortener/internal/models"
)

type ShortenedUrlsByIDJSONType struct {
	ShortenedUrlsByIDType
}

// NewShortenedUrlsByIDJson return a new ShortenedUrlsByIDJson
func NewShortenedUrlsByIDJson(fname string) *ShortenedUrlsByIDJSONType {
	ShortenedUrlsByIDJSON := &ShortenedUrlsByIDJSONType{*NewShortenedUrlsByID()}
	err := ShortenedUrlsByIDJSON.LoadJSON(fname)
	if err != nil {
		panic(err)
	}
	return ShortenedUrlsByIDJSON
}

// Save data into a json file
func (s *ShortenedUrlsByIDJSONType) SaveJSON(fname string) error {
	res, err := json.MarshalIndent(s.Data, "", "   ")
	if err != nil {
		return err
	}
	return os.WriteFile(fname, res, 0666)
}

// Load data from a json file
func (s *ShortenedUrlsByIDJSONType) LoadJSON(fname string) error {
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
func (s *ShortenedUrlsByIDJSONType) Save(ctx context.Context, id string, url string) error {
	s.Data[id] = url
	return s.SaveJSON(config.Config.FileStoragePath)
}

func (s *ShortenedUrlsByIDJSONType) BatchSave(ctx context.Context, shortenedUrlsByIds []models.BatchPortionShortenResponse) error {
	for _, portion := range shortenedUrlsByIds {
		s.Data[portion.ID] = portion.OriginalURL
	}
	if err := s.SaveJSON(config.Config.FileStoragePath); err != nil {
		return err
	}
	return nil
}

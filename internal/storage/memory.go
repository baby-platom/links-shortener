package storage

import (
	"context"

	"github.com/baby-platom/links-shortener/internal/models"
)

// ShortenedUrlsByIDMemoryStorer implements ShortenedUrlsByIDStorer interface using temproray memory.
type ShortenedUrlsByIDMemoryStorer struct {
	Data map[string]string
}

// CreateNewShortenedUrlsByID return a new ShortenedUrlsByID
func CreateNewShortenedUrlsByIDMemoryStorer() *ShortenedUrlsByIDMemoryStorer {
	return &ShortenedUrlsByIDMemoryStorer{Data: make(map[string]string)}
}

// Save creates new id:url relation
func (s *ShortenedUrlsByIDMemoryStorer) Save(ctx context.Context, id string, url string) error {
	s.Data[id] = url
	return nil
}

// Get returns url by id
func (s *ShortenedUrlsByIDMemoryStorer) Get(ctx context.Context, id string) (string, bool) {
	url, ok := s.Data[id]
	return url, ok
}

func (s *ShortenedUrlsByIDMemoryStorer) BatchSave(ctx context.Context, shortenedUrlsByIds []models.BatchPortionShortenResponse) error {
	for _, portion := range shortenedUrlsByIds {
		s.Data[portion.ID] = portion.OriginalURL
	}
	return nil
}

func (s *ShortenedUrlsByIDMemoryStorer) GetIDByURL(ctx context.Context, initialURL string) (string, error) {
	for id, url := range s.Data {
		if url == initialURL {
			return id, nil
		}
	}
	return "", nil
}

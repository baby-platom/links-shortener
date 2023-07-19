package shortid

import (
	"context"

	"github.com/baby-platom/links-shortener/internal/models"
)

// ShortenedUrlsByIDInterface represents ShortenedUrlsByID behaviour
type ShortenedUrlsByIDInterface interface {
	Save(ctx context.Context, id string, url string) error
	Get(ctx context.Context, id string) (string, bool)
	BatchSave(ctx context.Context, shortenedUrlsByIds []models.BatchPortionShortenResponse) error
}

// ShortenedUrlsByIDType stores id:url
type ShortenedUrlsByIDType struct {
	Data map[string]string
}

// NewShortenedUrlsByID return a new ShortenedUrlsByID
func NewShortenedUrlsByID() *ShortenedUrlsByIDType {
	return &ShortenedUrlsByIDType{Data: make(map[string]string)}
}

// Save creates new id:url relation
func (s *ShortenedUrlsByIDType) Save(ctx context.Context, id string, url string) error {
	s.Data[id] = url
	return nil
}

// Get returns url by id
func (s *ShortenedUrlsByIDType) Get(ctx context.Context, id string) (string, bool) {
	url, ok := s.Data[id]
	return url, ok
}

func (s *ShortenedUrlsByIDType) BatchSave(ctx context.Context, shortenedUrlsByIds []models.BatchPortionShortenResponse) error {
	for _, portion := range shortenedUrlsByIds {
		s.Data[portion.ID] = portion.OriginalURL
	}
	return nil
}

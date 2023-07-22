package storage

import (
	"context"

	"github.com/baby-platom/links-shortener/internal/models"
)

// ShortenedUrlsByIDStorer stores and provides access to shortened urls by ID.
type ShortenedUrlsByIDStorer interface {
	Save(ctx context.Context, id string, url string) error
	Get(ctx context.Context, id string) (string, bool)
	BatchSave(ctx context.Context, shortenedUrlsByIds []models.BatchPortionShortenResponse) error
	GetIDByURL(ctx context.Context, url string) (string, error)
}

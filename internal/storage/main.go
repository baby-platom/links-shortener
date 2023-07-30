package storage

import (
	"context"

	"github.com/baby-platom/links-shortener/internal/models"
)

// ShortenedUrlsByIDStorer stores and provides access to shortened urls by ID.
type ShortenedUrlsByIDStorer interface {
	Save(ctx context.Context, id string, url string, userID int) error
	Get(ctx context.Context, id string, userID int) (string, bool)
	BatchSave(ctx context.Context, shortenedUrlsByIds []models.BatchPortionShortenResponse, userID int) error
	GetIDByURL(ctx context.Context, url string, userID int) (string, error)
	GetUserShortenURLsList(ctx context.Context, baseAddress string, userIDToFind int) ([]models.UserShortenURLsListResponse, error)
}

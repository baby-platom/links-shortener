package storage

import (
	"context"

	"github.com/baby-platom/links-shortener/internal/models"
)

type deleteData struct {
	ids    []string
	userID string
}

// ShortenedUrlsByIDStorer stores and provides access to shortened urls by ID.
type ShortenedUrlsByIDStorer interface {
	Save(ctx context.Context, id string, url string, userID string) error
	Get(ctx context.Context, id string) (string, bool, bool)
	BatchSave(ctx context.Context, shortenedUrlsByIds []models.BatchPortionShortenResponse, userID string) error
	GetIDByURL(ctx context.Context, url string) (string, error)
	GetUserShortenURLsListResponse(ctx context.Context, baseAddress string, userIDToFind string) ([]models.UserShortenURLsListResponse, error)
	GetUserShortenURLsList(ctx context.Context, userIDToFind string) ([]string, error)
	BatchDelete(ctx context.Context, data []deleteData) error
	Delete(ctx context.Context, ids []string, userID string)
	GetDeleteCh() chan deleteData
}

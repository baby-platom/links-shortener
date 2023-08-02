package storage

import (
	"context"

	"github.com/baby-platom/links-shortener/internal/models"
)

type deleteData struct {
	ids    []string
	userID int
}

// ShortenedUrlsByIDStorer stores and provides access to shortened urls by ID.
type ShortenedUrlsByIDStorer interface {
	Save(ctx context.Context, id string, url string, userID int) error
	Get(ctx context.Context, id string) (string, bool, bool)
	BatchSave(ctx context.Context, shortenedUrlsByIds []models.BatchPortionShortenResponse, userID int) error
	GetIDByURL(ctx context.Context, url string) (string, error)
	GetUserShortenURLsListResponse(ctx context.Context, baseAddress string, userIDToFind int) ([]models.UserShortenURLsListResponse, error)
	GetUserShortenURLsList(ctx context.Context, userIDToFind int) ([]string, error)
	BacthDelete(ctx context.Context, data []deleteData) error
	MonitorDeleted(ctx context.Context)
	Delete(ctx context.Context, ids []string, userID int)
}

package storage

import (
	"context"

	"github.com/baby-platom/links-shortener/internal/database"
	"github.com/baby-platom/links-shortener/internal/models"
)

// ShortenedUrlsByIDDBStorer implements ShortenedUrlsByIDStorer interface using database.
type ShortenedUrlsByIDDBStorer struct {
	ShortenedUrlsByIDMemoryStorer
}

// CreateNewShortenedUrlsByIDDBStorer return a new ShortenedUrlsByIDDBStorer
func CreateNewShortenedUrlsByIDDBStorer() *ShortenedUrlsByIDDBStorer {
	return &ShortenedUrlsByIDDBStorer{*CreateNewShortenedUrlsByIDMemoryStorer()}
}

// Save creates new id:url relation and saves it to the json file
func (s *ShortenedUrlsByIDDBStorer) Save(ctx context.Context, id string, url string, userID int) error {
	return database.Connection.WriteShortenedURL(ctx, id, url, userID)
}

// Get returns url by id
func (s *ShortenedUrlsByIDDBStorer) Get(ctx context.Context, id string, userID int) (string, bool) {
	url, err := database.Connection.GetInitialURLByID(ctx, id, userID)
	if err != nil {
		panic(err)
	}

	ok := true
	if url == "" {
		ok = false
	}
	return url, ok
}

func (s *ShortenedUrlsByIDDBStorer) BatchSave(ctx context.Context, shortenedUrlsByIds []models.BatchPortionShortenResponse, userID int) error {
	if err := database.Connection.WriteBatchOfShortenedURL(ctx, shortenedUrlsByIds, userID); err != nil {
		return err
	}
	return nil
}

func (s *ShortenedUrlsByIDDBStorer) GetIDByURL(ctx context.Context, initialURL string, userID int) (string, error) {
	return database.Connection.GetIDByInitialURL(ctx, initialURL, userID)
}

func (s *ShortenedUrlsByIDDBStorer) GetUserShortenURLsList(ctx context.Context, baseAddress string, userID int) ([]models.UserShortenURLsListResponse, error) {
	return database.Connection.GetUserShortenURLsList(ctx, baseAddress, userID)
}

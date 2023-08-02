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
func (s *ShortenedUrlsByIDDBStorer) Save(ctx context.Context, id string, url string, userID string) error {
	return database.Connection.WriteShortenedURL(ctx, id, url, userID)
}

// Get returns url by id
func (s *ShortenedUrlsByIDDBStorer) Get(ctx context.Context, id string) (string, bool, bool) {
	url, deleted, err := database.Connection.GetInitialURLByID(ctx, id)
	if err != nil {
		panic(err)
	}

	ok := true
	if url == "" {
		ok = false
	}
	return url, ok, deleted
}

func (s *ShortenedUrlsByIDDBStorer) BatchSave(ctx context.Context, shortenedUrlsByIds []models.BatchPortionShortenResponse, userID string) error {
	if err := database.Connection.WriteBatchOfShortenedURL(ctx, shortenedUrlsByIds, userID); err != nil {
		return err
	}
	return nil
}

func (s *ShortenedUrlsByIDDBStorer) GetIDByURL(ctx context.Context, initialURL string) (string, error) {
	return database.Connection.GetIDByInitialURL(ctx, initialURL)
}

func (s *ShortenedUrlsByIDDBStorer) GetUserShortenURLsListResponse(ctx context.Context, baseAddress string, userID string) ([]models.UserShortenURLsListResponse, error) {
	return database.Connection.GetUserShortenURLsListResponse(ctx, baseAddress, userID)
}

func (s *ShortenedUrlsByIDDBStorer) GetUserShortenURLsList(ctx context.Context, userIDToFind string) ([]string, error) {
	return database.Connection.GetUserShortenURLsList(ctx, userIDToFind)
}

func (s *ShortenedUrlsByIDDBStorer) BacthDelete(ctx context.Context, data []deleteData) error {
	var ids []string
	for _, piece := range data {
		ids = append(ids, piece.ids...)
	}
	return database.Connection.BatchDelete(ctx, ids)
}

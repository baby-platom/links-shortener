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
func (s *ShortenedUrlsByIDDBStorer) Save(ctx context.Context, id string, url string) error {
	return database.Connection.WriteShortenedURL(ctx, id, url)
}

// Get returns url by id
func (s *ShortenedUrlsByIDDBStorer) Get(ctx context.Context, id string) (string, bool) {
	url, err := database.Connection.GetInitialURLByID(ctx, id)
	if err != nil {
		panic(err)
	}

	ok := true
	if url == "" {
		ok = false
	}
	return url, ok
}

func (s *ShortenedUrlsByIDDBStorer) BatchSave(ctx context.Context, shortenedUrlsByIds []models.BatchPortionShortenResponse) error {
	if err := database.Connection.WriteBatchOfShortenedURL(ctx, shortenedUrlsByIds); err != nil {
		return err
	}
	return nil
}

func (s *ShortenedUrlsByIDDBStorer) GetIDByURL(ctx context.Context, initialURL string) (string, error) {
	return database.Connection.GetIDByInitialURL(ctx, initialURL)
}

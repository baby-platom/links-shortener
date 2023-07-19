package shortid

import (
	"context"

	"github.com/baby-platom/links-shortener/internal/database"
	"github.com/baby-platom/links-shortener/internal/models"
)

type ShortenedUrlsByIDDatabaseType struct {
	ShortenedUrlsByIDType
}

// NewShortenedUrlsByIDDatabase return a new ShortenedUrlsByIDDatabase
func NewShortenedUrlsByIDDatabase() *ShortenedUrlsByIDDatabaseType {
	return &ShortenedUrlsByIDDatabaseType{*NewShortenedUrlsByID()}
}

// Save creates new id:url relation and saves it to the json file
func (s *ShortenedUrlsByIDDatabaseType) Save(ctx context.Context, id string, url string) error {
	return database.Connection.WriteShortenedURL(ctx, id, url)
}

// Get returns url by id
func (s *ShortenedUrlsByIDDatabaseType) Get(ctx context.Context, id string) (string, bool) {
	url, err := database.Connection.GetInitialURLLByIDByID(ctx, id)
	if err != nil {
		panic(err)
	}

	ok := true
	if url == "" {
		ok = false
	}
	return url, ok
}

func (s *ShortenedUrlsByIDDatabaseType) BatchSave(ctx context.Context, shortenedUrlsByIds []models.BatchPortionShortenResponse) error {
	if err := database.Connection.WriteBatchOfShortenedURL(ctx, shortenedUrlsByIds); err != nil {
		return err
	}
	return nil
}

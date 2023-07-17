package shortid

import (
	"context"

	"github.com/baby-platom/links-shortener/internal/database"
)

type ShortenedUrlsByIDDatabaseType struct {
	ShortenedUrlsByIDType
}

// NewShortenedUrlsByIDDatabase return a new ShortenedUrlsByIDDatabase
func NewShortenedUrlsByIDDatabase() *ShortenedUrlsByIDDatabaseType {
	return &ShortenedUrlsByIDDatabaseType{*NewShortenedUrlsByID()}
}

// Save creates new id:url relation and saves it to the json file
func (s *ShortenedUrlsByIDDatabaseType) Save(ctx context.Context, id string, url string) {
	err := database.Connection.WriteShortenedURL(ctx, id, url)
	if err != nil {
		panic(err)
	}
}

// Get returns url by id
func (s *ShortenedUrlsByIDDatabaseType) Get(ctx context.Context, id string) (string, bool) {
	err, url := database.Connection.GetShortenedURL(ctx, id)
	if err != nil {
		panic(err)
	}

	ok := true
	if url == "" {
		ok = false
	}
	return url, ok
}

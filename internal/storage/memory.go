package storage

import (
	"context"
	"fmt"

	"github.com/baby-platom/links-shortener/internal/models"
)

// ShortenedUrlsByIDMemoryStorer implements ShortenedUrlsByIDStorer interface using temproray memory.
type ShortenedUrlsByIDMemoryStorer struct {
	Data               map[string]string
	UsersShortenedURLs map[int][]string
}

// CreateNewShortenedUrlsByID return a new ShortenedUrlsByID
func CreateNewShortenedUrlsByIDMemoryStorer() *ShortenedUrlsByIDMemoryStorer {
	return &ShortenedUrlsByIDMemoryStorer{
		Data:               make(map[string]string),
		UsersShortenedURLs: make(map[int][]string),
	}
}

// Save creates new id:url relation
func (s *ShortenedUrlsByIDMemoryStorer) Save(ctx context.Context, id string, url string, userID int) error {
	s.Data[id] = url
	s.UsersShortenedURLs[userID] = append(s.UsersShortenedURLs[userID], id)
	return nil
}

// Get returns url by id
func (s *ShortenedUrlsByIDMemoryStorer) Get(ctx context.Context, id string) (string, bool) {
	url, ok := s.Data[id]
	return url, ok
}

func (s *ShortenedUrlsByIDMemoryStorer) BatchSave(ctx context.Context, shortenedUrlsByIds []models.BatchPortionShortenResponse, userID int) error {
	for _, portion := range shortenedUrlsByIds {
		s.Data[portion.ID] = portion.OriginalURL
		s.UsersShortenedURLs[userID] = append(s.UsersShortenedURLs[userID], portion.ID)
	}
	return nil
}

func (s *ShortenedUrlsByIDMemoryStorer) GetIDByURL(ctx context.Context, initialURL string) (string, error) {
	for id, url := range s.Data {
		if url == initialURL {
			return id, nil
		}
	}
	return "", nil
}

func (s *ShortenedUrlsByIDMemoryStorer) GetUserShortenURLsList(ctx context.Context, baseAddress string, userIDToFind int) ([]models.UserShortenURLsListResponse, error) {
	result := make([]models.UserShortenURLsListResponse, 0)

	shortURLs, ok := s.UsersShortenedURLs[userIDToFind]
	if ok {
		for _, shortURL := range shortURLs {
			initialURL := s.Data[shortURL]
			result = append(
				result,
				models.UserShortenURLsListResponse{
					ShortURL:    fmt.Sprintf("%s/%s", baseAddress, shortURL),
					OriginalURL: initialURL,
				},
			)
		}
	}

	return result, nil
}

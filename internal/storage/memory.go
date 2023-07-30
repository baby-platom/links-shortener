package storage

import (
	"context"
	"fmt"

	"github.com/baby-platom/links-shortener/internal/models"
)

// ShortenedUrlsByIDMemoryStorer implements ShortenedUrlsByIDStorer interface using temproray memory.
type ShortenedUrlsByIDMemoryStorer struct {
	Data map[int](map[string]string)
}

// CreateNewShortenedUrlsByID return a new ShortenedUrlsByID
func CreateNewShortenedUrlsByIDMemoryStorer() *ShortenedUrlsByIDMemoryStorer {
	return &ShortenedUrlsByIDMemoryStorer{Data: make(map[int](map[string]string))}
}

// Save creates new id:url relation
func (s *ShortenedUrlsByIDMemoryStorer) Save(ctx context.Context, id string, url string, userID int) error {
	if _, ok := s.Data[userID]; !ok {
		s.Data[userID] = make(map[string]string)
	}
	s.Data[userID][id] = url
	return nil
}

// Get returns url by id
func (s *ShortenedUrlsByIDMemoryStorer) Get(ctx context.Context, id string, userID int) (string, bool) {
	var url string
	data, ok := s.Data[userID]
	if ok {
		url, ok = data[id]
	}
	return url, ok
}

func (s *ShortenedUrlsByIDMemoryStorer) BatchSave(ctx context.Context, shortenedUrlsByIds []models.BatchPortionShortenResponse, userID int) error {
	if _, ok := s.Data[userID]; !ok {
		s.Data[userID] = make(map[string]string)
	}
	for _, portion := range shortenedUrlsByIds {
		s.Data[userID][portion.ID] = portion.OriginalURL
	}
	return nil
}

func (s *ShortenedUrlsByIDMemoryStorer) GetIDByURL(ctx context.Context, initialURL string, userIDToFind int) (string, error) {
	for userID, data := range s.Data {
		if userID == userIDToFind {
			for id, url := range data {
				if url == initialURL {
					return id, nil
				}
			}
		}
	}
	return "", nil
}

func (s *ShortenedUrlsByIDMemoryStorer) GetUserShortenURLsList(ctx context.Context, baseAddress string, userIDToFind int) ([]models.UserShortenURLsListResponse, error) {
	result := make([]models.UserShortenURLsListResponse, 0)
	for userID, data := range s.Data {
		if userID == userIDToFind {
			for id, url := range data {
				result = append(
					result,
					models.UserShortenURLsListResponse{
						ShortURL:    fmt.Sprintf("%s/%s", baseAddress, id),
						OriginalURL: url,
					},
				)
			}
			break
		}
	}
	return result, nil
}

package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/baby-platom/links-shortener/internal/logger"
	"github.com/baby-platom/links-shortener/internal/models"
)

type Data struct {
	initialURL string
	deleted    bool
}

// ShortenedUrlsByIDMemoryStorer implements ShortenedUrlsByIDStorer interface using temproray memory.
type ShortenedUrlsByIDMemoryStorer struct {
	Data               map[string]Data
	UsersShortenedURLs map[int][]string
	deleteCh           chan deleteData
}

// CreateNewShortenedUrlsByID return a new ShortenedUrlsByID
func CreateNewShortenedUrlsByIDMemoryStorer() *ShortenedUrlsByIDMemoryStorer {
	return &ShortenedUrlsByIDMemoryStorer{
		Data:               make(map[string]Data),
		UsersShortenedURLs: make(map[int][]string),
		deleteCh:           make(chan deleteData, 32),
	}
}

// Save creates new id:url relation
func (s *ShortenedUrlsByIDMemoryStorer) Save(ctx context.Context, id string, url string, userID int) error {
	s.Data[id] = Data{initialURL: url, deleted: false}
	s.UsersShortenedURLs[userID] = append(s.UsersShortenedURLs[userID], id)
	return nil
}

// Get returns url by id
func (s *ShortenedUrlsByIDMemoryStorer) Get(ctx context.Context, id string) (string, bool, bool) {
	data, ok := s.Data[id]
	return data.initialURL, ok, data.deleted
}

func (s *ShortenedUrlsByIDMemoryStorer) BatchSave(ctx context.Context, shortenedUrlsByIds []models.BatchPortionShortenResponse, userID int) error {
	for _, portion := range shortenedUrlsByIds {
		s.Data[portion.ID] = Data{initialURL: portion.OriginalURL, deleted: false}
		s.UsersShortenedURLs[userID] = append(s.UsersShortenedURLs[userID], portion.ID)
	}
	return nil
}

func (s *ShortenedUrlsByIDMemoryStorer) GetIDByURL(ctx context.Context, initialURL string) (string, error) {
	for id, data := range s.Data {
		if data.initialURL == initialURL {
			return id, nil
		}
	}
	return "", nil
}

func (s *ShortenedUrlsByIDMemoryStorer) GetUserShortenURLsListResponse(ctx context.Context, baseAddress string, userIDToFind int) ([]models.UserShortenURLsListResponse, error) {
	result := make([]models.UserShortenURLsListResponse, 0)

	shortURLs, ok := s.UsersShortenedURLs[userIDToFind]
	if ok {
		for _, shortURL := range shortURLs {
			initialURL := s.Data[shortURL].initialURL
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

func (s *ShortenedUrlsByIDMemoryStorer) GetUserShortenURLsList(ctx context.Context, userIDToFind int) ([]string, error) {
	shortURLs, _ := s.UsersShortenedURLs[userIDToFind]
	return shortURLs, nil
}

func (s *ShortenedUrlsByIDMemoryStorer) BacthDelete(ctx context.Context, data []deleteData) error {
	for _, piece := range data {
		var urlsToDelete []string
		for _, id := range piece.ids {
			new := s.Data[id]
			new.deleted = true
			s.Data[id] = new

			urlsToDelete = append(urlsToDelete, id)
		}

		var newURLs []string
		for _, url := range s.UsersShortenedURLs[piece.userID] {
			found := false
			for _, urlToDelete := range urlsToDelete {
				if url == urlToDelete {
					found = true
				}
			}

			if !found {
				newURLs = append(newURLs, url)
			}
		}
		s.UsersShortenedURLs[piece.userID] = newURLs
	}

	return nil
}

func (s *ShortenedUrlsByIDMemoryStorer) MonitorDeleted(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)

	var toDelete []deleteData

	for {
		select {
		case new := <-s.deleteCh:
			toDelete = append(toDelete, new)
		case <-ticker.C:
			if len(toDelete) == 0 {
				continue
			}

			err := s.BacthDelete(ctx, toDelete)
			if err != nil {
				logger.Log.Errorf("error occurred while deleting batch of shortened urls: %v", err)
			}
			toDelete = nil
		}
	}
}

func (s *ShortenedUrlsByIDMemoryStorer) Delete(ctx context.Context, ids []string, userID int) {
	data := deleteData{
		ids:    ids,
		userID: userID,
	}
	s.deleteCh <- data
}

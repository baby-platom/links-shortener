package storage

import (
	"context"
	"time"

	"github.com/baby-platom/links-shortener/internal/logger"
)

func MonitorDeleted(ctx context.Context, s ShortenedUrlsByIDStorer) {
	ticker := time.NewTicker(10 * time.Second)
	deleteCh := s.GetDeleteCh()

	var toDelete []deleteData

	for {
		select {
		case new := <-deleteCh:
			toDelete = append(toDelete, new)
		case <-ticker.C:
			if len(toDelete) == 0 {
				continue
			}

			err := s.BatchDelete(ctx, toDelete)
			if err != nil {
				logger.Log.Errorf("error occurred while deleting batch of shortened urls: %v", err)
			}
			toDelete = nil
		}
	}
}

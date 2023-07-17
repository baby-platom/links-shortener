package shortid

import (
    "context"
)

// ShortenedUrlsByIDInterface represents ShortenedUrlsByID behaviour
type ShortenedUrlsByIDInterface interface {
	Save(ctx context.Context, id string, url string)
	Get(ctx context.Context, id string) (string, bool)
}

// ShortenedUrlsByIDType stores id:url
type ShortenedUrlsByIDType struct {
	Data map[string]string
}

// NewShortenedUrlsByID return a new ShortenedUrlsByID
func NewShortenedUrlsByID() *ShortenedUrlsByIDType {
	return &ShortenedUrlsByIDType{Data: make(map[string]string)}
}

// Save creates new id:url relation
func (s *ShortenedUrlsByIDType) Save(ctx context.Context, id string, url string) {
	s.Data[id] = url
}

// Get returns url by id
func (s *ShortenedUrlsByIDType) Get(ctx context.Context, id string) (string, bool) {
	url, ok := s.Data[id]
	return url, ok
}

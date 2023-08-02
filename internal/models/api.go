package models

// ShortenRequest - inbound data
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenResponse - outbound data
type ShortenResponse struct {
	Result string `json:"result"`
}

// BatchPortionShortenRequest - inbound data
type BatchPortionShortenRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchPortionShortenResponse - outbound data
type BatchPortionShortenResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
	ID            string `json:"-"`
	OriginalURL   string `json:"-"`
}

// UserShortenURLsListResponse - outbound data
type UserShortenURLsListResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// BatchDeleteRequest - inbound data
type BatchDeleteRequest []string

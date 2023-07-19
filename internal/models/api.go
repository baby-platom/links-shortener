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
	CorrelationId string `json:"correlation_id"`
	OriginalUrl   string `json:"original_url"`
}

// BatchPortionShortenResponse - outbound data
type BatchPortionShortenResponse struct {
	CorrelationId string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
	ID            string `json:"-"`
	OriginalUrl   string `json:"-"`
}

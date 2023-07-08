package models

// ShortenRequest - inbound data
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenResponse - outbound data
type ShortenResponse struct {
	Result string `json:"result"`
}

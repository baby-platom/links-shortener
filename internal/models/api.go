package models

// ShortentRequest - inbound data
type ShortentRequest struct {
	URL string `json:"url"`
}

// ShortentResponse - outbound data
type ShortentResponse struct {
	Result string `json:"result"`
}

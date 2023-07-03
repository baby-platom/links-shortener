package shortid

import (
	"encoding/json"
	"os"
)

// ShortenedUrlsByIDType map id:initalUrl
type ShortenedUrlsByIDType map[string]string

// Save data into file
func (data ShortenedUrlsByIDType) Save(fname string) error {
	res, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		return err
	}
	return os.WriteFile(fname, res, 0666)
}

// Load data from a file
func (data *ShortenedUrlsByIDType) Load(fname string) error {
	res, err := os.ReadFile(fname)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(res, data); err != nil {
		return err
	}
	return nil
}

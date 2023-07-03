package main

import (
	"encoding/json"
	"fmt"
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

func main() {
	var ShortenedUrlsByID = make(ShortenedUrlsByIDType)
	var fname = "data.json"
	ShortenedUrlsByID["one"] = "id1"
	ShortenedUrlsByID["two"] = "id2"
	err := ShortenedUrlsByID.Save(fname)
	fmt.Println(err)
	var NewShortenedUrlsByID = make(ShortenedUrlsByIDType)
	err = NewShortenedUrlsByID.Load(fname)
	fmt.Println(err)
	for k, v := range NewShortenedUrlsByID {
		fmt.Printf("%s: %s\n", k, v)
	}
}

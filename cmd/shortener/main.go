package main

import (
	"github.com/baby-platom/links-shortener/internal/app"
	"github.com/baby-platom/links-shortener/internal/config"
)

func main() {
	config.ParseFlags()

	err := app.Run()
	if err != nil {
		panic(err)
	}
}

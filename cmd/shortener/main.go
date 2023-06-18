package main

import (
	"github.com/baby-platom/links-shortener/internal/app"
	"github.com/baby-platom/links-shortener/internal/config"
)

func main() {
	config.ParseFlags()

	if err := app.Run(); err != nil {
		panic(err)
	}
}

package config

import (
	"flag"
	"os"
)

// Config includes variables parsed from flags
var Config struct {
	Address     string
	BaseAddress string
}

// ParseFlags parses flags into the Config
func ParseFlags() {
	flag.StringVar(&Config.Address, "a", ":8080", "address and port to run server")
	flag.StringVar(&Config.BaseAddress, "b", "http://localhost:8080", "base address for shortened URLs")
	flag.Parse()

	if envAddress := os.Getenv("SERVER_ADDRESS"); envAddress != "" {
		Config.Address = envAddress
	}
	if envBaseAddress := os.Getenv("BASE_URL"); envBaseAddress != "" {
		Config.BaseAddress = envBaseAddress
	}
}

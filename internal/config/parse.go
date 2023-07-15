package config

import (
	"flag"
	"os"
)

// Config includes variables parsed from flags
var Config struct {
	Address         string
	BaseAddress     string
	LogLevel        string
	FileStoragePath string
	DatabaseDSN string
}

// ParseFlags parses flags into the Config
func ParseFlags() {
	flag.StringVar(&Config.Address, "a", ":8080", "address and port to run server")
	flag.StringVar(&Config.BaseAddress, "b", "http://localhost:8080", "base address for shortened URLs")
	flag.StringVar(&Config.LogLevel, "l", "info", "log level")
	flag.StringVar(&Config.FileStoragePath, "f", "/tmp/short-url-db.json", "file name for storing short URLs")
	flag.StringVar(&Config.DatabaseDSN, "d", "", "database DSN")
	flag.Parse()

	if envAddress := os.Getenv("SERVER_ADDRESS"); envAddress != "" {
		Config.Address = envAddress
	}
	if envBaseAddress := os.Getenv("BASE_URL"); envBaseAddress != "" {
		Config.BaseAddress = envBaseAddress
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		Config.LogLevel = envLogLevel
	}
	if fileStoragePath := os.Getenv("FILE_STORAGE_PATH"); fileStoragePath != "" {
		Config.FileStoragePath = fileStoragePath
	}
	if databaseDSN := os.Getenv("DATABASE_DSN"); databaseDSN != "" {
		Config.DatabaseDSN = databaseDSN
	}
}

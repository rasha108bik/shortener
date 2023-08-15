package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

// Config include all parametrs for run server
// SERVER_ADDRESS, BASE_URL, FILE_STORAGE_PATH, DATABASE_DSN.
type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
	DatabaseDSN     string `env:"DATABASE_DSN" envDefault:""`
	EnableHTTPS     string `env:"ENABLE_HTTPS" envDefault:""`
}

var (
	serverAddress     string
	baseURL           string
	fileStoragePath   string
	databaseDSN       string
	listenAndServeTLS string
)

// NewConfig returns a newly initialized Config objects.
func NewConfig() *Config {
	flag.StringVar(&serverAddress, "a", "", "server address")
	flag.StringVar(&baseURL, "b", "", "base URL")
	flag.StringVar(&fileStoragePath, "f", "", "file storage path")
	flag.StringVar(&databaseDSN, "d", "", "databaseDSN path")
	flag.StringVar(&listenAndServeTLS, "s", "", "start ListenAndServeTLS")

	flag.Parse()
	log.Printf("server address: %s, base URL: %s, file storagePath: %s databaseDSN: %s\n", serverAddress, baseURL, fileStoragePath, databaseDSN)

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if serverAddress != "" {
		cfg.ServerAddress = serverAddress
	}
	if baseURL != "" {
		cfg.BaseURL = baseURL
	}

	cfg.FileStoragePath = fileStoragePath
	cfg.DatabaseDSN = databaseDSN

	if listenAndServeTLS != "" {
		cfg.EnableHTTPS = listenAndServeTLS
	}

	return &cfg
}

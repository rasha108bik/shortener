package main

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"

	"github.com/rasha108bik/tiny_url/config"
	"github.com/rasha108bik/tiny_url/internal/app"
)

var (
	serverAddress   string
	baseURL         string
	fileStoragePath string
)

func init() {
	flag.StringVar(&serverAddress, "a", "", "server address")
	flag.StringVar(&baseURL, "b", "", "base URL")
	flag.StringVar(&fileStoragePath, "f", "", "file storage path")
}

func main() {
	flag.Parse()
	log.Printf("server address: %s, base URL: %s, file storagePath: %s\n", serverAddress, baseURL, fileStoragePath)

	var cfg config.Config
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
	if fileStoragePath != "" {
		cfg.FileStoragePath = fileStoragePath
	}

	log.Printf("%+v\n", cfg)

	app.Run(&cfg)
}

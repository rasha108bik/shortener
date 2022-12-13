package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"

	"github.com/rasha108bik/tiny_url/config"
	"github.com/rasha108bik/tiny_url/internal/server"
	"github.com/rasha108bik/tiny_url/internal/server/handlers"
	"github.com/rasha108bik/tiny_url/internal/storage"
	"github.com/rasha108bik/tiny_url/pkg/storagefile"
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

	// TODO: create app direction
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

	fileName := cfg.FileStoragePath
	producer, err := storagefile.NewProducer(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()
	consumer, err := storagefile.NewConsumer(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	db := storage.NewStorage()
	h := handlers.NewHandler(db, &cfg, producer, consumer)
	serv := server.NewServer(h)

	r := chi.NewRouter()
	r.MethodNotAllowed(serv.Handlers.ErrorHandler)
	r.Get("/{id}", serv.Handlers.GetOriginalURL)
	r.Post("/api/shorten", serv.Handlers.CreateShorten)
	r.Post("/", serv.Handlers.CreateShortLink)

	err = http.ListenAndServe(cfg.ServerAddress, r)
	if err != nil {
		log.Fatal(err)
	}
}

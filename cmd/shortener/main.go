package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/rasha108bik/tiny_url/config"
	"github.com/rasha108bik/tiny_url/internal/middleware"
	"github.com/rasha108bik/tiny_url/internal/server"
	"github.com/rasha108bik/tiny_url/internal/server/handlers"
	"github.com/rasha108bik/tiny_url/internal/storage"
	"github.com/rasha108bik/tiny_url/pkg/storagefile"
)

func main() {
	cfg := config.NewConfig()

	log.Printf("%+v\n", cfg)

	fileName := cfg.FileStoragePath
	writer, err := storagefile.NewProducer(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer writer.Close()
	reader, err := storagefile.NewConsumer(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	db := storage.NewStorage()
	h := handlers.NewHandler(db, cfg, writer, reader)
	serv := server.NewServer(h)

	r := chi.NewRouter()
	r.Use(middleware.GzipHandle)
	r.Use(middleware.GzipRequest)
	r.MethodNotAllowed(serv.Handlers.ErrorHandler)
	r.Get("/{id}", serv.Handlers.GetOriginalURL)
	r.Post("/api/shorten", serv.Handlers.CreateShorten)
	r.Post("/", serv.Handlers.CreateShortLink)

	err = http.ListenAndServe(cfg.ServerAddress, r)
	if err != nil {
		log.Fatal(err)
	}
}

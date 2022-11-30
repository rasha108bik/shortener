package app

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

func Run(cfg *config.Config) {
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
	h := handlers.NewHandler(db, cfg, producer, consumer)
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

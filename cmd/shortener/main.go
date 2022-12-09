package main

import (
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"

	"github.com/rasha108bik/tiny_url/config"
	"github.com/rasha108bik/tiny_url/internal/server"
	"github.com/rasha108bik/tiny_url/internal/server/handlers"
	"github.com/rasha108bik/tiny_url/internal/storage"
)

func main() {
	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v\n", cfg)

	db := storage.NewStorage()
	h := handlers.NewHandler(db, &cfg)
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

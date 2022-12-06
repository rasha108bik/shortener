package main

import (
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"

	"github.com/rasha108bik/tiny_url/internal/config"
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

	if cfg.BaseURL == "" {
		cfg.BaseURL = "/"
	}
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = "127.0.0.1:8080"
	}
	log.Printf("%+v\n", cfg)

	db := storage.NewStorage()
	h := handlers.NewHandler(db, &cfg)
	serv := server.NewServer(h)

	r := chi.NewRouter()
	r.MethodNotAllowed(serv.Handlers.ErrorHandler)
	r.Route(cfg.BaseURL, func(r chi.Router) {
		r.Get("/{id}", serv.Handlers.GetOriginalURL)
		r.Post("/", serv.Handlers.CreateShortLink)
		r.Post("/api/shorten", serv.Handlers.CreateShorten)
	})

	err = http.ListenAndServe(cfg.ServerAddress, r)
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/rasha108bik/tiny_url/internal/server"
	"github.com/rasha108bik/tiny_url/internal/server/handlers"
	"github.com/rasha108bik/tiny_url/internal/storage"
)

func main() {
	r := chi.NewRouter()

	db := storage.NewStorage()
	h := handlers.NewHandler(db)
	serv := server.NewServer(h)

	r.MethodNotAllowed(serv.Handlers.ErrorHandler)
	r.Get("/{id}", serv.Handlers.GetOriginalURL)
	r.Post("/", serv.Handlers.CreateShortLink)

	err := http.ListenAndServe("127.0.0.1:8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

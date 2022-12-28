package main

import (
	"context"
	"log"
	"net/http"

	"github.com/rasha108bik/tiny_url/config"
	"github.com/rasha108bik/tiny_url/internal/router"
	"github.com/rasha108bik/tiny_url/internal/server"
	"github.com/rasha108bik/tiny_url/internal/server/handlers"
	storage "github.com/rasha108bik/tiny_url/internal/storage/db"
	filestorage "github.com/rasha108bik/tiny_url/internal/storage/file"
	"github.com/rasha108bik/tiny_url/internal/storage/postgres"
)

func main() {
	cfg := config.NewConfig()

	log.Printf("%+v\n", cfg)

	pg, err := postgres.New(context.Background(), cfg.DatabaseDSN, postgres.MaxPoolSize(4))
	if err != nil {
		log.Printf("postgres.New: %v\n", err)
	}
	defer pg.Close()

	filestorage, err := filestorage.NewFileStorage(cfg.FileStoragePath)
	if err != nil {
		log.Fatal(err)
	}
	defer filestorage.Close()

	db := storage.NewStorage()
	h := handlers.NewHandler(cfg, db, filestorage, pg)
	serv := server.NewServer(h)
	r := router.NewRouter(serv)

	err = http.ListenAndServe(cfg.ServerAddress, r)
	if err != nil {
		log.Fatal(err)
	}
}

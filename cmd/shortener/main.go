package main

import (
	"context"
	"fmt"
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
		log.Fatal(fmt.Errorf("postgres.New: %w", err))
	}
	defer pg.Close()

	fileName := cfg.FileStoragePath
	filestorage, err := filestorage.NewFileStorage(fileName)
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

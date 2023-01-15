package main

import (
	"log"
	"net/http"

	"github.com/rasha108bik/tiny_url/config"
	"github.com/rasha108bik/tiny_url/internal/router"
	"github.com/rasha108bik/tiny_url/internal/server"
	"github.com/rasha108bik/tiny_url/internal/server/handlers"
	storage "github.com/rasha108bik/tiny_url/internal/storage/db"
	filestorage "github.com/rasha108bik/tiny_url/internal/storage/file"
)

func main() {
	cfg := config.NewConfig()

	log.Printf("%+v\n", cfg)

	fileName := cfg.FileStoragePath
	filestorage, err := filestorage.NewFileStorage(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer filestorage.Close()

	db := storage.NewStorage()
	h := handlers.NewHandler(cfg, db, filestorage)
	serv := server.NewServer(h)
	r := router.NewRouter(serv)

	err = http.ListenAndServe(cfg.ServerAddress, r)
	if err != nil {
		log.Fatal(err)
	}
}

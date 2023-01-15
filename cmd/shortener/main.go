package main

import (
	"os"

	"github.com/rasha108bik/tiny_url/config"
	"github.com/rasha108bik/tiny_url/internal/server"
	"github.com/rasha108bik/tiny_url/internal/server/handlers"
	"github.com/rasha108bik/tiny_url/internal/storager"
	"github.com/rs/zerolog"
)

func main() {
	log := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)

	cfg := config.NewConfig()

	log.Printf("%+v\n", cfg)

	str, err := storager.NewStorager(cfg)
	if err != nil {
		log.Printf("pgDB.New: %v\n", err)
	}
	defer str.Close()

	h := handlers.NewHandler(cfg, str)
	serv := server.NewServer(h, cfg.ServerAddress)

	err = serv.ListenAndServe()
	if err != nil {
		log.Fatal().Err(err)
	}
}

package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"

	"github.com/rasha108bik/tiny_url/config"
	"github.com/rasha108bik/tiny_url/internal/server"
	"github.com/rasha108bik/tiny_url/internal/server/handlers"
	"github.com/rasha108bik/tiny_url/internal/storager"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n",
		buildVersion, buildDate, buildCommit)

	log := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)

	cfg := config.NewConfig()

	log.Printf("%+v\n", cfg)

	str, err := storager.NewStorager(cfg)
	if err != nil {
		log.Printf("pgDB.New: %v\n", err)
	}
	defer str.Close()

	h := handlers.NewHandler(&log, cfg, str)

	err = server.NewServer(h, cfg.ServerAddress, cfg.EnableHTTPS).Start(&log)
	if err != nil {
		log.Fatal().Err(err)
	}
}

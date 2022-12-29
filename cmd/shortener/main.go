package main

import (
	"log"
	"net/http"

	//	"github.com/golang-migrate/migrate/v4"
	//	"github.com/golang-migrate/migrate/v4/database/postgres"
	//	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rasha108bik/tiny_url/config"
	"github.com/rasha108bik/tiny_url/internal/router"
	"github.com/rasha108bik/tiny_url/internal/server"
	"github.com/rasha108bik/tiny_url/internal/server/handlers"
	filestorage "github.com/rasha108bik/tiny_url/internal/storage/file"
	storage "github.com/rasha108bik/tiny_url/internal/storage/memdb"
	pgDB "github.com/rasha108bik/tiny_url/internal/storage/postgres"
)

func main() {
	cfg := config.NewConfig()

	log.Printf("%+v\n", cfg)

	pg, err := pgDB.New(cfg.DatabaseDSN)
	if err != nil {
		log.Printf("pgDB.New: %v\n", err)
	}
	defer pg.Close()

	//	if cfg.DatabaseDSN != "" {
	//		driver, err := postgres.WithInstance(pg.Postgres, &postgres.Config{})
	//		if err != nil {
	//			log.Printf("postgres.WithInstance: %v\n", err)
	//		}
	//
	//		m, err := migrate.NewWithDatabaseInstance(
	//			"./migrations/",
	//			"postgres", driver)
	//		if err != nil {
	//			log.Printf("migrate.NewWithDatabaseInstance: %v\n", err)
	//		}
	//
	//		err = m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run
	//		if err != nil && err != migrate.ErrNoChange {
	//			log.Fatal(fmt.Errorf("migrate failed: %v", err))
	//		}
	//	}

	filestorage, err := filestorage.NewFileStorage(cfg.FileStoragePath)
	if err != nil {
		log.Fatal(err)
	}
	defer filestorage.Close()

	memDB := storage.NewMemDB()
	h := handlers.NewHandler(cfg, memDB, filestorage, pg)
	serv := server.NewServer(h)
	r := router.NewRouter(serv)

	err = http.ListenAndServe(cfg.ServerAddress, r)
	if err != nil {
		log.Fatal(err)
	}
}

package storager

import (
	"context"
	"log"

	"github.com/rasha108bik/tiny_url/config"
	filestorage "github.com/rasha108bik/tiny_url/internal/storager/file"
	storage "github.com/rasha108bik/tiny_url/internal/storager/memdb"
	pgDB "github.com/rasha108bik/tiny_url/internal/storager/postgres"
)

//go:generate bin/mockgen -source=postgres.go -package=$GOPACKAGE -destination=postgres_interface_mock.go

// Storager interface for storage object.
type Storager interface {
	StoreURL(ctx context.Context, originalURL string, shortURL string) error
	GetOriginalURLByShortURL(ctx context.Context, shortURL string) (string, error)
	GetAllURLs(ctx context.Context) (map[string]string, error)
	GetShortURLByOriginalURL(ctx context.Context, originalURL string) (string, error)
	DeleteURLByShortURL(ctx context.Context, shortlURLs string) error
	Ping(ctx context.Context) error
	Close() error
}

// StorageType type storage.
type StorageType int

// Generic storager type.
const (
	InMemoryStorage StorageType = 1 << iota
	FileStorage
	PsgStorage
)

// NewStorager  returns a newly initialized Storager objects that implements the Storager
// interface.
func NewStorager(cfg *config.Config) (Storager, error) {
	switch getStoragerType(*cfg) {
	case PsgStorage:
		return makePostgres(cfg.DatabaseDSN), nil
	case FileStorage:
		return makeFileStorage(cfg.FileStoragePath), nil
	default:
		return makeInMemroy(), nil
	}
}

func getStoragerType(cfg config.Config) StorageType {
	if cfg.DatabaseDSN != "" {
		return PsgStorage
	} else if cfg.FileStoragePath != "" {
		return FileStorage
	}
	return InMemoryStorage
}

func makePostgres(databaseDSN string) Storager {
	psg, err := pgDB.NewPostgres(databaseDSN)
	if err != nil {
		log.Printf("pgDB.New: %v\n", err)
	}

	return psg
}

func makeFileStorage(fileStoragePath string) Storager {
	fileStorage, err := filestorage.NewFileStorage(fileStoragePath)
	if err != nil {
		log.Printf("filestorage.NewFileStorage: %v", err)
	}

	return fileStorage
}

func makeInMemroy() Storager {
	return storage.NewMemDB()
}

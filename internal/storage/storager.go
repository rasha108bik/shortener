package storage

import "context"

//go:generate mockgen -source=Storager -package=$GOPACKAGE -destination=postgres_interface_mock.go

type Storager interface {
	StoreURL(originalURL string, shortURL string) error
	GetOriginalURLByShortURL(shortURL string) (string, error)
	GetAllURLs() (map[string]string, error)
	GetShortURLByOriginalURL(originalURL string) (string, error)
	Ping(ctx context.Context) error
}

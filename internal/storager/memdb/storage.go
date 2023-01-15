package storage

import (
	"context"
	"net/url"

	appErr "github.com/rasha108bik/tiny_url/internal/errors"
)

type memDB struct {
	Locations map[string]string
}

func NewMemDB() *memDB {
	return &memDB{
		Locations: make(map[string]string),
	}
}

func (f *memDB) StoreURL(originalURL string, shortURL string) error {
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		return err
	}

	f.Locations[shortURL] = originalURL
	return nil
}

func (f *memDB) GetOriginalURLByShortURL(shortURL string) (string, error) {
	if url, ok := f.Locations[shortURL]; ok {
		return url, nil
	}
	return "", appErr.ErrNoSuchID
}

func (f *memDB) GetAllURLs() (map[string]string, error) {
	return f.Locations, nil
}

func (f *memDB) GetShortURLByOriginalURL(originalURL string) (string, error) {
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		return "", err
	}
	for shURL, orURL := range f.Locations {
		if orURL == originalURL {
			return shURL, appErr.ErrOriginalURLExist
		}
	}

	return "", appErr.ErrNoSuchID
}

func (f *memDB) Ping(ctx context.Context) error {
	return nil
}

func (f *memDB) Close() error {
	return nil
}

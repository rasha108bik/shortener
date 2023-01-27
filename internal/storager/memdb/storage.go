package storage

import (
	"context"
	"net/url"
	"sync"

	appErr "github.com/rasha108bik/tiny_url/internal/errors"
)

type memDB struct {
	locations map[string]string
	mtx       sync.RWMutex
}

func NewMemDB() *memDB {
	return &memDB{
		locations: make(map[string]string),
	}
}

func (m *memDB) StoreURL(originalURL string, shortURL string) error {
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		return err
	}

	m.mtx.RLock()
	m.locations[shortURL] = originalURL
	m.mtx.RUnlock()

	return nil
}

func (m *memDB) GetOriginalURLByShortURL(shortURL string) (string, error) {
	if url, ok := m.locations[shortURL]; ok {
		return url, nil
	}
	return "", appErr.ErrNoSuchID
}

func (m *memDB) GetAllURLs() (map[string]string, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return m.locations, nil
}

func (m *memDB) GetShortURLByOriginalURL(originalURL string) (string, error) {
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		return "", err
	}

	m.mtx.RLock()
	defer m.mtx.RUnlock()
	for shURL, orURL := range m.locations {
		if orURL == originalURL {
			return shURL, appErr.ErrOriginalURLExist
		}
	}

	return "", appErr.ErrNoSuchID
}

func (m *memDB) Ping(ctx context.Context) error {
	return nil
}

func (m *memDB) Close() error {
	return nil
}

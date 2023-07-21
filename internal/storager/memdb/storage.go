// module memDB, which implement the Storager interface methods
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

// NewMemDB returns a newly initialized memDB objects that implements the Storager
// interface.
func NewMemDB() *memDB {
	return &memDB{
		locations: make(map[string]string),
	}
}

// StoreURL save original URL and short URL in memdb.
func (m *memDB) StoreURL(_ context.Context, originalURL string, shortURL string) error {
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		return err
	}

	m.mtx.RLock()
	m.locations[shortURL] = originalURL
	m.mtx.RUnlock()

	return nil
}

// GetOriginalURLByShortURL get original URL by short URL.
func (m *memDB) GetOriginalURLByShortURL(_ context.Context, shortURL string) (string, error) {
	if url, ok := m.locations[shortURL]; ok {
		return url, nil
	}

	return "", appErr.ErrURLDeleted
}

// GetAllURLs get all URLs from memdb.
func (m *memDB) GetAllURLs(_ context.Context) (map[string]string, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return m.locations, nil
}

// GetShortURLByOriginalURL get URL by original URL
func (m *memDB) GetShortURLByOriginalURL(_ context.Context, originalURL string) (string, error) {
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

// Ping not implement.
func (m *memDB) Ping(_ context.Context) error {
	return nil
}

// Close not implement.
func (m *memDB) Close() error {
	return nil
}

// DeleteURLByShortURL delete long URL by short URL.
func (m *memDB) DeleteURLByShortURL(_ context.Context, shortURL string) error {
	delete(m.locations, shortURL)
	return nil
}

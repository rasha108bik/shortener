package storage

import (
	"errors"
	"net/url"
)

type memDB struct {
	Locations map[string]string
}

func NewMemDB() *memDB {
	return &memDB{
		Locations: make(map[string]string),
	}
}

func (f *memDB) StoreURL(longURL string, shortURL string) error {
	if _, err := url.ParseRequestURI(longURL); err != nil {
		return err
	}

	// lastID := len(f.Locations)
	// newID := fmt.Sprint(lastID + 1)
	f.Locations[shortURL] = longURL
	return nil
}

func (f *memDB) GetURLShortID(id string) (string, error) {
	if url, ok := f.Locations[id]; ok {
		return url, nil
	}
	return "", errors.New("no such id")
}

func (f *memDB) GetURLsShort() map[string]string {
	return f.Locations
}

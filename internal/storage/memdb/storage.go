package storage

import (
	"errors"
	"fmt"
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

func (f *memDB) StoreURL(longURL string) (string, error) {
	if _, err := url.ParseRequestURI(longURL); err != nil {
		return "", err
	}

	lastID := len(f.Locations)
	newID := fmt.Sprint(lastID + 1)
	f.Locations[newID] = longURL
	return newID, nil
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

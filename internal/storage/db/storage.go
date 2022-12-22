package storage

import (
	"errors"
	"fmt"
	"net/url"
)

type storage struct {
	Locations map[string]string
}

func NewStorage() *storage {
	return &storage{
		Locations: make(map[string]string),
	}
}

func (f *storage) StoreURL(longURL string) (string, error) {
	if _, err := url.ParseRequestURI(longURL); err != nil {
		return "", err
	}

	lastID := len(f.Locations)
	newID := fmt.Sprint(lastID + 1)
	f.Locations[newID] = longURL
	return newID, nil
}

func (f *storage) GetURLShortID(id string) (string, error) {
	if url, ok := f.Locations[id]; ok {
		return url, nil
	}
	return "", errors.New("no such id")
}

func (f *storage) GetURLsShort() map[string]string {
	return f.Locations
}

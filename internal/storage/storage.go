package storage

import (
	"errors"
	"fmt"
	"net/url"
)

type Storage interface {
	StoreURL(longURL string) (string, error)
	GetURLShortID(id string) (string, error)
}

type storage struct {
	Locations map[string]string
}

func NewStorage() *storage {
	return &storage{
		Locations: make(map[string]string),
	}
}

func (s *storage) StoreURL(longURL string) (string, error) {
	if _, err := url.ParseRequestURI(longURL); err != nil {
		return "", err
	}

	lastID := len(s.Locations)
	newID := fmt.Sprint(lastID + 1)
	s.Locations[newID] = longURL
	return newID, nil
}

func (s *storage) GetURLShortID(id string) (string, error) {
	if url, ok := s.Locations[id]; ok {
		return url, nil
	}
	return "", errors.New("no such id")
}

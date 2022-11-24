package storage

import (
	"errors"
	"fmt"
	"net/url"
)

type Storage interface {
	Save(longURL string) (string, error)
	Get(id string) (string, error)
}

type storage struct {
	Locations map[string]string
}

func New() *storage {
	return &storage{
		Locations: make(map[string]string),
	}
}

func (s *storage) Save(longURL string) (string, error) {
	if _, err := url.ParseRequestURI(longURL); err != nil {
		return "", err
	}

	lastID := len(s.Locations)
	newID := fmt.Sprint(lastID + 1)
	s.Locations[newID] = longURL
	return newID, nil
}

func (s *storage) Get(id string) (string, error) {
	if url, ok := s.Locations[id]; ok {
		return url, nil
	}
	return "", errors.New("no such id")
}

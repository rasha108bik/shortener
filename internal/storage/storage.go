package storage

import (
	"errors"
	"fmt"
	"net/url"
)

type Storage struct {
	Locations map[string]string
}

func New() *Storage {
	return &Storage{
		Locations: make(map[string]string),
	}
}

func (s *Storage) Save(longURL string) (string, error) {
	if _, err := url.ParseRequestURI(longURL); err != nil {
		return "", err
	}

	lastID := len(s.Locations)
	newID := fmt.Sprint(lastID + 1)
	s.Locations[newID] = longURL
	return newID, nil
}

func (s *Storage) Get(id string) (string, error) {
	if url, ok := s.Locations[id]; ok {
		return url, nil
	}
	return "", errors.New("no such id")
}

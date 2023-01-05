package storage

import (
	"encoding/json"
	"errors"
	"os"
)

type Event struct {
	URL string
}

type fileStorage struct {
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

func NewFileStorage(filename string) (*fileStorage, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}
	return &fileStorage{
		file:    file,
		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),
	}, nil
}

func (f *fileStorage) StoreURL(originalURL string, shortURL string) error {
	return f.encoder.Encode(&shortURL)
}

func (f *fileStorage) GetURLShortID(id string) (string, error) {
	events := []Event{}
	if err := f.decoder.Decode(&events); err != nil {
		return "", err
	}

	for _, event := range events {
		if event.URL == id {
			return event.URL, nil
		}
	}

	return "", errors.New("no such id")
}

func (f *fileStorage) GetURLsShort() map[string]string {
	return nil
}

func (f *fileStorage) Close() error {
	return f.file.Close()
}

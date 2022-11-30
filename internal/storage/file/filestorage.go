package storage

import (
	"context"
	"encoding/json"
	"os"

	appErr "github.com/rasha108bik/tiny_url/internal/errors"
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

func (f *fileStorage) GetOriginalURLByShortURL(shortURL string) (string, error) {
	events := []Event{}
	if err := f.decoder.Decode(&events); err != nil {
		return "", err
	}

	for _, event := range events {
		if event.URL == shortURL {
			return event.URL, nil
		}
	}

	return "", appErr.ErrNoSuchID
}

func (f *fileStorage) GetAllURLs() (map[string]string, error) {
	return nil, nil
}

func (f *fileStorage) GetShortURLByOriginalURL(originalURL string) (string, error) {
	return "", nil
}

func (f *fileStorage) Ping(ctx context.Context) error {
	return nil
}

func (f *fileStorage) Close() error {
	return f.file.Close()
}

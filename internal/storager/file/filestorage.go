package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"

	appErr "github.com/rasha108bik/tiny_url/internal/errors"
)

type Event struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type fileStorage struct {
	file *os.File
}

func NewFileStorage(filename string) (*fileStorage, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &fileStorage{
		file: file,
	}, nil
}

func (f *fileStorage) StoreURL(originalURL string, shortURL string) error {
	data, err := json.Marshal(
		&Event{
			OriginalURL: originalURL,
			ShortURL:    shortURL,
		},
	)
	if err != nil {
		return err
	}

	_, err = f.file.Write(data)
	if err != nil {
		return err
	}

	_, err = f.file.Write([]byte("\n"))
	if err != nil {
		return err
	}

	err = f.file.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (f *fileStorage) GetOriginalURLByShortURL(shortURL string) (string, error) {
	f.file.Seek(0, io.SeekStart)

	scanner := bufio.NewScanner(f.file)
	for scanner.Scan() {
		data := scanner.Bytes()
		event := Event{}
		err := json.Unmarshal(data, &event)
		if err != nil {
			break
		}

		if event.ShortURL == shortURL {
			return event.OriginalURL, nil
		}
	}

	return "", appErr.ErrNoSuchID
}

func (f *fileStorage) GetAllURLs() (map[string]string, error) {
	f.file.Seek(0, io.SeekStart)

	res := make(map[string]string)

	scanner := bufio.NewScanner(f.file)
	for scanner.Scan() {
		data := scanner.Bytes()
		event := Event{}
		err := json.Unmarshal(data, &event)
		if err != nil {
			break
		}
		res[event.ShortURL] = event.OriginalURL
	}

	return res, nil
}

func (f *fileStorage) GetShortURLByOriginalURL(originalURL string) (string, error) {
	f.file.Seek(0, io.SeekStart)

	scanner := bufio.NewScanner(f.file)
	for scanner.Scan() {
		data := scanner.Bytes()
		event := Event{}
		err := json.Unmarshal(data, &event)
		if err != nil {
			break
		}

		if event.OriginalURL == originalURL {
			return event.ShortURL, appErr.ErrOriginalURLExist
		}
	}

	return "", appErr.ErrNoSuchID
}

func (f *fileStorage) Close() error {
	return f.file.Close()
}

func (f *fileStorage) Ping(ctx context.Context) error {
	return nil
}

// module file, which implement the Storager interface methods
package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"

	appErr "github.com/rasha108bik/tiny_url/internal/errors"
)

// models for filestorage
type (
	// Event model for save in file the fiels: short_url, original_url
	Event struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	fileStorage struct {
		file *os.File
	}
)

// NewFileStorage returns a newly initialized fileStorage objects that implements the Storager
// interface.
func NewFileStorage(filename string) (*fileStorage, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &fileStorage{
		file: file,
	}, nil
}

// StoreURL save original URL and short URL in file.
func (f *fileStorage) StoreURL(ctx context.Context, originalURL string, shortURL string) error {
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

// GetOriginalURLByShortURL get original URL by short URL.
func (f *fileStorage) GetOriginalURLByShortURL(ctx context.Context, shortURL string) (string, error) {
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

// GetAllURLs get all URLs from file.
func (f *fileStorage) GetAllURLs(ctx context.Context) (map[string]string, error) {
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

// GetShortURLByOriginalURL get URL by original URL
func (f *fileStorage) GetShortURLByOriginalURL(ctx context.Context, originalURL string) (string, error) {
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

// Close closing file descriptor.
func (f *fileStorage) Close() error {
	return f.file.Close()
}

// Ping not implement.
func (f *fileStorage) Ping(ctx context.Context) error {
	return nil
}

// DeleteURLByShortURL not implement.
func (f *fileStorage) DeleteURLByShortURL(ctx context.Context, shortlURL string) error {
	return nil
}

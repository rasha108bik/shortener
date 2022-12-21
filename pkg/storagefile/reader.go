package storagefile

import (
	"encoding/json"
	"os"
)

type Reader interface {
	ReadEvent() ([]Event, error)
	Close() error
}

type reader struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(fileName string) (*reader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}
	return &reader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *reader) ReadEvent() ([]Event, error) {
	event := []Event{}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *reader) Close() error {
	return c.file.Close()
}

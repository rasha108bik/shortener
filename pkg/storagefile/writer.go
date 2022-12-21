package storagefile

import (
	"encoding/json"
	"os"
)

type Writer interface {
	WriteEvent(event *Event) error
	Close() error
}

type Event struct {
	URL string
}

type writer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(filename string) (*writer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &writer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *writer) WriteEvent(event *Event) error {
	return p.encoder.Encode(&event)
}

func (p *writer) Close() error {
	return p.file.Close()
}

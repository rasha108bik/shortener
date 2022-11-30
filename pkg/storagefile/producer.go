package storagefile

import (
	"encoding/json"
	"os"
)

type Producer interface {
	WriteEvent(event *Event) error
	Close() error
}

type Event struct {
	URL string
}

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(filename string) (*producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *producer) WriteEvent(event *Event) error {
	return p.encoder.Encode(&event)
}

func (p *producer) Close() error {
	return p.file.Close()
}

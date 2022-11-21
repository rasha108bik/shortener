package storage

import (
	"encoding/base64"
	"errors"
)

type Storage interface {
	Save(string) (string, error)
	Get(string) (string, error)
}

type DB struct {
	db map[string]string
}

func New() Storage {
	m := make(map[string]string, 0)
	return &DB{
		db: m,
	}
}

func (d *DB) Save(data string) (string, error) {
	sEnc := base64.StdEncoding.EncodeToString([]byte(data))
	d.db[sEnc] = data

	return sEnc, nil
}

func (d *DB) Get(data string) (string, error) {
	// for k, v := range d.db {
	// 	fmt.Printf("k: %s v: %s \n", k, v)
	// }

	origURL := d.db[data]
	if origURL == "" {
		return "", errors.New("key not found")
	}

	return origURL, nil
}

package storage

import (
	"math/rand"

	"github.com/catinello/base62"
)

func GenerateUniqKey() (string, error) {
	return base62.Encode(rand.Int()), nil
}

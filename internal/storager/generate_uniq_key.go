package storager

import (
	"github.com/google/uuid"
)

func GenerateUniqKey() (string, error) {
	return uuid.NewString(), nil
}

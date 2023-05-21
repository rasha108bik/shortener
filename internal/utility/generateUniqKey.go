package utility

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jxskiss/base62"
)

// GenerateUniqKey generate unique key
func GenerateUniqKey() string {
	id, _ := uuid.NewRandom()
	date := time.Now()
	str := fmt.Sprintf("%s:%s", id.String(), date.String())
	hash := md5.Sum([]byte(str))
	encodedHash := base62.Encode([]byte(hash[:]))
	return string(encodedHash[:6])
}

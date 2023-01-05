package storage

type Storager interface {
	StoreURL(originalURL string, shortURL string) error
	GetURLShortID(id string) (string, error)
	GetURLsShort() map[string]string
}

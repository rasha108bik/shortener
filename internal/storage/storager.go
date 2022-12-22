package storage

type Storager interface {
	StoreURL(data string) (string, error)
	GetURLShortID(id string) (string, error)
	GetURLsShort() map[string]string
}

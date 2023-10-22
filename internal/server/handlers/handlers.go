package handlers

import (
	"net/http"

	"github.com/rs/zerolog"

	"github.com/rasha108bik/tiny_url/config"
	"github.com/rasha108bik/tiny_url/internal/storager"
)

// Handlers all API public methods
// interface.
type Handlers interface {
	ErrorHandler(w http.ResponseWriter, r *http.Request)
	CreateShortLink(w http.ResponseWriter, r *http.Request)
	CreateShorten(w http.ResponseWriter, r *http.Request)
	GetOriginalURL(w http.ResponseWriter, r *http.Request)
	Stats(w http.ResponseWriter, r *http.Request)
	FetchURLs(w http.ResponseWriter, r *http.Request)
	Ping(w http.ResponseWriter, r *http.Request)
	ShortenBatch(w http.ResponseWriter, r *http.Request)
	DeleteURLs(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	log     *zerolog.Logger
	cfg     *config.Config
	storage storager.Storager
}

// NewHandler returns a newly initialized handler objects that implements the Handlers
// interface.
func NewHandler(
	log *zerolog.Logger,
	cfg *config.Config,
	storage storager.Storager,
) *handler {
	return &handler{
		log:     log,
		cfg:     cfg,
		storage: storage,
	}
}

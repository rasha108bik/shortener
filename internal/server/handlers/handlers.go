package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/rasha108bik/tiny_url/config"
	"github.com/rasha108bik/tiny_url/internal/storage"
	"github.com/rasha108bik/tiny_url/internal/storage/postgres"
)

type Handlers interface {
	CreateShorten(w http.ResponseWriter, r *http.Request)
	CreateShortLink(w http.ResponseWriter, r *http.Request)
	GetOriginalURL(w http.ResponseWriter, r *http.Request)
	GetOriginalURLs(w http.ResponseWriter, r *http.Request)
	ErrorHandler(w http.ResponseWriter, r *http.Request)
	Ping(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	cfg         *config.Config
	memDB       storage.Storager
	fileStorage storage.Storager
	pg          postgres.Postgres
}

func NewHandler(
	cfg *config.Config,
	memDB storage.Storager,
	fileStorage storage.Storager,
	pg postgres.Postgres,
) *handler {
	return &handler{
		cfg:         cfg,
		memDB:       memDB,
		fileStorage: fileStorage,
		pg:          pg,
	}
}

func (h *handler) ErrorHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "wrong method", http.StatusBadRequest)
}

func (h *handler) CreateShortLink(w http.ResponseWriter, r *http.Request) {
	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	originalURL := string(resBody)
	shortURL, err := storage.GenerateUniqKey()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.pg.StoreURL(originalURL, shortURL)
	if err != nil {
		log.Printf("pg.StoreURL: %v\n", err)
	}

	err = h.memDB.StoreURL(originalURL, shortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.fileStorage.StoreURL(originalURL, shortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(h.cfg.BaseURL + "/" + shortURL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handler) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "id")
	if shortURL == "" {
		http.Error(w, "id emtpy", http.StatusBadRequest)
		return
	}

	originalURL, err := h.pg.GetOriginalURLByShortURL(shortURL)
	if err != nil {
		log.Printf("pg.GetOriginalURLByShortURL: %v\n", originalURL)

		originalURL, err = h.memDB.GetURLShortID(shortURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}

func (h *handler) CreateShorten(w http.ResponseWriter, r *http.Request) {
	m := ReqCreateShorten{}
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL, err := storage.GenerateUniqKey()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.pg.StoreURL(m.URL, shortURL)
	if err != nil {
		log.Printf("pg.StoreURL: %v\n", err)
	}

	err = h.memDB.StoreURL(m.URL, shortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.fileStorage.StoreURL(m.URL, shortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respRCS := RespReqCreateShorten{Result: h.cfg.BaseURL + "/" + shortURL}
	response, err := json.Marshal(respRCS)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handler) GetOriginalURLs(w http.ResponseWriter, r *http.Request) {
	mapURLs, err := h.pg.GetAllURLs()
	if err != nil {
		log.Printf("pg.mapURLs: %v\n and data urls: %v\n", err, mapURLs)
	}

	mapURLs = h.memDB.GetURLsShort()
	if len(mapURLs) == 0 {
		http.Error(w, errors.New("urls is empty").Error(), http.StatusNoContent)
		return
	}

	mapData := mapperGetOriginalURLs(mapURLs, h.cfg.BaseURL)
	res, err := json.Marshal(mapData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func mapperGetOriginalURLs(data map[string]string, baseURL string) []RespGetOriginalURLs {
	res := make([]RespGetOriginalURLs, 0)
	for k, v := range data {
		res = append(res, RespGetOriginalURLs{
			ShortURL:    baseURL + "/" + k,
			OriginalURL: v,
		})
	}
	return res
}

func (h *handler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := h.pg.Ping(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

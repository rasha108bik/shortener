package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/rasha108bik/tiny_url/config"
	appErr "github.com/rasha108bik/tiny_url/internal/errors"
	"github.com/rasha108bik/tiny_url/internal/storage"
)

type Handlers interface {
	ErrorHandler(w http.ResponseWriter, r *http.Request)
	CreateShortLink(w http.ResponseWriter, r *http.Request)
	CreateShorten(w http.ResponseWriter, r *http.Request)
	GetOriginalURL(w http.ResponseWriter, r *http.Request)
	FetchURLs(w http.ResponseWriter, r *http.Request)
	Ping(w http.ResponseWriter, r *http.Request)
	ShortenBatch(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	cfg         *config.Config
	memDB       storage.Storager
	fileStorage storage.Storager
	pg          storage.Storager
	pgcon       bool
}

func NewHandler(
	cfg *config.Config,
	memDB storage.Storager,
	fileStorage storage.Storager,
	pg storage.Storager,
	pgcon bool,
) *handler {
	return &handler{
		cfg:         cfg,
		memDB:       memDB,
		fileStorage: fileStorage,
		pg:          pg,
		pgcon:       pgcon,
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

	if h.pgcon {
		shrURL, err := h.pg.GetShortURLByOriginalURL(originalURL)
		if err != nil {
			if err == sql.ErrNoRows {
				err = h.pg.StoreURL(originalURL, shortURL)
				if err != nil {
					log.Printf("pg.StoreURL: %v\n", err)
				}
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		if shrURL != "" {
			respCreateShortLink(w, http.StatusConflict, h.cfg.BaseURL, shrURL)
			return
		}
	}

	shrURL, err := h.memDB.GetShortURLByOriginalURL(originalURL)
	if err != nil {
		if errors.Is(err, appErr.ErrOriginalURLExist) {
			respCreateShortLink(w, http.StatusConflict, h.cfg.BaseURL, shrURL)
			return
		}
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

	respCreateShortLink(w, http.StatusCreated, h.cfg.BaseURL, shortURL)
}

func respCreateShortLink(w http.ResponseWriter, statusCode int, baseURL string, shortURL string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(baseURL + "/" + shortURL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handler) CreateShorten(w http.ResponseWriter, r *http.Request) {
	m := ReqCreateShorten{}
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("ReqCreateShorten: %v", m)
	shortURL, err := storage.GenerateUniqKey()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.pgcon {
		shrURL, err := h.pg.GetShortURLByOriginalURL(m.URL)
		if err != nil {
			if err == sql.ErrNoRows {
				err = h.pg.StoreURL(m.URL, shortURL)
				if err != nil {
					log.Printf("pg.StoreURL: %v\n", err)
				}
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		if shrURL != "" {
			respCreateShorten(w, http.StatusConflict, h.cfg.BaseURL, shrURL)
			return
		}
	}

	shrURL, err := h.memDB.GetShortURLByOriginalURL(m.URL)
	if err != nil {
		if errors.Is(err, appErr.ErrOriginalURLExist) {
			respCreateShorten(w, http.StatusConflict, h.cfg.BaseURL, shrURL)
			return
		}
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

	respCreateShorten(w, http.StatusCreated, h.cfg.BaseURL, shortURL)
}

func respCreateShorten(w http.ResponseWriter, statusCode int, baseURL string, shortURL string) {
	respRCS := RespReqCreateShorten{Result: baseURL + "/" + shortURL}
	log.Printf("respRCS: %v", respRCS)

	response, err := json.Marshal(respRCS)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(response)
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

	if h.pgcon {
		originalURL, err := h.pg.GetOriginalURLByShortURL(shortURL)
		if err != nil {
			log.Printf("pg.GetOriginalURLByShortURL: %v\n", originalURL)
		}
	}

	originalURL, err := h.memDB.GetOriginalURLByShortURL(shortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}

func (h *handler) FetchURLs(w http.ResponseWriter, r *http.Request) {
	mapURLs, _ := h.memDB.GetAllURLs()
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

func (h *handler) ShortenBatch(w http.ResponseWriter, r *http.Request) {
	m := []ReqShortenBatch{}
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respShortenBatch := []RespShortenBatch{}
	for _, v := range m {
		shortURL, err := storage.GenerateUniqKey()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if h.pgcon {
			err = h.pg.StoreURL(v.OriginalURL, shortURL)
			if err != nil {
				log.Printf("pg.StoreURL: %v\n", err)
			}
		}

		err = h.memDB.StoreURL(v.OriginalURL, shortURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.fileStorage.StoreURL(v.OriginalURL, shortURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		respShortenBatch = append(respShortenBatch, RespShortenBatch{
			CorrelationID: v.CorrelationID,
			ShortURL:      h.cfg.BaseURL + "/" + shortURL,
		})
	}

	res, err := json.Marshal(respShortenBatch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

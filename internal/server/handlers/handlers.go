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
	appErr "github.com/rasha108bik/tiny_url/internal/errors"
	"github.com/rasha108bik/tiny_url/internal/storager"
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
	cfg *config.Config
	str storager.Storager
}

func NewHandler(
	cfg *config.Config,
	str storager.Storager,
) *handler {
	return &handler{
		cfg: cfg,
		str: str,
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
	shortURL, err := storager.GenerateUniqKey()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shrURL, err := h.str.GetShortURLByOriginalURL(originalURL)
	if err != nil {
		if errors.Is(err, appErr.ErrOriginalURLExist) {
			respCreateShortLink(w, http.StatusConflict, h.cfg.BaseURL, shrURL)
			return
		}
	}

	err = h.str.StoreURL(originalURL, shortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	shortURL, err := storager.GenerateUniqKey()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shrURL, err := h.str.GetShortURLByOriginalURL(m.URL)
	if err != nil {
		if errors.Is(err, appErr.ErrOriginalURLExist) {
			respCreateShorten(w, http.StatusConflict, h.cfg.BaseURL, shrURL)
			return
		}
	}

	err = h.str.StoreURL(m.URL, shortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	originalURL, err := h.str.GetOriginalURLByShortURL(shortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}

func (h *handler) FetchURLs(w http.ResponseWriter, r *http.Request) {
	mapURLs, _ := h.str.GetAllURLs()
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

	err := h.str.Ping(ctx)
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
		shortURL, err := storager.GenerateUniqKey()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		shrURL, err := h.str.GetShortURLByOriginalURL(v.OriginalURL)
		if err != nil {
			if errors.Is(err, appErr.ErrOriginalURLExist) {
				respCreateShorten(w, http.StatusConflict, h.cfg.BaseURL, shrURL)
				return
			}
		}

		err = h.str.StoreURL(v.OriginalURL, shortURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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

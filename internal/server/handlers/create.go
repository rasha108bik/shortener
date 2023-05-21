package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	appErr "github.com/rasha108bik/tiny_url/internal/errors"
	"github.com/rasha108bik/tiny_url/internal/server/handlers/models"
	"github.com/rasha108bik/tiny_url/internal/utility"
)

// CreateShortLink create short link and save in DB
func (h *handler) CreateShortLink(w http.ResponseWriter, r *http.Request) {
	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	originalURL := string(resBody)
	shortURL := utility.GenerateUniqKey()

	shrURL, err := h.storage.GetShortURLByOriginalURL(r.Context(), originalURL)
	if err != nil {
		if errors.Is(err, appErr.ErrOriginalURLExist) {
			respCreateShortLink(w, http.StatusConflict, h.cfg.BaseURL, shrURL)
			return
		}
	}

	err = h.storage.StoreURL(r.Context(), originalURL, shortURL)
	if err != nil {
		h.log.Error().Err(err).Msg("StoreURL")
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

// CreateShorten create short link and save in DB with ReqCreateShorten model
func (h *handler) CreateShorten(w http.ResponseWriter, r *http.Request) {
	m := models.ReqCreateShorten{}
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("ReqCreateShorten: %v", m)
	shortURL := utility.GenerateUniqKey()

	shrURL, err := h.storage.GetShortURLByOriginalURL(r.Context(), m.URL)
	if err != nil {
		if errors.Is(err, appErr.ErrOriginalURLExist) {
			respCreateShorten(w, http.StatusConflict, h.cfg.BaseURL, shrURL)
			return
		}
	}

	err = h.storage.StoreURL(r.Context(), m.URL, shortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respCreateShorten(w, http.StatusCreated, h.cfg.BaseURL, shortURL)
}

func respCreateShorten(w http.ResponseWriter, statusCode int, baseURL string, shortURL string) {
	respRCS := models.RespReqCreateShorten{Result: baseURL + "/" + shortURL}
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

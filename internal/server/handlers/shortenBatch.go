package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	appErr "github.com/rasha108bik/tiny_url/internal/errors"
	"github.com/rasha108bik/tiny_url/internal/server/handlers/models"
	"github.com/rasha108bik/tiny_url/internal/utility"
)

// ShortenBatch save URLs which include ReqShortenBatch model
func (h *handler) ShortenBatch(w http.ResponseWriter, r *http.Request) {
	m := []models.ReqShortenBatch{}
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respShortenBatch := []models.RespShortenBatch{}
	for _, v := range m {
		shortURL := utility.GenerateUniqKey()

		shrURL, err := h.storage.GetShortURLByOriginalURL(r.Context(), v.OriginalURL)
		if err != nil {
			if errors.Is(err, appErr.ErrOriginalURLExist) {
				respCreateShorten(w, http.StatusConflict, h.cfg.BaseURL, shrURL)
				return
			}
		}

		err = h.storage.StoreURL(r.Context(), v.OriginalURL, shortURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		respShortenBatch = append(respShortenBatch, models.RespShortenBatch{
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

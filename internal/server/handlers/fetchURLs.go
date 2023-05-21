package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rasha108bik/tiny_url/internal/server/handlers/models"
)

// FetchURLs get all URLs from DB
func (h *handler) FetchURLs(w http.ResponseWriter, r *http.Request) {
	mapURLs, _ := h.storage.GetAllURLs(r.Context())
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

func mapperGetOriginalURLs(data map[string]string, baseURL string) []models.RespGetOriginalURLs {
	res := make([]models.RespGetOriginalURLs, 0)
	for k, v := range data {
		res = append(res, models.RespGetOriginalURLs{
			ShortURL:    baseURL + "/" + k,
			OriginalURL: v,
		})
	}
	return res
}

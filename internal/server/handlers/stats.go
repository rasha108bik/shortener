package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rasha108bik/tiny_url/internal/server/handlers/models"
)

// Stats get count short urls and users.
func (h *handler) Stats(w http.ResponseWriter, r *http.Request) {
	cntShortURLs, err := h.storage.GetCountShortURLAndUsers(r.Context())
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, errors.New("short urls is empty").Error(), http.StatusNoContent)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	mapData := models.ResponseStats{
		URLs:  cntShortURLs,
		Users: cntShortURLs,
	}

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

package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	appErr "github.com/rasha108bik/tiny_url/internal/errors"
)

func (h *handler) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "id")
	if shortURL == "" {
		http.Error(w, "id emtpy", http.StatusBadRequest)
		return
	}

	originalURL, err := h.storage.GetOriginalURLByShortURL(r.Context(), shortURL)
	if err != nil {
		if err == appErr.ErrURLDeleted {
			w.WriteHeader(http.StatusGone)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}

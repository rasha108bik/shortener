package handlers

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/rasha108bik/tiny_url/internal/storage"
)

type Handlers interface {
	CreateShortLink(w http.ResponseWriter, r *http.Request)
	GetOriginalURL(w http.ResponseWriter, r *http.Request)
	ErrorHandler(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	storage storage.Storage
}

func NewHandler(storage storage.Storage) *handler {
	return &handler{
		storage: storage,
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

	res, err := h.storage.StoreURL(string(resBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("http://127.0.0.1:8080/" + res))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handler) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id emtpy", http.StatusBadRequest)
		return
	}

	url, err := h.storage.GetURLShortID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

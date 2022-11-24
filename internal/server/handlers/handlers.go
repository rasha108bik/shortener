package handlers

import (
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/rasha108bik/tiny_url/internal/storage"
)

type Handlers interface {
	PostHandler() http.HandlerFunc
	GetHandler() http.HandlerFunc
}

type handler struct {
	storage storage.Storage
}

func New(storage storage.Storage) *handler {
	return &handler{
		storage: storage,
	}
}

func URLErrorHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "wrong method", http.StatusBadRequest)
}

func (h *handler) PostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		res, err := h.storage.Save(string(resBody))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		makeResponse(w, "application/json", []byte("http://127.0.0.1:8080/"+res), http.StatusCreated)
	}
}

func (h *handler) GetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "id emtpy", http.StatusBadRequest)
			return
		}

		url, err := h.storage.Get(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func makeResponse(w http.ResponseWriter, contenType string, body []byte, statusCode int) {
	w.Header().Set("Content-Type", contenType)
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		log.Println(err)
	}
}

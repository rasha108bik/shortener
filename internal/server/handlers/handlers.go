package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/rasha108bik/tiny_url/config"
	"github.com/rasha108bik/tiny_url/internal/storage"
)

type Handlers interface {
	CreateShorten(w http.ResponseWriter, r *http.Request)
	CreateShortLink(w http.ResponseWriter, r *http.Request)
	GetOriginalURL(w http.ResponseWriter, r *http.Request)
	ErrorHandler(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	cfg         *config.Config
	db          storage.Storager
	fileStorage storage.Storager
}

func NewHandler(
	cfg *config.Config,
	db storage.Storager,
	fileStorage storage.Storager,
) *handler {
	return &handler{
		cfg:         cfg,
		db:          db,
		fileStorage: fileStorage,
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

	res, err := h.db.StoreURL(string(resBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.fileStorage.StoreURL(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(h.cfg.BaseURL + "/" + res))
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

	url, err := h.db.GetURLShortID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *handler) CreateShorten(w http.ResponseWriter, r *http.Request) {
	m := ReqCreateShorten{}
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newURL, err := h.db.StoreURL(m.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.fileStorage.StoreURL(newURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respRCS := RespReqCreateShorten{Result: h.cfg.BaseURL + "/" + newURL}
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

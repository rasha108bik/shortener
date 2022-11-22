package handler

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"

	"github.com/rasha108bik/short_url/internal/app/server"
	"github.com/rasha108bik/short_url/internal/app/storage"
)

type Handler struct {
	// logger  *zap.Logger
	router  *chi.Mux
	storage storage.Storage
}

func NewHandler(
	// logger *zap.Logger,
	r *chi.Mux,
	db storage.Storage,
) server.HandlerV1 {
	return &Handler{
		// logger:  logger,
		router:  r,
		storage: db,
	}
}

func (h *Handler) PostShortURL(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			// h.logger.Sugar().Infof("failed read body", "body: ", string(b))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		data := string(b)
		_, err = url.Parse(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		db := h.storage
		dataEnc, err := db.Save(data)
		if err != nil {
			// h.logger.Sugar().Infof("failed save in db: %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		makeResponse(w, "application/json", []byte("http://127.0.0.1:8080/"+dataEnc), http.StatusCreated)
	}
}

func makeResponse(w http.ResponseWriter, contenType string, body []byte, statusCode int) {
	w.Header().Set("content-type", contenType)
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		log.Println(err)
	}
}

func (h *Handler) GetOriginalURL(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "id is empty", http.StatusBadRequest)
			return
		}
		fmt.Printf("chi id param: %s \n", id)
		// h.logger.Sugar().Infof("id: %s", id)

		url, err := h.storage.Get(id)
		if err != nil {
			fmt.Printf("url param: %s \n", url)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// fmt.Printf("url param: %s \n", url)

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/rasha108bik/tiny_url/internal/middleware"
	"github.com/rasha108bik/tiny_url/internal/server"
)

func NewRouter(s *server.Server) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.GzipHandle)
	r.Use(middleware.GzipRequest)
	r.Use(middleware.SetUserCookie)
	r.MethodNotAllowed(s.Handlers.ErrorHandler)
	r.Get("/ping", s.Handlers.Ping)
	r.Get("/{id}", s.Handlers.GetOriginalURL)
	r.Get("/api/user/urls", s.Handlers.FetchURLs)
	r.Post("/api/shorten", s.Handlers.CreateShorten)
	r.Post("/", s.Handlers.CreateShortLink)

	return r
}

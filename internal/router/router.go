package router

import (
	"github.com/go-chi/chi/v5"
	middlewareChi "github.com/go-chi/chi/v5/middleware"

	"github.com/rasha108bik/tiny_url/internal/middleware"
	"github.com/rasha108bik/tiny_url/internal/server/handlers"
)

// NewRouter returns a newly *chi.Muxo objects that registery pattern and middleware
func NewRouter(s handlers.Handlers) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.GzipHandle)
	r.Use(middleware.GzipRequest)
	r.Use(middleware.SetUserCookie)
	r.Mount("/debug", middlewareChi.Profiler())

	r.MethodNotAllowed(s.ErrorHandler)
	r.Get("/ping", s.Ping)
	r.Get("/{id}", s.GetOriginalURL)
	r.Get("/api/user/urls", s.FetchURLs)
	r.Get("/api/internal/stats", s.Stats)
	r.Post("/api/shorten", s.CreateShorten)
	r.Post("/", s.CreateShortLink)
	r.Post("/api/shorten/batch", s.ShortenBatch)
	r.Delete("/api/user/urls", s.DeleteURLs)

	return r
}

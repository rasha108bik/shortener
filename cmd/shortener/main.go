package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	handler "github.com/rasha108bik/short_url/internal/app/handlers"
	"github.com/rasha108bik/short_url/internal/app/storage"
)

func main() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// logger, err := zap.NewProduction()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer logger.Sync() // nolint

	db := storage.New()
	h := handler.NewHandler(r, db)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	r.Route("/", func(r chi.Router) {
		r.Post("/", h.PostShortURL(ctx))
		r.Get("/{id}", h.GetOriginalURL(ctx))
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
	}()

	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/rasha108bik/tiny_url/config"
	"github.com/rasha108bik/tiny_url/internal/storager"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func BenchmarkCreateShortLink(b *testing.B) {
	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", cfg)

	str, err := storager.NewStorager(&cfg)
	if err != nil {
		log.Printf("pgDB.New: %v\n", err)
	}
	defer str.Close()

	log := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)
	handler := NewHandler(&log, &cfg, str)

	var shortenURL string
	var originalURL string

	b.Run("save", func(b *testing.B) {
		originalURL = "http://jqymby.biz/wruxoh/1"

		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(originalURL))
		w := httptest.NewRecorder()
		h := http.HandlerFunc(handler.CreateShortLink)
		h(w, request)
		result := w.Result()

		userResult, _ := io.ReadAll(result.Body)
		result.Body.Close()

		shortenURL = string(userResult)
	})

	b.Run("get", func(b *testing.B) {
		uri, err := url.Parse(shortenURL)
		require.NoError(b, err)

		request := httptest.NewRequest(http.MethodGet, "/"+uri.Path[1:], nil)
		w := httptest.NewRecorder()
		h := http.HandlerFunc(handler.GetOriginalURL)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", uri.Path[1:])
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
		h(w, request)
	})

	b.Run("fetch_urls", func(b *testing.B) {
		request := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
		w := httptest.NewRecorder()
		h := http.HandlerFunc(handler.FetchURLs)
		h(w, request)
	})

	b.Run("save shorten", func(b *testing.B) {
		reqBody, err := json.Marshal(map[string]string{
			"url": "http://fsdkfkldshfjs.ru/1",
		})
		require.NoError(b, err)

		request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()
		h := http.HandlerFunc(handler.CreateShorten)
		h(w, request)
	})
}

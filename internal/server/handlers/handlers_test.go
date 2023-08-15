package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rasha108bik/tiny_url/config"
	"github.com/rasha108bik/tiny_url/internal/server/handlers/models"
	"github.com/rasha108bik/tiny_url/internal/storager"
)

func TestHandlers(t *testing.T) {
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

	t.Run("save", func(t *testing.T) {
		originalURL = "http://jqymby.biz/wruxoh/eii7bbkvbz4oj"

		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(originalURL))
		w := httptest.NewRecorder()
		h := http.HandlerFunc(handler.CreateShortLink)
		h(w, request)
		result := w.Result()

		assert.Equal(t, http.StatusCreated, result.StatusCode)
		assert.Equal(t, "application/json", result.Header.Get("Content-Type"))

		userResult, err := io.ReadAll(result.Body)
		require.NoError(t, err)
		err = result.Body.Close()
		require.NoError(t, err)

		shortenURL = string(userResult)

		// проверяем URL на валидность
		_, urlParseErr := url.Parse(shortenURL)
		assert.NoErrorf(t, urlParseErr, "cannot parsee URL: %s ", shortenURL, err)
		fmt.Println("shortenURL:  ", shortenURL)
	})

	t.Run("get", func(t *testing.T) {
		uri, err := url.Parse(shortenURL)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodGet, "/"+uri.Path[1:], nil)
		w := httptest.NewRecorder()
		h := http.HandlerFunc(handler.GetOriginalURL)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", uri.Path[1:])
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

		h(w, request)
		result := w.Result()
		err = result.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusTemporaryRedirect, result.StatusCode)
		assert.Equalf(t, originalURL, result.Header.Get("Location"),
			"Несоответствие URL полученного в заголовке Location ожидаемому",
		)
	})

	t.Run("fetch_urls", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
		w := httptest.NewRecorder()
		h := http.HandlerFunc(handler.FetchURLs)

		h(w, request)
		result := w.Result()
		err = result.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, result.StatusCode)

		m := []models.RespGetOriginalURLs{}
		err = json.NewDecoder(result.Body).Decode(&m)
		require.NoError(t, err)

		expectedBody := []models.RespGetOriginalURLs{
			{
				ShortURL:    shortenURL,
				OriginalURL: originalURL,
			},
		}

		assert.Equalf(t, expectedBody, m,
			"Данные в теле ответа не соответствуют ожидаемым",
		)
	})

	t.Run("save shorten", func(t *testing.T) {
		reqBody, err := json.Marshal(map[string]string{
			"url": "http://fsdkfkldshfjs.ru/test",
		})
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()
		h := http.HandlerFunc(handler.CreateShorten)
		h(w, request)
		result := w.Result()

		err = result.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, result.StatusCode)
		assert.Equal(t, "application/json", result.Header.Get("Content-Type"))

		m := models.RespReqCreateShorten{}
		err = json.NewDecoder(result.Body).Decode(&m)
		require.NoError(t, err)

		// проверяем URL на валидность
		_, urlParseErr := url.Parse(m.Result)
		assert.NoErrorf(t, urlParseErr, "cannot parsee URL: %s ", m.Result, err)
	})

	t.Run("get short urls and users", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/api/internal/stats", nil)
		w := httptest.NewRecorder()
		h := http.HandlerFunc(handler.GetOriginalURL)
		h(w, request)
		result := w.Result()

		err = result.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, result.StatusCode)
	})
}

func TestHandlersStatusConflict(t *testing.T) {
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

	t.Run("save", func(t *testing.T) {
		originalURL = "http://jqymby.biz/wruxoh/eii7bbkvbz411"

		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(originalURL))
		w := httptest.NewRecorder()
		h := http.HandlerFunc(handler.CreateShortLink)
		h(w, request)
		result := w.Result()

		assert.Equal(t, http.StatusCreated, result.StatusCode)
		assert.Equal(t, "application/json", result.Header.Get("Content-Type"))

		userResult, err := io.ReadAll(result.Body)
		require.NoError(t, err)
		err = result.Body.Close()
		require.NoError(t, err)

		shortenURL = string(userResult)

		// проверяем URL на валидность
		_, urlParseErr := url.Parse(shortenURL)
		assert.NoErrorf(t, urlParseErr, "cannot parsee URL: %s ", shortenURL, err)
		fmt.Println("shortenURL:  ", shortenURL)
	})

	t.Run("save status_conflict", func(t *testing.T) {
		originalURL = "http://jqymby.biz/wruxoh/eii7bbkvbz411"

		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(originalURL))
		w := httptest.NewRecorder()
		h := http.HandlerFunc(handler.CreateShortLink)
		h(w, request)
		result := w.Result()

		assert.Equal(t, http.StatusConflict, result.StatusCode)
		assert.Equal(t, "application/json", result.Header.Get("Content-Type"))

		userResult, err := io.ReadAll(result.Body)
		require.NoError(t, err)
		err = result.Body.Close()
		require.NoError(t, err)

		shortenURL = string(userResult)

		// проверяем URL на валидность
		_, urlParseErr := url.Parse(shortenURL)
		assert.NoErrorf(t, urlParseErr, "cannot parsee URL: %s ", shortenURL, err)
		fmt.Println("shortenURL:  ", shortenURL)
	})
}

func TestHandlersBatchRequest(t *testing.T) {
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

	requestData := []models.ReqShortenBatch{
		{
			CorrelationID: "1",
			OriginalURL:   "http://fsdkfkldshfjs.ru/test1",
		},
		{
			CorrelationID: "2",
			OriginalURL:   "http://fsdkfkldshfjs.ru/test2",
		},
		{
			CorrelationID: "3",
			OriginalURL:   "http://okdak1f046v.biz/njt7g7efm49f",
		},
		{
			CorrelationID: "4",
			OriginalURL:   "http://l1syjj.biz",
		},
	}
	// correlations between originalURLs and shortURLs
	correlations := make(map[string]string)

	t.Run("shorten_batch", func(t *testing.T) {
		reqBody, err := json.Marshal(requestData)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()
		h := http.HandlerFunc(handler.ShortenBatch)
		h(w, request)
		result := w.Result()

		err = result.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, result.StatusCode)
		assert.Equal(t, "application/json", result.Header.Get("Content-Type"))

		m := []models.RespShortenBatch{}
		err = json.NewDecoder(result.Body).Decode(&m)
		require.NoError(t, err)

		allCorrelationsFound := true
		for _, respPair := range m {
			var originalURL string
			for _, reqPair := range requestData {
				originalURL = reqPair.OriginalURL
				break
			}

			found := assert.NotEmptyf(t, originalURL, "not found original_url by correlation ID: %s", respPair.CorrelationID)
			if !found {
				allCorrelationsFound = false
			}

			correlations[respPair.ShortURL] = originalURL
		}

		if !allCorrelationsFound {
			dump := dumpRequest(request, true)
			jsonBody, _ := json.Marshal(requestData)
			t.Logf("Оригинальный запрос:\n\n%s\n\nТело запроса:\n\n%s", dump, jsonBody)
		}
	})
}

func dumpRequest(req *http.Request, body bool) (dump []byte) {
	if req != nil {
		dump, _ = httputil.DumpRequest(req, body)
	}
	return
}

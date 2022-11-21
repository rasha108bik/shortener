package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/rasha108bik/tiny_url/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlers(t *testing.T) {
	db := storage.New()
	handler := New(db)

	var shortenURL string
	var originalURL string

	t.Run("save	", func(t *testing.T) {
		originalURL = "http://jqymby.biz/wruxoh/eii7bbkvbz4oj"

		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(originalURL))
		w := httptest.NewRecorder()
		h := http.HandlerFunc(handler.PostHandler())
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
	})

	t.Run("get", func(t *testing.T) {
		uri, err := url.Parse(shortenURL)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodGet, "/"+uri.Path[1:], nil)
		w := httptest.NewRecorder()
		h := http.HandlerFunc(handler.GetHandler())

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
}

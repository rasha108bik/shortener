package server

import (
	"context"
	"net/http"
)

type HandlerV1 interface {
	PostShortURL(ctx context.Context) http.HandlerFunc
	GetOriginalURL(ctx context.Context) http.HandlerFunc
}

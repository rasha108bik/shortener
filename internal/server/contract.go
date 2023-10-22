package server

import "net/http"

type httpRoute interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type grpcRoute interface {
}

package server

import (
	"net/http"

	"github.com/rasha108bik/tiny_url/internal/router"
	"github.com/rasha108bik/tiny_url/internal/server/handlers"
)

func NewServer(
	h handlers.Handlers,
	serverAddress string,
) http.Server {
	r := router.NewRouter(h)
	return http.Server{Addr: serverAddress, Handler: r}
}

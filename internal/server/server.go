package server

import (
	"github.com/rasha108bik/tiny_url/internal/server/handlers"
)

type Server struct {
	Handlers handlers.Handlers
}

func NewServer(h handlers.Handlers) *Server {
	return &Server{
		Handlers: h,
	}
}

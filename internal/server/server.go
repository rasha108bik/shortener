package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rasha108bik/tiny_url/internal/router"
	"github.com/rasha108bik/tiny_url/internal/server/handlers"
	"golang.org/x/crypto/acme/autocert"
)

type Server struct {
	serv        http.Server
	enableHTTPS string
}

// NewServer returns a newly initialized http.Server objects
func (s Server) NewServer(
	h handlers.Handlers,
	serverAddress,
	enableHTTPS string,
) *Server {
	s.enableHTTPS = enableHTTPS
	r := router.NewRouter(h)
	s.buildHttpServer(r, serverAddress)

	return &s
}

func (s Server) buildHttpServer(
	r *chi.Mux,
	serverAddress string,
) *Server {
	if s.enableHTTPS != "" {
		manager := &autocert.Manager{
			// директория для хранения сертификатов
			Cache: autocert.DirCache("cache-dir"),
			// функция, принимающая Terms of Service издателя сертификатов
			Prompt: autocert.AcceptTOS,
			// перечень доменов, для которых будут поддерживаться сертификаты
			HostPolicy: autocert.HostWhitelist("mysite.ru", "www.mysite.ru"),
		}
		// конструируем сервер с поддержкой TLS
		s.serv = http.Server{
			Addr:    ":443",
			Handler: r,
			// для TLS-конфигурации используем менеджер сертификатов
			TLSConfig: manager.TLSConfig(),
		}
	}
	s.serv = http.Server{Addr: serverAddress, Handler: r}
	return &s
}

func (s Server) Start() error {
	var err error

	if s.enableHTTPS != "" {
		err = s.serv.ListenAndServeTLS("", "")
		if err != nil {
			return err
		}
	}

	err = s.serv.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

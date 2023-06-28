package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rasha108bik/tiny_url/internal/router"
	"github.com/rasha108bik/tiny_url/internal/server/handlers"
	"golang.org/x/crypto/acme/autocert"
)

type server struct {
	serv        http.Server
	enableHTTPS string
}

// NewServer returns a newly initialized http.Server objects
func NewServer(
	h handlers.Handlers,
	serverAddress,
	enableHTTPS string,
) *server {
	r := router.NewRouter(h)
	return &server{
		serv:        buildHTTPServer(r, serverAddress, enableHTTPS),
		enableHTTPS: enableHTTPS,
	}
}

func buildHTTPServer(
	r *chi.Mux,
	serverAddress,
	enableHTTPS string,
) http.Server {
	if enableHTTPS != "" {
		manager := &autocert.Manager{
			// директория для хранения сертификатов
			Cache: autocert.DirCache("cache-dir"),
			// функция, принимающая Terms of Service издателя сертификатов
			Prompt: autocert.AcceptTOS,
			// перечень доменов, для которых будут поддерживаться сертификаты
			HostPolicy: autocert.HostWhitelist("mysite.ru", "www.mysite.ru"),
		}
		// конструируем сервер с поддержкой TLS
		return http.Server{
			Addr:    ":443",
			Handler: r,
			// для TLS-конфигурации используем менеджер сертификатов
			TLSConfig: manager.TLSConfig(),
		}
	}

	return http.Server{Addr: serverAddress, Handler: r}
}

// Start is running server
func (s *server) Start() error {
	var err error

	if s.enableHTTPS != "" {
		err = s.serv.ListenAndServeTLS("", "")
		if err != nil {
			return err
		}
	} else {
		err = s.serv.ListenAndServe()
		if err != nil {
			return err
		}
	}

	return nil
}

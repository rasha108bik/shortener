package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/rasha108bik/tiny_url/internal/router"
	"github.com/rasha108bik/tiny_url/internal/server/handlers"
	"github.com/rs/zerolog"
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
	enableHTTPS,
	hostWhitelist string,
) *server {
	r := router.NewRouter(h)
	return &server{
		serv:        buildHTTPServer(r, serverAddress, enableHTTPS, hostWhitelist),
		enableHTTPS: enableHTTPS,
	}
}

func buildHTTPServer(
	r *chi.Mux,
	serverAddress,
	enableHTTPS,
	hostWhitelist string,
) http.Server {
	if enableHTTPS != "" {
		manager := &autocert.Manager{
			// директория для хранения сертификатов
			Cache: autocert.DirCache("cache-dir"),
			// функция, принимающая Terms of Service издателя сертификатов
			Prompt: autocert.AcceptTOS,
			// перечень доменов, для которых будут поддерживаться сертификаты
			HostPolicy: autocert.HostWhitelist(hostWhitelist),
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
func (s *server) Start(
	log *zerolog.Logger,
) error {
	var err error

	idleConnsClosed := make(chan struct{})
	sigs := make(chan os.Signal, 3)

	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		// read from interrupt channel
		<-sigs
		// recieved the os.interrupt signal, start the graceful shutdown procedure
		if err := s.serv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

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

	// wait for the graceful shutdown procedure to complete
	<-idleConnsClosed
	fmt.Println("Server Shutdown gracefully")

	return nil
}

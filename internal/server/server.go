package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rasha108bik/tiny_url/api/tinyurl/generated"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
)

type server struct {
	serv        http.Server
	enableHTTPS string
}

// NewServer returns a newly initialized http.Server objects
func NewServer(
	_ context.Context,
	httpRoute httpRoute,
	grpcRoute grpcRoute,
	serverAddress,
	enableHTTPS string,
) *server {
	return &server{
		serv:        buildHTTPServer(httpRoute, serverAddress, enableHTTPS),
		enableHTTPS: enableHTTPS,
	}
}

func buildHTTPServer(
	r httpRoute,
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

	lis, err := netListenerRun(log, s.enableHTTPS, s.serv.Addr)
	if err != nil {
		return err
	}

	// grpc server run
	grpcServer := buildGRPCServer()
	log.Info().Msg("grpcServer.Serve failed")
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			log.Error().Err(err).Msg("grpcServer server failed")
		}
	}()

	// wait for the graceful shutdown procedure to complete
	<-idleConnsClosed
	fmt.Println("Server Shutdown gracefully")

	return nil
}

func netListenerRun(log *zerolog.Logger, enableHTTPS, addr string) (net.Listener, error) {
	if enableHTTPS != "" {
		lis, err := tls.Listen("tcp", addr, nil)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to listen")
			return nil, err
		}

		return lis, nil

		// err = s.serv.ListenAndServeTLS("", "")
		// if err != nil {
		// 	return err
		// }
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
		return nil, err
	}

	return lis, nil

	// err = s.serv.ListenAndServe()
	// if err != nil {
	// 	return err
	// }
}

func buildGRPCServer() *grpc.Server {
	grpcServ := generated.UnimplementedApiServiceServer{}
	baseGrpcServer := grpc.NewServer()
	generated.RegisterApiServiceServer(baseGrpcServer, &grpcServ)

	return baseGrpcServer
}

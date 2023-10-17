package server

import (
	"context"
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

	// grpc server run
	grpcServer := buildGRPCServer()
	go func() {
		lis, err := net.Listen("tcp", ":8084")
		if err != nil {
			log.Error().Err(err).Msg("GRPC listen server failed")
		}

		if err = grpcServer.Serve(lis); err != nil {
			log.Error().Err(err).Msg("grpcServer server failed")
		}
	}()

	err = netListenerRun(log, s.enableHTTPS, s.serv.Addr, &s.serv)
	if err != nil {
		return err
	}

	// wait for the graceful shutdown procedure to complete
	<-idleConnsClosed
	fmt.Println("Server Shutdown gracefully")

	return nil
}

func netListenerRun(log *zerolog.Logger, enableHTTPS, addr string, serv *http.Server) error {
	if enableHTTPS != "" {
		err := serv.ListenAndServeTLS("", "")
		if err != nil {
			return err
		}
	}

	err := serv.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func buildGRPCServer() *grpc.Server {
	grpcServ := generated.UnimplementedApiServiceServer{}
	baseGrpcServer := grpc.NewServer()
	generated.RegisterApiServiceServer(baseGrpcServer, &grpcServ)

	return baseGrpcServer
}

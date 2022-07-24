package grpc

import (
	"net"
	"nexdata/pkg/config"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
)

var (
	customFunc grpc_recovery.RecoveryHandlerFunc
)

type Handler interface {
	ServeGRPC(*grpc.Server)
}

type Server struct {
	grpc    *grpc.Server
	Addr    string
	Handler Handler

	done chan struct{}
	eg   errgroup.Group
}

func NewServer(handler Handler) *Server {
	// Define customfunc to handle panic
	customFunc = func(p interface{}) (err error) {
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	}

	// Shared options for the logger, with a custom gRPC code to log level function.
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(customFunc),
	}

	// Create a server. Recovery handlers should typically be last in the chain so that other middleware
	// (e.g. logging) can operate on the recovered state instead of being directly affected by any panic
	srv := grpc.NewServer(
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(opts...),
		),
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(opts...),
		),
	)

	return &Server{
		grpc:    srv,
		Addr:    config.Source.GRPC.Host,
		Handler: handler,
		done:    make(chan struct{}),
	}
}

func (srv *Server) listenAndServe() error {
	lis, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}

	srv.Handler.ServeGRPC(srv.grpc)

	return srv.grpc.Serve(lis)
}

// Start Start HTTP server
func (srv *Server) Start() error {
	go srv.graceful()

	srv.eg.Go(func() error {
		log.Info().Msgf("Starting gRPC server [%s]", config.Source.GRPC.Host)
		return srv.listenAndServe()
	})

	if err := srv.eg.Wait(); err != nil {
		return err
	}

	<-srv.done

	return nil
}

// graceful Graceful shutdown
func (srv *Server) graceful() {
	sigint := make(chan os.Signal, 1)

	// interrupt signal sent from terminal
	signal.Notify(sigint, os.Interrupt)

	// sigterm signal sent from kubernetes
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(sigint)

	<-sigint

	log.Info().Msg("received an interrupt signal, shut down gRPC server.")
	// We received an interrupt signal, shut down.
	srv.grpc.GracefulStop()

	close(srv.done)
}

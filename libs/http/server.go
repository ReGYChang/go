package http

import (
	"context"
	"net/http"
	"nexdata/pkg/config"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"golang.org/x/sync/errgroup"
)

// Server HTTP server object
type Server struct {
	*http.Server

	done chan struct{}
	eg   errgroup.Group
}

// NewServer Create a new HTTP server object
func NewServer(handler http.Handler) *Server {
	server := &Server{
		Server: &http.Server{
			Addr:    config.Source.HTTP.Host,
			Handler: handler,
		},
		done: make(chan struct{}),
	}

	return server
}

// Start Start HTTP server
func (srv *Server) Start() error {
	go srv.graceful()

	srv.eg.Go(func() error {
		log.Info().Msg("Starting server ... !!")
		return srv.Server.ListenAndServe()
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

	log.Info().Msg("received an interrupt signal, shut down the server.")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// We received an interrupt signal, shut down.
	if err := srv.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		log.Printf("HTTP server Shutdown: %v", err)
	}

	close(srv.done)
}

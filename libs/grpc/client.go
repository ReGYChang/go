package grpc

import (
	"nexdata/pkg/config"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/rs/zerolog/log"
)

var (
	mux sync.RWMutex

	connMap = make(map[string]*grpc.ClientConn, 0)
)

// Dial creates a client connection to the given target.
func Dial(target string) *grpc.ClientConn {
	opts := []grpc_retry.CallOption{
		grpc_retry.WithMax(config.Destination.Retry.Max),                                    // 重試最大上限
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(config.Destination.Retry.Interval)), // 重試間隔時間
	}

	clientConn, err := grpc.Dial(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)),
	)
	if err != nil {
		log.Fatal().Err(err).Msgf("did not connect[%s]", target)
	}

	go graceful(clientConn)

	return clientConn
}

// GetConn Initialized and exposed a singleton object of gRPC client.
func GetConn(target string) *grpc.ClientConn {
	mux.RLock()
	if cc, ok := connMap[target]; ok {
		mux.RUnlock()
		return cc
	}

	mux.RUnlock()

	mux.Lock()
	defer mux.Unlock()
	if cc, ok := connMap[target]; ok {
		return cc
	}

	cc := Dial(target)
	connMap[target] = cc

	return cc
}

// graceful Graceful shutdown
func graceful(conn *grpc.ClientConn) {
	sigint := make(chan os.Signal, 1)

	// interrupt signal sent from terminal
	signal.Notify(sigint, os.Interrupt)

	// sigterm signal sent from kubernetes
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(sigint)

	<-sigint

	log.Info().Msg("received an interrupt signal, close the gRPC client.")
	conn.Close()
}

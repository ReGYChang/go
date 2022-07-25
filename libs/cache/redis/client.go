package redis

import (
	"context"
	"nexdata/pkg/config"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
)

var (
	// Used to create a singleton object of Elasticsearch client.
	// Initialized and exposed through GetClient().
	client *redis.Client

	// Used to execute client creation procedure only once.
	once sync.Once
)

func GetClient() (*redis.Client, error) {
	var err error

	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:        config.Redis.Host,
			DB:          config.Redis.DB,
			Username:    config.Redis.Username,
			Password:    config.Redis.Password,
			ReadTimeout: config.Redis.ReadTimeout,
		})

		go graceful()

		_, err := client.Ping(context.Background()).Result()
		if err != nil {
			log.Fatal().Msgf("error creating the redis client: %v", err)
		}

		return
	})

	return client, err
}

// graceful shutdown
func graceful() {
	sigint := make(chan os.Signal, 1)

	// interrupt signal sent from terminal
	signal.Notify(sigint, os.Interrupt)

	// sigterm signal sent from kubernetes
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(sigint)

	<-sigint

	log.Info().Msg("received an interrupt signal, disconnect the MongoDB client.")
	if err := client.Close(); err != nil {
		log.Panic().Err(err).Msg("")
	}
}

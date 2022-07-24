package mongodb

import (
	"context"
	"fmt"
	"nexdata/pkg/config"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	APP_NAME string = "nexconvot"
)

var (
	// Used to create a singleton object of MongoDB client.
	// Initialized and exposed through GetClient().
	client *mongo.Client

	// Used to execute client creation procedure only once.
	once sync.Once
)

// GetClient Initialized and exposed a singleton object of MongoDB client.
func GetClient() (*mongo.Client, error) {
	var err error

	once.Do(func() {
		// Set client options
		uri := fmt.Sprintf("mongodb://%s", config.Mongo.Host)
		credential := options.Credential{
			Username: config.Mongo.Username,
			Password: config.Mongo.Password,
		}
		opts := options.Client().
			ApplyURI(uri).
			SetAppName(APP_NAME).
			SetAuth(credential)

		// Connect to MongoDB
		client, err = mongo.Connect(context.Background(), opts)
		if err != nil {
			log.Err(err).Msg("")
			return
		}

		go graceful()

		// Check the connection
		err = client.Ping(context.Background(), readpref.Primary())
		if err != nil {
			log.Printf("opon mongodb fail: %v", err)
			return
		}
	})

	return client, err
}

// graceful Graceful shutdown
func graceful() {
	sigint := make(chan os.Signal, 1)

	// interrupt signal sent from terminal
	signal.Notify(sigint, os.Interrupt)

	// sigterm signal sent from kubernetes
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(sigint)

	<-sigint

	log.Info().Msg("received an interrupt signal, disconnect the MongoDB client.")
	if err := client.Disconnect(context.Background()); err != nil {
		log.Panic().Err(err).Msg("")
	}
}

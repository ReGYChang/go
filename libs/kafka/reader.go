package kafka

import (
	"context"
	"nexdata/pkg/config"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

// Reader Kafka reader object
type Reader struct {
	reader *kafka.Reader

	done    chan struct{}
	Message chan kafka.Message
}

// NewReader Make a new reader
func NewReader() *Reader {
	r := &Reader{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  config.Source.Kafka.Brokers,
			GroupID:  config.Source.Kafka.ConsumerGroup,
			Topic:    config.Source.Kafka.Topic,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
		done:    make(chan struct{}),
		Message: make(chan kafka.Message),
	}

	go r.graceful()

	return r
}

// Start Start Kafka reader
func (r *Reader) Start() error {
	log.Info().Msg("start consuming ... !!")
	for {
		m, err := r.reader.ReadMessage(context.Background())
		if err != nil {
			log.Err(err).Msg("")
			break
		}

		r.Message <- m
	}

	return nil
}

func (r *Reader) Done() <-chan struct{} {
	if r.done == nil {
		r.done = make(chan struct{})
	}

	return r.done
}

// graceful Graceful close
func (r *Reader) graceful() {
	sigint := make(chan os.Signal, 1)

	// interrupt signal sent from terminal
	signal.Notify(sigint, os.Interrupt)

	// sigterm signal sent from kubernetes
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(sigint)

	<-sigint

	log.Info().Msg("received an interrupt signal, shut down kafka reader.")
	// We received an interrupt signal, shut down.
	if err := r.reader.Close(); err != nil {
		log.Fatal().Err(err).Msg("failed to close reader")
	}

	close(r.done)
}

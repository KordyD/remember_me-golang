package main

import (
	"context"
	"fmt"
	"github.com/kordyd/remember_me-golang/services/event_ingestion/pkg/models"
	"github.com/segmentio/kafka-go"
	"log"
)

func main() {
	server := New()
	server.Save()
}

type Server struct {
	Consumer *kafka.Reader
}

type EventSaver interface {
	Save(event models.Event) error
}

func New() *Server {
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "events",
	})
	return &Server{
		Consumer: consumer,
	}
}

func (s *Server) Save() {
	for {
		message, err := s.Consumer.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("Failed to read message: %v", err)
		}
		fmt.Printf("Received message: value=%s, partition=%d, offset=%d\n",
			string(message.Value), message.Partition, message.Offset)
	}
}

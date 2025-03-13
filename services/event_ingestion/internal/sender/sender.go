package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kordyd/remember_me-golang/event_ingestion/pkg/models"
	"github.com/segmentio/kafka-go"
	"log/slog"
)

type Sender struct {
	log      *slog.Logger
	producer *kafka.Writer
}

func New(log *slog.Logger) *Sender {
	producer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "events",
		Balancer: &kafka.RoundRobin{},
	}
	// TODO producer.Close()
	return &Sender{
		log:      log,
		producer: producer,
	}
}

func (s *Sender) SendEvent(event models.Event) error {
	const op = "sender.SendEvent"
	jsonEvent, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	err = s.producer.WriteMessages(context.TODO(),
		kafka.Message{
			Value: jsonEvent,
		},
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

package sender

import (
	"github.com/kordyd/remember_me-golang/event_ingestion/internal/models"
	"log/slog"
)

type Sender struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Sender {
	return &Sender{
		log: log,
	}
}

func (s *Sender) SendEvent(event models.Event) error {
	//TODO implement me
	s.log.Debug("SendEvent", "event", event.EventType)
	return nil
}

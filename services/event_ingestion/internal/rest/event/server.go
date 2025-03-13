package event

import (
	"encoding/json"
	"github.com/kordyd/remember_me-golang/services/event_ingestion/pkg/models"
	"io"
	"log/slog"
	"net/http"
)

type Sender interface {
	SendEvent(event models.Event) error
}

type ServerAPI struct {
	sender Sender
	log    *slog.Logger
}

func NewServer(log *slog.Logger, sender Sender) *ServerAPI {
	return &ServerAPI{
		sender: sender,
		log:    log,
	}
}

func (s ServerAPI) HandleSendEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			s.log.Error("failed to close body", "error", err)
		}
		return
	}(r.Body)

	if event.PageUrl == "" {
		http.Error(w, "missing page_url", http.StatusBadRequest)
		return
	}

	if event.UserId == "" {
		http.Error(w, "missing user_id", http.StatusBadRequest)
		return
	}

	if event.EventType == "" {
		http.Error(w, "missing event_type", http.StatusBadRequest)
		return
	}

	if event.Id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	if event.Timestamp.IsZero() {
		http.Error(w, "missing timestamp", http.StatusBadRequest)
		return
	}

	err = s.sender.SendEvent(event)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write([]byte(event.EventType))
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
}

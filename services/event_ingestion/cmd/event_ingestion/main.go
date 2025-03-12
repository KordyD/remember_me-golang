package main

import (
	"github.com/kordyd/remember_me-golang/event_ingestion/internal/app"
	"github.com/kordyd/remember_me-golang/event_ingestion/internal/config"
	"github.com/kordyd/remember_me-golang/event_ingestion/internal/logger"
	"github.com/kordyd/remember_me-golang/event_ingestion/internal/sender"
	"os"
	"os/signal"
	"syscall"
)

type EventProvider interface {
	GetEvents( /*opts*/ )
}

func main() {
	cfg := config.MustSetup()
	log := logger.MustSetup(cfg.Env)
	senderMock := sender.New(log)
	application := app.New(log, senderMock, cfg.Port)
	go func() {
		application.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	application.MustStop()
}

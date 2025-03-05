package main

import (
	"fmt"
	"github.com/kordyd/remember_me-golang/sso/internal/app"
	"github.com/kordyd/remember_me-golang/sso/internal/config"
	"github.com/kordyd/remember_me-golang/sso/internal/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
	application := app.New(log, cfg.GRPCPort, psqlInfo, cfg.TokenTTL)
	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	application.GRPCServer.Stop()
	log.Info("Gracefully stopped")
}

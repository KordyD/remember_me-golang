package app

import (
	"fmt"
	grpcapp "github.com/kordyd/remember_me-golang/sso/internal/app/grpc"
	"github.com/kordyd/remember_me-golang/sso/internal/logger"
	"github.com/kordyd/remember_me-golang/sso/internal/services/auth"
	"github.com/kordyd/remember_me-golang/sso/internal/storage/postgres"
	"log/slog"
	"time"
)

type App struct {
	db         *postgres.Postgres
	gRPCServer *grpcapp.App
	log        *slog.Logger
}

func New(log *slog.Logger,
	grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	db, err := postgres.New(storagePath)
	if err != nil {
		log.Error("failed to connect to database", logger.Err(err))
		panic(err)
	}
	authService := auth.New(log, db, db, tokenTTL)
	grpcApp := grpcapp.New(log, grpcPort, authService)
	return &App{
		db:         db,
		log:        log,
		gRPCServer: grpcApp,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		a.log.Error("failed to run app", "err", logger.Err(err))
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "App.Run"
	err := a.gRPCServer.Run()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() error {
	const op = "App.Stop"
	a.log.With(slog.String("op", op)).Info("stopping app")
	a.gRPCServer.Stop()
	err := a.db.Close()
	if err != nil {
		a.log.With(slog.String("op", op)).Error("failed to close database")
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

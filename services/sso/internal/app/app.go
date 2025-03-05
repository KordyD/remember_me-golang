package app

import (
	"github.com/kordyd/remember_me-golang/sso/internal/services/auth"
	"github.com/kordyd/remember_me-golang/sso/internal/storage/postgres"
	"log/slog"
	"time"
)
import grpcapp "github.com/kordyd/remember_me-golang/sso/internal/app/grpc"

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger,
	grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	db, err := postgres.New(storagePath)
	if err != nil {
		panic(err)
	}
	authService := auth.New(log, db, db, tokenTTL)
	grpcApp := grpcapp.New(log, grpcPort, authService)
	return &App{
		GRPCServer: grpcApp,
	}
}

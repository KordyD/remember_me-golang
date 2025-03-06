package grpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	authgrpc "github.com/kordyd/remember_me-golang/sso/internal/grpc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net"
)

type App struct {
	gRPCServer *grpc.Server
	log        *slog.Logger
	port       int
}

func New(log *slog.Logger, port int, authService authgrpc.Auth) *App {
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) error {
			log.Error("recovery from panic", "panic", p)
			return status.Errorf(codes.Internal, "server internal error")
		}),
	}
	loggerOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent),
	}
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, fields ...any) {
			log.Log(ctx, slog.Level(level), msg, fields...)
		}), loggerOpts...)))
	authgrpc.Register(grpcServer, authService)
	return &App{
		gRPCServer: grpcServer,
		log:        log,
		port:       port,
	}
}

func (a *App) Run() error {
	const op = "gRPCApp.Run"
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	a.log.Info("grpc server started", slog.String("addr", lis.Addr().String()))
	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "gRPCApp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()

	a.log.With(slog.String("op", op)).
		Info("gRPC server stopped gracefully", slog.Int("port", a.port))

}

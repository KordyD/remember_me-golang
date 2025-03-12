package app

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kordyd/remember_me-golang/event_ingestion/internal/rest/event"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	log    *slog.Logger
	Server *http.Server
}

func New(log *slog.Logger, sender event.Sender, port int) *App {
	eventServer := event.NewServer(log, sender)
	r := chi.NewRouter()
	middleware.DefaultLogger = loggerMiddleware(log)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Post("/send-event", eventServer.HandleSendEvent)

	return &App{
		log: log,
		Server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: r,
		},
	}
}

func loggerMiddleware(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var bodyBytes []byte
			if r.Body != nil {
				bodyBytes, _ = io.ReadAll(r.Body)
			}
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			log.Info("request received",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("protocol", r.Proto),
				slog.String("remote", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("body", string(bodyBytes)),
			)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()
			next.ServeHTTP(ww, r)
			log.Info("response send",
				slog.Int("status", ww.Status()),
				slog.String("latency", time.Since(start).String()),
			)
		})
	}
}

func (a *App) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}

func (a *App) run() error {
	const op = "app.run"
	a.log.Info(fmt.Sprintf("%s starting", op))
	err := a.Server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) MustStop() {
	if err := a.stop(); err != nil {
		panic(err)
	}
}

func (a *App) stop() error {
	const op = "app.stop"
	a.log.Info(fmt.Sprintf("%s starting", op))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := a.Server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	a.log.Info("Server shut down gracefully")
	return nil
}

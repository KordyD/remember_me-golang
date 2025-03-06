package logger

import (
	"github.com/kordyd/remember_me-golang/protos/gens/go/sso"
	"log/slog"
	"os"
	"strings"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger
	// TODO log to file or db
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					return sanitize(a)
				},
			}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					return sanitize(a)
				},
			}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					return sanitize(a)
				},
			}),
		)
	}

	return log
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func sanitize(a slog.Attr) slog.Attr {
	if a.Key == "grpc.request.content" {
		if req, ok := a.Value.Any().(*sso.LoginRequest); ok {
			maskedReq := &sso.LoginRequest{
				Email:    req.GetEmail(),
				Password: maskString(req.GetPassword()),
				AppId:    req.GetAppId(),
			}
			a.Value = slog.AnyValue(maskedReq)
		}
		if req, ok := a.Value.Any().(*sso.RegisterRequest); ok {
			maskedReq := &sso.RegisterRequest{
				Email:    req.GetEmail(),
				Password: maskString(req.GetPassword()),
			}
			a.Value = slog.AnyValue(maskedReq)
		}
	}
	return a
}

func maskString(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]) + strings.Repeat("*", len(s)-1)
}

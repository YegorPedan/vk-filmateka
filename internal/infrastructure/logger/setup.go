package logger

import (
	"io"
	"log/slog"
	"os"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

var log *slog.Logger = nil

// NewLogger use onlu after SetupLogger func
func NewLogger() *slog.Logger {
	return log
}

func SetupLogger(env string) *slog.Logger {
	if log != nil {
		return log
	}

	switch env {
	case EnvLocal:
		log = setupLocalLog(os.Stdout)
	case EnvDev:
		log = setupDevLog(os.Stdout)
	case EnvProd:
		log = setupProdLog(os.Stdout)
	}

	return log
}

func setupLocalLog(out io.Writer) *slog.Logger {
	jsonHandler := slog.NewJSONHandler(out, &slog.HandlerOptions{Level: slog.LevelDebug})
	return slog.New(jsonHandler)
}

func setupDevLog(out io.Writer) *slog.Logger {
	jsonHandler := slog.NewJSONHandler(out, &slog.HandlerOptions{Level: slog.LevelInfo})
	return slog.New(jsonHandler)
}

func setupProdLog(out io.Writer) *slog.Logger {
	jsonHandler := slog.NewJSONHandler(out, &slog.HandlerOptions{Level: slog.LevelError})
	return slog.New(jsonHandler)
}

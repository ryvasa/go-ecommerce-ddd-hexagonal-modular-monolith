package logger

import (
	"log/slog"
	"os"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/logger"
)

type SlogLogger struct {
	log *slog.Logger
}

func NewLogger() logger.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return &SlogLogger{
		log: slog.New(handler),
	}
}

func (l *SlogLogger) Info(msg string, fields ...any) {
	l.log.Info(msg, fields...)
}

func (l *SlogLogger) Error(msg string, fields ...any) {
	l.log.Error(msg, fields...)
}

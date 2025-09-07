package logging

import (
	"log/slog"
	"os"

	"github.com/umutcomlekci/automated-messaging-system/internal/config"
)

func NewLogger(loggerName string) *slog.Logger {
	logger := slog.New(
		slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		})).
		With("app_name", config.GetAppName()).
		With("logger", loggerName).
		With("release", config.GetRelease())
	return logger
}

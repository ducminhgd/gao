package log

import (
	"log/slog"
	"os"
	"strings"

	gormslog "github.com/onrik/gorm-slog"
)

var (
	// DefaultSlogLogger is a default slog.Logger instance for backward compatibility
	DefaultSlogLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
)

func ConvertSlogLevel(l string) slog.Level {
	switch strings.ToUpper(l) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

func NewGORMLogger(loglevel string) *gormslog.Logger {
	return gormslog.New(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     ConvertSlogLevel(loglevel),
	})))
}

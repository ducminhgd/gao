package log

import (
	"context"
	"log/slog"
	"os"
)

// slogLogger wraps slog.Logger to implement our Logger interface
type slogLogger struct {
	logger           *slog.Logger
	contextExtractor ContextExtractor
}

// newSlogLogger creates a new slog-based logger
func newSlogLogger(config Config) (Logger, error) {
	// Parse log level
	level := parseSlogLevel(config.Level)

	// Create handler options
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: config.EnableCaller,
	}

	// Choose handler based on format
	var handler slog.Handler
	if config.Format == FormatConsole {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	// Create base logger
	baseLogger := slog.New(handler)

	// Add initial fields
	var initialAttrs []any
	if config.ServiceName != "" {
		initialAttrs = append(initialAttrs, slog.String("service", config.ServiceName))
	}
	if config.ServiceVersion != "" {
		initialAttrs = append(initialAttrs, slog.String("version", config.ServiceVersion))
	}
	if config.Environment != "" {
		initialAttrs = append(initialAttrs, slog.String("environment", config.Environment))
	}

	// Add additional fields from config
	for _, field := range config.AdditionalFields {
		initialAttrs = append(initialAttrs, slog.Any(field.Key, field.Value))
	}

	if len(initialAttrs) > 0 {
		baseLogger = baseLogger.With(initialAttrs...)
	}

	return &slogLogger{
		logger:           baseLogger,
		contextExtractor: config.ContextExtractor,
	}, nil
}

// parseSlogLevel converts our level string to slog level
func parseSlogLevel(level string) slog.Level {
	switch level {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	case LevelFatal:
		// slog doesn't have a fatal level, use error + 4
		return slog.LevelError + 4
	default:
		return slog.LevelInfo
	}
}

// fieldsToSlogAttrs converts our Field slice to slog attributes
func fieldsToSlogAttrs(fields []Field) []any {
	attrs := make([]any, len(fields))
	for i, f := range fields {
		attrs[i] = slog.Any(f.Key, f.Value)
	}
	return attrs
}

func (l *slogLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, fieldsToSlogAttrs(fields)...)
}

func (l *slogLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, fieldsToSlogAttrs(fields)...)
}

func (l *slogLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, fieldsToSlogAttrs(fields)...)
}

func (l *slogLogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, fieldsToSlogAttrs(fields)...)
}

func (l *slogLogger) Fatal(msg string, fields ...Field) {
	// slog doesn't have Fatal, so we log at the highest level and exit
	l.logger.Log(context.Background(), parseSlogLevel(LevelFatal), msg, fieldsToSlogAttrs(fields)...)
	os.Exit(1)
}

func (l *slogLogger) WithFields(fields ...Field) Logger {
	return &slogLogger{
		logger:           l.logger.With(fieldsToSlogAttrs(fields)...),
		contextExtractor: l.contextExtractor,
	}
}

func (l *slogLogger) WithContext(ctx context.Context) Logger {
	if ctx == nil || l.contextExtractor == nil {
		return l
	}

	fields := l.contextExtractor(ctx)
	if len(fields) == 0 {
		return l
	}

	return l.WithFields(fields...)
}

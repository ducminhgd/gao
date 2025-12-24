package log

import (
	"context"
	"fmt"
)

const (
	KindZap  = "zap"
	KindSlog = "slog"
)

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelFatal = "fatal"
)

const (
	FormatJSON    = "json"
	FormatConsole = "console"
)

// Logger is the main interface for structured logging
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)

	// WithFields returns a new logger with preset fields
	WithFields(fields ...Field) Logger

	// WithContext extracts correlation IDs from context (trace_id, request_id, etc.)
	WithContext(ctx context.Context) Logger
}

// Field is a key-value pair for structured logging
type Field struct {
	Key   string
	Value any
}

var std Logger

func Debug(msg string, fields ...Field) {
	if std != nil {
		std.Debug(msg, fields...)
	}
}

func Info(msg string, fields ...Field) {
	if std != nil {
		std.Info(msg, fields...)
	}
}

func Warn(msg string, fields ...Field) {
	if std != nil {
		std.Warn(msg, fields...)
	}
}

func Error(msg string, fields ...Field) {
	if std != nil {
		std.Error(msg, fields...)
	}
}

func Fatal(msg string, fields ...Field) {
	if std != nil {
		std.Fatal(msg, fields...)
	}
}

func WithFields(fields ...Field) Logger {
	if std != nil {
		return std.WithFields(fields...)
	}
	return nil
}

func WithContext(ctx context.Context) Logger {
	if std != nil {
		return std.WithContext(ctx)
	}
	return nil
}

// ContextExtractor extracts fields from a context
type ContextExtractor func(ctx context.Context) []Field

// DefaultContextExtractor returns a context extractor that looks for common trace/request IDs
func DefaultContextExtractor() ContextExtractor {
	return func(ctx context.Context) []Field {
		var fields []Field

		// Extract trace_id if present
		if traceID := ctx.Value("trace_id"); traceID != nil {
			fields = append(fields, Field{Key: "trace_id", Value: traceID})
		}

		// Extract request_id if present
		if requestID := ctx.Value("request_id"); requestID != nil {
			fields = append(fields, Field{Key: "request_id", Value: requestID})
		}

		// Extract user_id if present
		if userID := ctx.Value("user_id"); userID != nil {
			fields = append(fields, Field{Key: "user_id", Value: userID})
		}

		return fields
	}
}

type Config struct {
	Kind   string // KindZap, KindSlog
	Level  string // LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal
	Format string // FormatJSON or FormatConsole

	Development      bool
	ServiceName      string
	ServiceVersion   string
	Environment      string
	EnableCaller     bool
	EnableStacktrace bool
	AdditionalFields []Field

	ContextExtractor ContextExtractor
}

// New creates a logger. First logger becomes the global default.
func New(config Config) (Logger, error) {
	if config.Kind == "" {
		config.Kind = KindZap
	}
	if config.Level == "" {
		config.Level = LevelInfo
	}
	if config.Format == "" {
		config.Format = FormatJSON
	}
	if config.ContextExtractor == nil {
		config.ContextExtractor = DefaultContextExtractor()
	}

	var log Logger
	var err error

	switch config.Kind {
	case KindZap:
		log, err = newZapLogger(config)
	case KindSlog:
		log, err = newSlogLogger(config)
	default:
		return nil, fmt.Errorf("unsupported logger kind: %s (supported: %s, %s)", config.Kind, KindZap, KindSlog)
	}

	if err != nil {
		return nil, err
	}

	if std == nil {
		std = log
	}

	return log, nil
}

func SetStd(l Logger) {
	std = l
}

// DefaultConfig returns the default production configuration
func DefaultConfig() Config {
	return Config{
		Kind:             KindZap,
		Level:            LevelInfo,
		Format:           FormatJSON,
		Development:      false,
		EnableCaller:     false,
		EnableStacktrace: true,
	}
}

// DevelopmentConfig returns a configuration suitable for development
func DevelopmentConfig() Config {
	return Config{
		Kind:             KindZap,
		Level:            LevelDebug,
		Format:           FormatConsole,
		Development:      true,
		EnableCaller:     false,
		EnableStacktrace: true,
	}
}

// SimpleConfig returns a minimal configuration
func SimpleConfig() Config {
	return Config{
		Kind:         KindZap,
		Level:        LevelInfo,
		Format:       FormatJSON,
		EnableCaller: false,
	}
}

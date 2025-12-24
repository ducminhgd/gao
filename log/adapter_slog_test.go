package log

import (
	"context"
	"log/slog"
	"testing"
)

func TestNewSlogLogger(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "default config",
			config: Config{
				Kind:  KindSlog,
				Level: LevelInfo,
			},
		},
		{
			name: "debug level",
			config: Config{
				Kind:  KindSlog,
				Level: LevelDebug,
			},
		},
		{
			name: "console format",
			config: Config{
				Kind:   KindSlog,
				Level:  LevelInfo,
				Format: FormatConsole,
			},
		},
		{
			name: "json format",
			config: Config{
				Kind:   KindSlog,
				Level:  LevelInfo,
				Format: FormatJSON,
			},
		},
		{
			name: "with caller enabled",
			config: Config{
				Kind:         KindSlog,
				Level:        LevelInfo,
				EnableCaller: true,
			},
		},
		{
			name: "with service metadata",
			config: Config{
				Kind:           KindSlog,
				Level:          LevelInfo,
				ServiceName:    "test-service",
				ServiceVersion: "1.0.0",
				Environment:    "test",
			},
		},
		{
			name: "with additional fields",
			config: Config{
				Kind:  KindSlog,
				Level: LevelInfo,
				AdditionalFields: []Field{
					{Key: "region", Value: "us-east-1"},
					{Key: "datacenter", Value: "dc1"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := newSlogLogger(tt.config)
			if err != nil {
				t.Fatalf("failed to create slog logger: %v", err)
			}
			if logger == nil {
				t.Fatal("expected logger, got nil")
			}

			// Test that we can call log methods without panic
			logger.Debug("debug message", Field{Key: "test", Value: "value"})
			logger.Info("info message", Field{Key: "test", Value: "value"})
			logger.Warn("warn message", Field{Key: "test", Value: "value"})
			logger.Error("error message", Field{Key: "test", Value: "value"})
		})
	}
}

func TestParseSlogLevel(t *testing.T) {
	tests := []struct {
		level    string
		expected slog.Level
	}{
		{LevelDebug, slog.LevelDebug},
		{LevelInfo, slog.LevelInfo},
		{LevelWarn, slog.LevelWarn},
		{LevelError, slog.LevelError},
		{LevelFatal, slog.LevelError + 4},
		{"unknown", slog.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			result := parseSlogLevel(tt.level)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSlogLoggerWithFields(t *testing.T) {
	cfg := Config{
		Kind:  KindSlog,
		Level: LevelDebug,
	}

	logger, err := newSlogLogger(cfg)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	// Test WithFields
	newLogger := logger.WithFields(
		Field{Key: "user_id", Value: "123"},
		Field{Key: "request_id", Value: "req-456"},
	)

	if newLogger == nil {
		t.Fatal("WithFields returned nil")
	}

	// Should be able to log with the new logger
	newLogger.Info("test message")
}

func TestSlogLoggerWithContext(t *testing.T) {
	cfg := Config{
		Kind:  KindSlog,
		Level: LevelDebug,
	}

	logger, err := newSlogLogger(cfg)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	// Test with nil context
	nilCtxLogger := logger.WithContext(nil)
	if nilCtxLogger != logger {
		t.Error("WithContext(nil) should return the same logger")
	}

	// Test with context containing values
	ctx := context.WithValue(context.Background(), "trace_id", "trace-123")
	ctx = context.WithValue(ctx, "request_id", "req-456")

	ctxLogger := logger.WithContext(ctx)
	if ctxLogger == nil {
		t.Fatal("WithContext returned nil")
	}

	// Should be able to log with the context logger
	ctxLogger.Info("test message with context")
}

func TestFieldsToSlogAttrs(t *testing.T) {
	fields := []Field{
		{Key: "string", Value: "test"},
		{Key: "int", Value: 42},
		{Key: "bool", Value: true},
		{Key: "float", Value: 3.14},
	}

	attrs := fieldsToSlogAttrs(fields)
	if len(attrs) != len(fields) {
		t.Errorf("expected %d attrs, got %d", len(fields), len(attrs))
	}
}

func TestSlogLoggerAllLevels(t *testing.T) {
	cfg := Config{
		Kind:  KindSlog,
		Level: LevelDebug,
	}

	logger, err := newSlogLogger(cfg)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	// Test all log levels (except Fatal which calls os.Exit)
	tests := []struct {
		name string
		fn   func()
	}{
		{
			name: "debug",
			fn: func() {
				logger.Debug("debug message", Field{Key: "level", Value: "debug"})
			},
		},
		{
			name: "info",
			fn: func() {
				logger.Info("info message", Field{Key: "level", Value: "info"})
			},
		},
		{
			name: "warn",
			fn: func() {
				logger.Warn("warn message", Field{Key: "level", Value: "warn"})
			},
		},
		{
			name: "error",
			fn: func() {
				logger.Error("error message", Field{Key: "level", Value: "error"})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			tt.fn()
		})
	}
}

func TestSlogLoggerTextHandler(t *testing.T) {
	cfg := Config{
		Kind:   KindSlog,
		Level:  LevelInfo,
		Format: FormatConsole,
	}

	logger, err := newSlogLogger(cfg)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	// Should create a text handler for console format
	logger.Info("text handler test", Field{Key: "format", Value: "console"})
}

func TestSlogLoggerJSONHandler(t *testing.T) {
	cfg := Config{
		Kind:   KindSlog,
		Level:  LevelInfo,
		Format: FormatJSON,
	}

	logger, err := newSlogLogger(cfg)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	// Should create a JSON handler for JSON format
	logger.Info("json handler test", Field{Key: "format", Value: "json"})
}

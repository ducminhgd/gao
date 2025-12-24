package log

import (
	"context"
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestNewZapLogger(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "default config",
			config: Config{
				Kind:  KindZap,
				Level: LevelInfo,
			},
		},
		{
			name: "debug level",
			config: Config{
				Kind:  KindZap,
				Level: LevelDebug,
			},
		},
		{
			name: "console format",
			config: Config{
				Kind:   KindZap,
				Level:  LevelInfo,
				Format: FormatConsole,
			},
		},
		{
			name: "json format",
			config: Config{
				Kind:   KindZap,
				Level:  LevelInfo,
				Format: FormatJSON,
			},
		},
		{
			name: "with caller enabled",
			config: Config{
				Kind:         KindZap,
				Level:        LevelInfo,
				EnableCaller: true,
			},
		},
		{
			name: "with stacktrace enabled",
			config: Config{
				Kind:             KindZap,
				Level:            LevelInfo,
				EnableStacktrace: true,
			},
		},
		{
			name: "development mode",
			config: Config{
				Kind:        KindZap,
				Level:       LevelDebug,
				Development: true,
			},
		},
		{
			name: "with service metadata",
			config: Config{
				Kind:           KindZap,
				Level:          LevelInfo,
				ServiceName:    "test-service",
				ServiceVersion: "1.0.0",
				Environment:    "test",
			},
		},
		{
			name: "with additional fields",
			config: Config{
				Kind:  KindZap,
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
			logger, err := newZapLogger(tt.config)
			if err != nil {
				t.Fatalf("failed to create zap logger: %v", err)
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

func TestParseZapLevel(t *testing.T) {
	tests := []struct {
		level    string
		expected zapcore.Level
	}{
		{LevelDebug, zapcore.DebugLevel},
		{LevelInfo, zapcore.InfoLevel},
		{LevelWarn, zapcore.WarnLevel},
		{LevelError, zapcore.ErrorLevel},
		{LevelFatal, zapcore.FatalLevel},
		{"unknown", zapcore.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			result := parseZapLevel(tt.level)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestZapLoggerWithFields(t *testing.T) {
	cfg := Config{
		Kind:  KindZap,
		Level: LevelDebug,
	}

	logger, err := newZapLogger(cfg)
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

func TestZapLoggerWithContext(t *testing.T) {
	cfg := Config{
		Kind:  KindZap,
		Level: LevelDebug,
	}

	logger, err := newZapLogger(cfg)
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

func TestFieldsToZap(t *testing.T) {
	fields := []Field{
		{Key: "string", Value: "test"},
		{Key: "int", Value: 42},
		{Key: "bool", Value: true},
		{Key: "float", Value: 3.14},
	}

	zapFields := fieldsToZap(fields)
	if len(zapFields) != len(fields) {
		t.Errorf("expected %d fields, got %d", len(fields), len(zapFields))
	}
}

func TestZapLoggerAllLevels(t *testing.T) {
	cfg := Config{
		Kind:  KindZap,
		Level: LevelDebug,
	}

	logger, err := newZapLogger(cfg)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	// Test all log levels
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

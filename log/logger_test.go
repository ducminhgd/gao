package log

import (
	"context"
	"testing"
)

func TestDefaultContextExtractor(t *testing.T) {
	extractor := DefaultContextExtractor()

	tests := []struct {
		name     string
		ctx      context.Context
		expected int
	}{
		{
			name:     "empty context",
			ctx:      context.Background(),
			expected: 0,
		},
		{
			name:     "context with trace_id",
			ctx:      context.WithValue(context.Background(), "trace_id", "trace-123"),
			expected: 1,
		},
		{
			name: "context with all fields",
			ctx: func() context.Context {
				ctx := context.WithValue(context.Background(), "trace_id", "trace-123")
				ctx = context.WithValue(ctx, "request_id", "req-456")
				ctx = context.WithValue(ctx, "user_id", "user-789")
				return ctx
			}(),
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := extractor(tt.ctx)
			if len(fields) != tt.expected {
				t.Errorf("expected %d fields, got %d", tt.expected, len(fields))
			}
		})
	}
}

func TestConfig(t *testing.T) {
	t.Run("DefaultConfig", func(t *testing.T) {
		cfg := DefaultConfig()
		if cfg.Kind != KindZap {
			t.Errorf("expected kind %s, got %s", KindZap, cfg.Kind)
		}
		if cfg.Level != LevelInfo {
			t.Errorf("expected level %s, got %s", LevelInfo, cfg.Level)
		}
		if cfg.Format != FormatJSON {
			t.Errorf("expected format %s, got %s", FormatJSON, cfg.Format)
		}
		if cfg.Development {
			t.Error("expected Development to be false")
		}
	})

	t.Run("DevelopmentConfig", func(t *testing.T) {
		cfg := DevelopmentConfig()
		if cfg.Kind != KindZap {
			t.Errorf("expected kind %s, got %s", KindZap, cfg.Kind)
		}
		if cfg.Level != LevelDebug {
			t.Errorf("expected level %s, got %s", LevelDebug, cfg.Level)
		}
		if cfg.Format != FormatConsole {
			t.Errorf("expected format %s, got %s", FormatConsole, cfg.Format)
		}
		if !cfg.Development {
			t.Error("expected Development to be true")
		}
	})

	t.Run("SimpleConfig", func(t *testing.T) {
		cfg := SimpleConfig()
		if cfg.Kind != KindZap {
			t.Errorf("expected kind %s, got %s", KindZap, cfg.Kind)
		}
		if cfg.Level != LevelInfo {
			t.Errorf("expected level %s, got %s", LevelInfo, cfg.Level)
		}
		if cfg.Format != FormatJSON {
			t.Errorf("expected format %s, got %s", FormatJSON, cfg.Format)
		}
	})
}

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		config    Config
		wantError bool
	}{
		{
			name: "valid zap logger",
			config: Config{
				Kind:  KindZap,
				Level: LevelInfo,
			},
			wantError: false,
		},
		{
			name: "valid slog logger",
			config: Config{
				Kind:  KindSlog,
				Level: LevelDebug,
			},
			wantError: false,
		},
		{
			name: "invalid kind",
			config: Config{
				Kind: "invalid",
			},
			wantError: true,
		},
		{
			name: "zap with custom fields",
			config: Config{
				Kind:           KindZap,
				Level:          LevelInfo,
				ServiceName:    "test-service",
				ServiceVersion: "1.0.0",
				Environment:    "test",
				AdditionalFields: []Field{
					{Key: "custom", Value: "value"},
				},
			},
			wantError: false,
		},
		{
			name: "slog with custom fields",
			config: Config{
				Kind:           KindSlog,
				Level:          LevelWarn,
				ServiceName:    "test-service",
				ServiceVersion: "1.0.0",
				AdditionalFields: []Field{
					{Key: "region", Value: "us-east-1"},
				},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.config)
			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if logger == nil {
					t.Error("expected logger, got nil")
				}
			}
		})
	}
}

func TestSetStd(t *testing.T) {
	// Save original std
	originalStd := std
	defer func() {
		std = originalStd
	}()

	cfg := SimpleConfig()
	cfg.Kind = KindSlog
	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	SetStd(logger)
	if std != logger {
		t.Error("SetStd did not set the global logger")
	}
}

func TestGlobalFunctions(t *testing.T) {
	// Save original std
	originalStd := std
	defer func() {
		std = originalStd
	}()

	// Create a test logger
	cfg := SimpleConfig()
	cfg.Kind = KindSlog
	cfg.Level = LevelDebug
	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	SetStd(logger)

	// Test that global functions don't panic
	t.Run("global functions", func(t *testing.T) {
		Debug("debug message", Field{Key: "test", Value: "value"})
		Info("info message", Field{Key: "test", Value: "value"})
		Warn("warn message", Field{Key: "test", Value: "value"})
		Error("error message", Field{Key: "test", Value: "value"})

		contextLogger := WithContext(context.Background())
		if contextLogger == nil {
			t.Error("WithContext returned nil")
		}

		fieldsLogger := WithFields(Field{Key: "test", Value: "value"})
		if fieldsLogger == nil {
			t.Error("WithFields returned nil")
		}
	})

	// Test with nil std
	t.Run("nil std", func(t *testing.T) {
		std = nil
		Debug("debug")
		Info("info")
		Warn("warn")
		Error("error")
		// These should not panic
	})
}

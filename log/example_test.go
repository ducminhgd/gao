package log_test

import (
	"context"

	"github.com/ducminhgd/gao/log"
)

// ExampleNew_zapLogger demonstrates creating a Zap logger
func ExampleNew_zapLogger() {
	logger, err := log.New(log.Config{
		Kind:   log.KindZap,
		Level:  log.LevelInfo,
		Format: log.FormatJSON,
	})
	if err != nil {
		panic(err)
	}

	logger.Info("Application started")
	logger.Info("User logged in", log.Field{Key: "user_id", Value: "12345"})
}

// ExampleNew_slogLogger demonstrates creating a slog logger
func ExampleNew_slogLogger() {
	logger, err := log.New(log.Config{
		Kind:   log.KindSlog,
		Level:  log.LevelInfo,
		Format: log.FormatJSON,
	})
	if err != nil {
		panic(err)
	}

	logger.Info("Application started")
	logger.Info("User logged in", log.Field{Key: "user_id", Value: "12345"})
}

// ExampleNew_withServiceMetadata demonstrates creating a logger with service metadata
func ExampleNew_withServiceMetadata() {
	logger, err := log.New(log.Config{
		Kind:           log.KindZap,
		Level:          log.LevelInfo,
		Format:         log.FormatJSON,
		ServiceName:    "my-api",
		ServiceVersion: "1.0.0",
		Environment:    "production",
	})
	if err != nil {
		panic(err)
	}

	logger.Info("Service initialized")
}

// ExampleNew_development demonstrates creating a development logger
func ExampleNew_development() {
	logger, err := log.New(log.DevelopmentConfig())
	if err != nil {
		panic(err)
	}

	logger.Debug("Debug information")
	logger.Info("Application started")
}

// ExampleLogger_WithFields demonstrates using WithFields
func ExampleLogger_WithFields() {
	logger, _ := log.New(log.SimpleConfig())

	// Create a logger with preset fields
	requestLogger := logger.WithFields(
		log.Field{Key: "request_id", Value: "req-123"},
		log.Field{Key: "user_id", Value: "user-456"},
	)

	// All logs from requestLogger will include request_id and user_id
	requestLogger.Info("Processing request")
	requestLogger.Info("Request completed")
}

// ExampleLogger_WithContext demonstrates using WithContext
func ExampleLogger_WithContext() {
	logger, _ := log.New(log.SimpleConfig())

	// Create a context with trace information
	ctx := context.WithValue(context.Background(), "trace_id", "trace-xyz")
	ctx = context.WithValue(ctx, "request_id", "req-123")

	// Extract fields from context
	contextLogger := logger.WithContext(ctx)
	contextLogger.Info("Processing with context")
}

// ExampleDefaultConfig demonstrates using default production config
func ExampleDefaultConfig() {
	logger, err := log.New(log.DefaultConfig())
	if err != nil {
		panic(err)
	}

	logger.Info("Production logger initialized")
}

// ExampleDevelopmentConfig demonstrates using development config
func ExampleDevelopmentConfig() {
	logger, err := log.New(log.DevelopmentConfig())
	if err != nil {
		panic(err)
	}

	logger.Debug("Development mode enabled")
}

// Example_globalLogger demonstrates using the global logger functions
func Example_globalLogger() {
	// Create and set the global logger
	_, err := log.New(log.Config{
		Kind:  log.KindZap,
		Level: log.LevelInfo,
	})
	if err != nil {
		panic(err)
	}

	// Use global functions
	log.Info("Application started")
	log.Info("User action", log.Field{Key: "action", Value: "login"})
	log.Error("An error occurred", log.Field{Key: "error", Value: "database connection failed"})
}

// Example_structuredLogging demonstrates structured logging with multiple fields
func Example_structuredLogging() {
	logger, _ := log.New(log.SimpleConfig())

	logger.Info("User registered",
		log.Field{Key: "user_id", Value: "user-123"},
		log.Field{Key: "email", Value: "user@example.com"},
		log.Field{Key: "registration_method", Value: "oauth"},
		log.Field{Key: "provider", Value: "google"},
	)

	logger.Error("Database query failed",
		log.Field{Key: "query", Value: "SELECT * FROM users WHERE id = ?"},
		log.Field{Key: "error", Value: "connection timeout"},
		log.Field{Key: "duration_ms", Value: 5000},
	)
}

// Example_customContextExtractor demonstrates using a custom context extractor
func Example_customContextExtractor() {
	customExtractor := func(ctx context.Context) []log.Field {
		var fields []log.Field

		// Extract custom fields from context
		if orgID := ctx.Value("org_id"); orgID != nil {
			fields = append(fields, log.Field{Key: "org_id", Value: orgID})
		}
		if tenantID := ctx.Value("tenant_id"); tenantID != nil {
			fields = append(fields, log.Field{Key: "tenant_id", Value: tenantID})
		}

		return fields
	}

	logger, _ := log.New(log.Config{
		Kind:             log.KindZap,
		Level:            log.LevelInfo,
		ContextExtractor: customExtractor,
	})

	ctx := context.WithValue(context.Background(), "org_id", "org-789")
	ctx = context.WithValue(ctx, "tenant_id", "tenant-456")

	logger.WithContext(ctx).Info("Multi-tenant operation")
}

// Example_differentLogLevels demonstrates using different log levels
func Example_differentLogLevels() {
	logger, _ := log.New(log.Config{
		Kind:  log.KindZap,
		Level: log.LevelDebug,
	})

	logger.Debug("Detailed debugging information", log.Field{Key: "variable", Value: "value"})
	logger.Info("Informational message")
	logger.Warn("Warning: resource usage high", log.Field{Key: "cpu_percent", Value: 85})
	logger.Error("Error processing request", log.Field{Key: "error", Value: "timeout"})
}

// Example_consoleFormat demonstrates using console format for development
func Example_consoleFormat() {
	logger, _ := log.New(log.Config{
		Kind:   log.KindZap,
		Level:  log.LevelInfo,
		Format: log.FormatConsole,
	})

	logger.Info("Console formatted output", log.Field{Key: "readable", Value: true})
}

// Example_jsonFormat demonstrates using JSON format for production
func Example_jsonFormat() {
	logger, _ := log.New(log.Config{
		Kind:   log.KindSlog,
		Level:  log.LevelInfo,
		Format: log.FormatJSON,
	})

	logger.Info("JSON formatted output", log.Field{Key: "structured", Value: true})
}

// Example_additionalFields demonstrates using additional fields in config
func Example_additionalFields() {
	logger, _ := log.New(log.Config{
		Kind:        log.KindZap,
		Level:       log.LevelInfo,
		ServiceName: "my-service",
		AdditionalFields: []log.Field{
			{Key: "region", Value: "us-east-1"},
			{Key: "availability_zone", Value: "us-east-1a"},
			{Key: "instance_id", Value: "i-1234567890abcdef0"},
		},
	})

	// All logs will include the additional fields
	logger.Info("Server started")
}

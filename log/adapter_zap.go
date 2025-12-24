package log

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger wraps zap.Logger to implement our Logger interface
type zapLogger struct {
	logger           *zap.Logger
	contextExtractor ContextExtractor
}

// newZapLogger creates a new zap-based logger
func newZapLogger(config Config) (Logger, error) {
	// Configure encoding
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Development mode adjustments
	if config.Development {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Choose encoder based on format
	var encoder zapcore.Encoder
	if config.Format == FormatConsole {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// Parse log level
	level := parseZapLevel(config.Level)

	// Create core
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		level,
	)

	// Build options
	opts := []zap.Option{}

	if config.EnableCaller {
		opts = append(opts, zap.AddCaller())
	}

	if config.EnableStacktrace {
		opts = append(opts, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	if config.Development {
		opts = append(opts, zap.Development())
	}

	// Create base logger
	baseLogger := zap.New(core, opts...)

	// Add initial fields
	var initialFields []zap.Field
	if config.ServiceName != "" {
		initialFields = append(initialFields, zap.String("service", config.ServiceName))
	}
	if config.ServiceVersion != "" {
		initialFields = append(initialFields, zap.String("version", config.ServiceVersion))
	}
	if config.Environment != "" {
		initialFields = append(initialFields, zap.String("environment", config.Environment))
	}

	// Add additional fields from config
	for _, field := range config.AdditionalFields {
		initialFields = append(initialFields, zap.Any(field.Key, field.Value))
	}

	if len(initialFields) > 0 {
		baseLogger = baseLogger.With(initialFields...)
	}

	return &zapLogger{
		logger:           baseLogger,
		contextExtractor: config.ContextExtractor,
	}, nil
}

// parseZapLevel converts our level string to zap level
func parseZapLevel(level string) zapcore.Level {
	switch level {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	case LevelFatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// fieldsToZap converts our Field slice to zap.Field slice
func fieldsToZap(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = zap.Any(f.Key, f.Value)
	}
	return zapFields
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, fieldsToZap(fields)...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, fieldsToZap(fields)...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, fieldsToZap(fields)...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, fieldsToZap(fields)...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, fieldsToZap(fields)...)
}

func (l *zapLogger) WithFields(fields ...Field) Logger {
	return &zapLogger{
		logger:           l.logger.With(fieldsToZap(fields)...),
		contextExtractor: l.contextExtractor,
	}
}

func (l *zapLogger) WithContext(ctx context.Context) Logger {
	if ctx == nil || l.contextExtractor == nil {
		return l
	}

	fields := l.contextExtractor(ctx)
	if len(fields) == 0 {
		return l
	}

	return l.WithFields(fields...)
}

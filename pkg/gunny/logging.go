package gunny

import (
	"context"
	"log/slog"
)

type contextKey string

const contextKeyLogger contextKey = "logger"

// Logger allows callers to specify their own logging interface.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// DefaultLogger returns Gunny's default [Logger] instance. By default, this
// simply uses [log/slog].
func DefaultLogger() Logger {
	return &SlogLogger{}
}

// NewContextWithLogger constructs a new [context.Context] with the given
// logger attached as a value.
func NewContextWithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKeyLogger, logger)
}

// LoggerFromContext attempts to get the [Logger] instance associated with the
// given context, if one exists. If no such logger exists, the result of
// [DefaultLogger] will be returned.
func LoggerFromContext(ctx context.Context) Logger {
	loggerValue := ctx.Value(contextKeyLogger)
	if loggerValue == nil {
		return DefaultLogger()
	}
	logger, ok := loggerValue.(Logger)
	if !ok {
		return DefaultLogger()
	}
	return logger
}

// SlogLogger is a trivial facade for [log/slog] that implements [Logger].
type SlogLogger struct{}

var _ Logger = (*SlogLogger)(nil)

// Debug implements Logger.
func (s *SlogLogger) Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

// Info implements Logger.
func (s *SlogLogger) Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

// Warn implements Logger.
func (s *SlogLogger) Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

// Error implements Logger.
func (s *SlogLogger) Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

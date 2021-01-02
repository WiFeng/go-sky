package log

import (
	"context"
)

type loggerKey struct{}

var activeLoggerKey = loggerKey{}

// ContextWithLogger function
func ContextWithLogger(ctx context.Context, logg Logger) context.Context {
	return context.WithValue(ctx, activeLoggerKey, logg)
}

// LoggerFromContext function
func LoggerFromContext(ctx context.Context) Logger {
	val := ctx.Value(activeLoggerKey)
	if logg, ok := val.(Logger); ok {
		return logg
	}
	return defaultLogger.With()
}

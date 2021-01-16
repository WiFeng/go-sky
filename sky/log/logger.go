package log

import (
	"context"

	"github.com/WiFeng/go-sky/sky/config"
	"github.com/opentracing/opentracing-go"
	jaegerclient "github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
)

var (
	// traceIDKey ...
	traceIDKey = "trace_id"
)

// Logger interface
type Logger interface {
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	DPanicw(msg string, keysAndValues ...interface{})
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Panic(args ...interface{})
	Panicf(template string, args ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Sync() error
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	With(args ...interface{}) Logger

	Log(keyvals ...interface{}) error
}

type logger struct {
	*zap.SugaredLogger
}

// default logger
var defaultLogger Logger

var (
	// TypeKey ...
	TypeKey = "file"
	// TypeValAccess ...
	TypeValAccess = "access.log"
	// TypeValRPC ...
	TypeValRPC = "rpc.log"
)

// Init ...
func Init(logConf config.Log) (Logger, error) {
	logger, err := NewLogger(logConf)
	if err != nil {
		return logger, err
	}

	SetDefaultLogger(logger)
	return logger, nil
}

// GetDefaultLogger return default logger
func GetDefaultLogger() Logger {
	return defaultLogger
}

// SetDefaultLogger set default logger
func SetDefaultLogger(logg Logger) {
	defaultLogger = logg
}

// NewLogger new Logger
func NewLogger(logConf config.Log) (Logger, error) {
	zapConf := config.NewZapConfig(logConf)
	zapOptions := []zap.Option{
		zap.AddCallerSkip(1),
	}
	zapLogger, err := zapConf.Build(zapOptions...)
	if err != nil {
		return nil, err
	}

	logger := logger{
		zapLogger.Sugar(),
	}
	return logger, nil
}

func (l logger) Log(keyvals ...interface{}) error {
	l.Info(keyvals...)
	return nil
}

func (l logger) With(args ...interface{}) Logger {
	sl := l.SugaredLogger.With(args...)

	logger := logger{
		sl,
	}

	return logger
}

// GetTraceID Get trace id from the context.
func GetTraceID(ctx context.Context) string {
	var traceID string
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		spanContext := span.Context()
		jeagerSpanContext, ok := spanContext.(jaegerclient.SpanContext)
		if ok {
			traceID = jeagerSpanContext.TraceID().String()
		}
	}

	return traceID
}

// BuildLogger ...
func BuildLogger(ctx context.Context) context.Context {
	newLogg := GetDefaultLogger()

	traceID := GetTraceID(ctx)
	if traceID != "" {
		newLogg = newLogg.With(traceIDKey, traceID)
	} else {
		newLogg = newLogg.With(traceIDKey, traceID)
	}

	newCtx := ContextWithLogger(ctx, newLogg)
	return newCtx
}

// Info ...
func Info(ctx context.Context, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Info(args...)
}

// Infof ...
func Infof(ctx context.Context, template string, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Infof(template, args...)
}

// Infow ...
func Infow(ctx context.Context, msg string, keysAndValues ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Infow(msg, keysAndValues...)
}

// Debug ...
func Debug(ctx context.Context, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Debug(args...)
}

// Debugf ...
func Debugf(ctx context.Context, template string, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Debugf(template, args...)
}

// Debugw ...
func Debugw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Debugw(msg, keysAndValues...)
}

// Warn ...
func Warn(ctx context.Context, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Warn(args...)
}

// Warnf ...
func Warnf(ctx context.Context, template string, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Warnf(template, args...)
}

// Warnw ...
func Warnw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Warnw(msg, keysAndValues...)
}

// Error ...
func Error(ctx context.Context, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Error(args...)
}

// Errorf ...
func Errorf(ctx context.Context, template string, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Errorf(template, args...)
}

// Errorw ...
func Errorw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Errorw(msg, keysAndValues...)
}

// DPanic ...
func DPanic(ctx context.Context, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.DPanic(args...)
}

// DPanicf ...
func DPanicf(ctx context.Context, template string, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.DPanicf(template, args...)
}

// DPanicw ...
func DPanicw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.DPanicw(msg, keysAndValues...)
}

// Panic ...
func Panic(ctx context.Context, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Panic(args...)
}

// Panicf ...
func Panicf(ctx context.Context, template string, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Panicf(template, args...)
}

// Panicw ...
func Panicw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Panicw(msg, keysAndValues...)
}

// Fatal ...
func Fatal(ctx context.Context, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Fatal(args...)
}

// Fatalf ...
func Fatalf(ctx context.Context, template string, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Fatalf(template, args...)
}

// Fatalw ...
func Fatalw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Fatalw(msg, keysAndValues...)
}

// Sync ...
func Sync(ctx context.Context) {
	logg := LoggerFromContext(ctx)
	logg.Sync()
}

// With ...
func With(ctx context.Context, args ...interface{}) Logger {
	logg := LoggerFromContext(ctx)
	newlogg := logg.With(args...)
	return newlogg
}

// Log ...
func Log(ctx context.Context, keyvals ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Log(keyvals...)
}

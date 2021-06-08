package log

import (
	"context"
	"fmt"
	"os"

	"github.com/WiFeng/go-sky/config"
	"github.com/WiFeng/go-sky/helper"
	skyprome "github.com/WiFeng/go-sky/metrics/prometheus"
	"github.com/opentracing/opentracing-go"
	jaegerclient "github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	jacklog "gopkg.in/natefinch/lumberjack.v2"
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
	// TypeValRedis ...
	TypeValRedis = "redis.log"
)

// Init ...
func Init(ctx context.Context, serviceName string, cfg config.Log) (logger Logger, err error) {
	logger, err = NewLogger(cfg)
	if err != nil {
		fmt.Println("Init logger error. ", err)
		os.Exit(1)
		return
	}

	SetDefaultLogger(logger)
	helper.AddDeferFunc(func() {
		logger.Sync()
	})

	return
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
func NewLogger(cfg config.Log) (Logger, error) {
	var filename string
	var enableStdout bool

	filename = "./logs/runtime.log"
	if cfg.OutputPath != "" {
		filename = cfg.OutputPath
	}

	enableStdout = false
	if cfg.Development {
		enableStdout = true
	}

	return newLogger(cfg, filename, enableStdout)
}

func newLogger(cfg config.Log, filename string, enableStdout bool) (Logger, error) {
	writeSyncer := zapcore.AddSync(&jacklog.Logger{
		Filename:   filename,
		MaxSize:    cfg.Rotate.MaxSize,
		MaxBackups: cfg.Rotate.MaxBackups,
		MaxAge:     cfg.Rotate.MaxAge,
		LocalTime:  cfg.Rotate.LocalTime,
		Compress:   cfg.Rotate.Compress,
	})

	levelEnabler, err := buildZapLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	options := buildZapOptions(cfg)
	options = append(options, zap.AddCallerSkip(1))

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(NewZapEncoderConfig()),
		writeSyncer,
		levelEnabler,
	)

	if enableStdout {
		stdoutWriteSyncer, _, _ := zap.Open("stdout")
		stdoutCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(NewZapEncoderConfig()),
			stdoutWriteSyncer,
			levelEnabler,
		)
		core = zapcore.NewTee(core, stdoutCore)
	}

	logger := logger{
		zap.New(core, options...).Sugar(),
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
	skyprome.LogCounter("INFO")
}

// Infof ...
func Infof(ctx context.Context, template string, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Infof(template, args...)
	skyprome.LogCounter("INFO")
}

// Infow ...
func Infow(ctx context.Context, msg string, keysAndValues ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Infow(msg, keysAndValues...)
	skyprome.LogCounter("INFO")
}

// Debug ...
func Debug(ctx context.Context, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Debug(args...)
	skyprome.LogCounter("DEBUG")
}

// Debugf ...
func Debugf(ctx context.Context, template string, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Debugf(template, args...)
	skyprome.LogCounter("DEBUG")
}

// Debugw ...
func Debugw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Debugw(msg, keysAndValues...)
	skyprome.LogCounter("DEBUG")
}

// Warn ...
func Warn(ctx context.Context, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Warn(args...)
	skyprome.LogCounter("WARN")
}

// Warnf ...
func Warnf(ctx context.Context, template string, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Warnf(template, args...)
	skyprome.LogCounter("WARN")
}

// Warnw ...
func Warnw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Warnw(msg, keysAndValues...)
	skyprome.LogCounter("WARN")
}

// Error ...
func Error(ctx context.Context, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Error(args...)
	skyprome.LogCounter("ERROR")
}

// Errorf ...
func Errorf(ctx context.Context, template string, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Errorf(template, args...)
	skyprome.LogCounter("ERROR")
}

// Errorw ...
func Errorw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Errorw(msg, keysAndValues...)
	skyprome.LogCounter("ERROR")
}

// DPanic ...
func DPanic(ctx context.Context, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.DPanic(args...)
	skyprome.LogCounter("PANIC")
}

// DPanicf ...
func DPanicf(ctx context.Context, template string, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.DPanicf(template, args...)
	skyprome.LogCounter("PANIC")
}

// DPanicw ...
func DPanicw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.DPanicw(msg, keysAndValues...)
	skyprome.LogCounter("PANIC")
}

// Panic ...
func Panic(ctx context.Context, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Panic(args...)
	skyprome.LogCounter("PANIC")
}

// Panicf ...
func Panicf(ctx context.Context, template string, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Panicf(template, args...)
	skyprome.LogCounter("PANIC")
}

// Panicw ...
func Panicw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Panicw(msg, keysAndValues...)
	skyprome.LogCounter("PANIC")
}

// Fatal ...
func Fatal(ctx context.Context, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Fatal(args...)
	skyprome.LogCounter("FATAL")
}

// Fatalf ...
func Fatalf(ctx context.Context, template string, args ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Fatalf(template, args...)
	skyprome.LogCounter("FATAL")
}

// Fatalw ...
func Fatalw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	logg := LoggerFromContext(ctx)
	logg.Fatalw(msg, keysAndValues...)
	skyprome.LogCounter("FATAL")
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

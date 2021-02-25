package log

import (
	"github.com/WiFeng/go-sky/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapEncoderConfig ...
func NewZapEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:     "time",
		LevelKey:    "level",
		NameKey:     "logger",
		CallerKey:   "caller",
		FunctionKey: zapcore.OmitKey,
		// FunctionKey:   "func",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		// EncodeLevel:    zapcore.LowercaseLevelEncoder,
		// EncodeTime:     zapcore.EpochTimeEncoder,
		// EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func buildZapOptions(cfg config.Log) []zap.Option {
	opts := []zap.Option{}

	if cfg.Development {
		opts = append(opts, zap.Development())
	}

	if !cfg.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}

	stackLevel := zap.ErrorLevel
	if cfg.Development {
		stackLevel = zap.WarnLevel
	}
	if !cfg.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}

	return opts
}

func buildZapLevel(level string) (zapcore.Level, error) {
	var alel zapcore.Level
	err := alel.UnmarshalText([]byte(level))
	return alel, err
}

func buildZapAtomicLevel(level string) (*zap.AtomicLevel, error) {
	alel := &zap.AtomicLevel{}
	err := alel.UnmarshalText([]byte(level))
	return alel, err
}

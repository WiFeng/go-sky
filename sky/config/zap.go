package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log log config
type Log struct {
	Level             zap.AtomicLevel        `toml:"level"`
	Development       bool                   `toml:"development"`
	DisableCaller     bool                   `toml:"disable_caller"`
	DisableStacktrace bool                   `toml:"disable_stacktrace"`
	Sampling          *zap.SamplingConfig    `toml:"sampling"`
	Encoding          string                 `toml:"encoding"`
	EncoderConfig     zapEncoderConfig       `toml:"encoder_config"`
	OutputPaths       []string               `toml:"output_paths"`
	ErrorOutputPaths  []string               `toml:"error_output_paths"`
	InitialFields     map[string]interface{} `toml:"initial_fields"`
}

type zapEncoderConfig struct {
	MessageKey    string `toml:"messageKey"`
	LevelKey      string `toml:"levelKey"`
	TimeKey       string `toml:"timeKey"`
	NameKey       string `toml:"nameKey"`
	CallerKey     string `toml:"callerKey"`
	StacktraceKey string `toml:"stacktraceKey"`
	LineEnding    string `toml:"lineEnding"`

	// EncodeLevel    LevelEncoder    `toml:"levelEncoder"`
	// EncodeTime     TimeEncoder     `toml:"timeEncoder"`
	// EncodeDuration DurationEncoder `toml:"durationEncoder"`
	// EncodeCaller   CallerEncoder   `toml:"callerEncoder"`

	// EncodeName NameEncoder `toml:"nameEncoder"`
}

// NewZapConfig new zap config
func NewZapConfig(zapConf Log) zap.Config {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	zapConfig := zap.Config{
		Level:             zapConf.Level,
		Development:       zapConf.Development,
		DisableCaller:     zapConf.DisableCaller,
		DisableStacktrace: zapConf.DisableStacktrace,
		Sampling:          zapConf.Sampling,
		Encoding:          zapConf.Encoding,
		EncoderConfig:     encoderConfig,
		OutputPaths:       zapConf.OutputPaths,
		ErrorOutputPaths:  zapConf.ErrorOutputPaths,
		InitialFields:     zapConf.InitialFields,
	}

	return zapConfig
}

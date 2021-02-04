package config

// Log log config
type Log struct {
	Level             string `toml:"level"`
	Development       bool   `toml:"development"`
	DisableCaller     bool   `toml:"disableCaller"`
	DisableStacktrace bool   `toml:"disableStacktrace"`
	OutputPath        string `toml:"outputPath"`
	Rotate            LogRotate
}

// LogRotate ...
type LogRotate struct {
	Filename   string
	MaxSize    int // megabytes
	MaxBackups int
	MaxAge     int // days
	LocalTime  bool
	Compress   bool
}

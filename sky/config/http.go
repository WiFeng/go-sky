package config

import "time"

// HTTP http config
type HTTP struct {
	Addr string
}

// PProf ...
type PProf struct {
	Host string
	Port int
}

// Client ...
type Client struct {
	Name     string
	Host     string
	Port     int
	Protocol string

	CustomTranport bool
	MillSecUnit    bool

	Timeout   time.Duration
	Transport HTTPTransport
}

// HTTPTransport ...
type HTTPTransport struct {
	Customized  bool // 定制Transport参数
	MillSecUnit bool // 是否使用毫秒单位

	IdleConnTimeout       time.Duration
	TLSHandshakeTimeout   time.Duration
	ExpectContinueTimeout time.Duration
	ResponseHeaderTimeout time.Duration

	MaxConnsPerHost     int
	MaxIdleConns        int
	MaxIdleConnsPerHost int

	DisableKeepAlives  bool
	DisableCompression bool
}

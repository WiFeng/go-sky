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

	Timeout        time.Duration
	TimeoutMillSec time.Duration
	Transport      HTTPTransport
}

// HTTPTransport ...
type HTTPTransport struct {
	Customized  bool
	MillSecUnit bool

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

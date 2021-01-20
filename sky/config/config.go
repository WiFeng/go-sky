package config

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
)

// Config config
type Config struct {
	Server        Server
	Redis         []Redis
	DB            []DB
	Client        []Client
	Elasticsearch []Elasticsearch
}

// Server server config
type Server struct {
	Name  string
	HTTP  HTTP
	PProf PProf
	Log   Log
}

// HTTP http config
type HTTP struct {
	Addr string
}

// PProf ...
type PProf struct {
	Host string
	Port int
}

// HTTPTransport ...
type HTTPTransport struct {
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

// Redis redis config
type Redis struct {
	Name string
	Host string
	Port int
	Auth string
	DB   int
}

// DB config ...
type DB struct {
	Name    string
	Driver  string
	Host    string
	Port    int
	User    string
	Pass    string
	DB      string
	Charset string
}

// Elasticsearch ...
type Elasticsearch struct {
	Addresses []string
	Username  string
	Password  string
	Transport HTTPTransport
}

// Init ...
func Init(dir, env string, conf *Config) (string, error) {
	return LoadConfig(dir, "config", env, conf)
}

// DecodeFile decode toml file
func DecodeFile(fpath string, v interface{}) (toml.MetaData, error) {
	return toml.DecodeFile(fpath, v)
}

// LoadConfig ...
func LoadConfig(dir string, name string, env string, conf interface{}) (string, error) {
	var confFile = fmt.Sprintf("%s/%s.toml", dir, name)

	if env != "" {
		confFile = fmt.Sprintf("./conf/%s_%s.toml", name, env)
	}

	if _, err := DecodeFile(confFile, conf); err != nil {
		return confFile, err
	}

	return confFile, nil
}

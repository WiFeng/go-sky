package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// Config config
type Config struct {
	Server        Server
	Redis         []Redis
	Database      []Database
	Client        []Client
	Elasticsearch []Elasticsearch
	Kafka         []Kafka
}

// Server server config
type Server struct {
	Name  string
	HTTP  HTTP
	PProf PProf
	Log   Log
}

// Init ...
func Init(dir, env string, conf *Config) (string, error) {
	return LoadConfig(dir, "config", env, conf)
}

func decodeFile(fpath string, v interface{}) (toml.MetaData, error) {
	return toml.DecodeFile(fpath, v)
}

// LoadConfig ...
func LoadConfig(dir string, name string, env string, conf interface{}) (string, error) {
	var confFile = fmt.Sprintf("%s/%s.toml", dir, name)

	if env != "" {
		confFile = fmt.Sprintf("./conf/%s_%s.toml", name, env)
	}

	if _, err := decodeFile(confFile, conf); err != nil {
		return confFile, err
	}

	return confFile, nil
}

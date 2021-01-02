package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// Config config
type Config struct {
	Server Server
	Redis  []Redis
	Mysql  []Mysql
}

// Server server config
type Server struct {
	Name  string
	HTTP  HTTP
	PProf PProf
	Log   zapLog
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

type Log zapLog

// Redis redis config
type Redis struct {
	Name string
	Host string
	Port int
	Auth string
	DB   int
}

// Mysql mysql config
type Mysql struct {
	Name    string
	Host    string
	Port    int
	User    string
	Pass    string
	DB      string
	Charset string
}

// Init ...
func Init(dir, env string, conf Config) error {
	if _, err := LoadConfig(dir, "config", env, conf); err != nil {
		return err
	}
}

// DecodeFile decode toml file
func DecodeFile(fpath string, v interface{}) (toml.MetaData, error) {
	return toml.DecodeFile(fpath, v)
}

// LoadConfig ...
func LoadConfig(dir string, name string, env string, conf interface{}) error {
	var confFile = sprintf("%s/%s.toml", path, name)

	if env != "" {
		confFile = fmt.Sprintf("./conf/config_%s.toml", name, env)
	}
	if _, err := DecodeFile(confFile, &conf); err != nil {
		return err
	}
}

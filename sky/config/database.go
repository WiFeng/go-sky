package config

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

package config

// Database config ...
type Database struct {
	Name    string
	Driver  string
	Host    string
	Port    int
	User    string
	Pass    string
	DB      string
	Charset string
}

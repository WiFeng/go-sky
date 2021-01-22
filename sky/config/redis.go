package config

// Redis redis config
type Redis struct {
	Name string
	Host string
	Port int
	Auth string
	DB   int
}

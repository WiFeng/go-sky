package config

// Database config ...
type Database struct {
	Name              string
	Driver            string
	DataSource        string
	Host              string
	Port              int
	User              string
	Pass              string
	DB                string
	Charset           string
	InterpolateParams string
}

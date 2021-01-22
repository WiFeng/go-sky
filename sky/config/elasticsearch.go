package config

// Elasticsearch ...
type Elasticsearch struct {
	Name           string
	Addrs          []string
	Username       string
	Password       string
	CustomTranport bool
	Transport      HTTPTransport
}

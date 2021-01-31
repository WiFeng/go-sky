package config

// Metrics ...
type Metrics struct {
	Prometheus Prometheus
}

// Prometheus ...
type Prometheus struct {
	Addr string
}

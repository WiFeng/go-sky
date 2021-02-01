package config

// Metrics ...
type Metrics struct {
	Prometheus Prometheus
}

// Prometheus ...
type Prometheus struct {
	Addr string

	DisableHTTPRequestsTotalCounter      bool
	DisableHTTPRequestsDurationHistogram bool
	DisableHTTPRequestsDurationSummary   bool

	HTTPRequestsDurationHistogramBuckets  []float64
	HTTPRequestsDurationSummaryObjectives map[float64]float64
}

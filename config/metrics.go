package config

// Metrics ...
type Metrics struct {
	Prometheus Prometheus
}

// Prometheus ...
type Prometheus struct {
	Addr string

	DisableHTTPServerRequestsTotalCounter      bool
	DisableHTTPServerRequestsDurationHistogram bool
	DisableHTTPServerRequestsDurationSummary   bool
	DisableHTTPClientRequestsTotalCounter      bool
	DisableHTTPClientRequestsDurationHistogram bool
	DisableHTTPClientRequestsDurationSummary   bool
	DisableLogTotalCounter                     bool

	HTTPServerRequestsDurationHistogramBuckets  []float64
	HTTPServerRequestsDurationSummaryObjectives map[float64]float64
	HTTPClientRequestsDurationHistogramBuckets  []float64
	HTTPClientRequestsDurationSummaryObjectives map[float64]float64
}

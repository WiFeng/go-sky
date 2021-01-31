package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotalCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "The total number of http requests",
		},
		[]string{"code", "method", "path"},
	)

	httpRequestsDurationHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds_histogram",
			Help: "A histogram of latencies for requests.",
			//Buckets: []float64{.25, .5, 1, 2.5, 5, 10},
		},
	)

	httpRequestsDurationSummary = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "http_request_duration_seconds_summary",
		Help:       "A summary of latencies for requests.",
		Objectives: map[float64]float64{0.5: 0.05, 0.8: 0.001, 0.9: 0.01, 0.95: 0.01},
	})
)

// HTTPRequestsTotalCounter ...
func HTTPRequestsTotalCounter(code int, method string, path string) {
	labels := prometheus.Labels{
		"code":   sanitizeCode(code),
		"method": sanitizeMethod(method),
		"path":   path,
	}
	httpRequestsTotalCounter.With(labels).Inc()
}

// HTTPRequestsDurationHistogram ...
func HTTPRequestsDurationHistogram(duration float64) {
	httpRequestsDurationHistogram.Observe(duration)
}

// HTTPRequestsDurationSummary ...
func HTTPRequestsDurationSummary(duration float64) {
	httpRequestsDurationSummary.Observe(duration)
}

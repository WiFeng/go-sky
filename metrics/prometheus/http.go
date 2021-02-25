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
		[]string{"service", "code", "method", "path"},
	)

	httpRequestsDurationHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds_histogram",
			Help:    "A histogram of latencies for requests.",
			Buckets: promecfg.HTTPRequestsDurationHistogramBuckets,
		},
		[]string{"service"},
	)

	httpRequestsDurationSummary = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_request_duration_seconds_summary",
			Help:       "A summary of latencies for requests.",
			Objectives: promecfg.HTTPRequestsDurationSummaryObjectives,
		},
		[]string{"service"},
	)
)

// HTTPRequestsTotalCounter ...
func HTTPRequestsTotalCounter(code int, method string, path string) {
	if promecfg.DisableHTTPRequestsTotalCounter {
		return
	}

	labels := prometheus.Labels{
		"service": service,
		"code":    sanitizeCode(code),
		"method":  sanitizeMethod(method),
		"path":    path,
	}
	httpRequestsTotalCounter.With(labels).Inc()
}

// HTTPRequestsDurationHistogram ...
func HTTPRequestsDurationHistogram(duration float64) {
	if promecfg.DisableHTTPRequestsDurationHistogram {
		return
	}

	labels := prometheus.Labels{
		"service": service,
	}
	httpRequestsDurationHistogram.With(labels).Observe(duration)
}

// HTTPRequestsDurationSummary ...
func HTTPRequestsDurationSummary(duration float64) {
	if promecfg.DisableHTTPRequestsDurationSummary {
		return
	}

	labels := prometheus.Labels{
		"service": service,
	}
	httpRequestsDurationSummary.With(labels).Observe(duration)
}

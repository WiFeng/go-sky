package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpServerRequestsTotalCounter      *prometheus.CounterVec
	httpServerRequestsDurationHistogram *prometheus.HistogramVec
	httpServerRequestsDurationSummary   *prometheus.SummaryVec
)

func HttpServerInit() {
	httpServerRequestsTotalCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_server_request_total",
			Help: "The total number of http requests",
		},
		[]string{"service", "code", "method", "path"},
	)

	httpServerRequestsDurationHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_server_request_duration_seconds_histogram",
			Help:    "A histogram of latencies for requests.",
			Buckets: promecfg.HTTPServerRequestsDurationHistogramBuckets,
		},
		[]string{"service", "code", "method", "path"},
	)

	httpServerRequestsDurationSummary = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_server_request_duration_seconds_summary",
			Help:       "A summary of latencies for requests.",
			Objectives: promecfg.HTTPServerRequestsDurationSummaryObjectives,
		},
		[]string{"service", "code", "method", "path"},
	)
}

// HTTPServerRequestsTotalCounter ...
func HTTPServerRequestsTotalCounter(code int, method string, path string) {
	if promecfg.DisableHTTPServerRequestsTotalCounter {
		return
	}

	labels := prometheus.Labels{
		"service": service,
		"code":    sanitizeCode(code),
		"method":  sanitizeMethod(method),
		"path":    path,
	}
	httpServerRequestsTotalCounter.With(labels).Inc()
}

// HTTPServerRequestsDurationHistogram ...
func HTTPServerRequestsDurationHistogram(code int, method string, path string, duration float64) {
	if promecfg.DisableHTTPServerRequestsDurationHistogram {
		return
	}

	labels := prometheus.Labels{
		"service": service,
		"code":    sanitizeCode(code),
		"method":  sanitizeMethod(method),
		"path":    path,
	}
	httpServerRequestsDurationHistogram.With(labels).Observe(duration)
}

// HTTPServerRequestsDurationSummary ...
func HTTPServerRequestsDurationSummary(code int, method string, path string, duration float64) {
	if promecfg.DisableHTTPServerRequestsDurationSummary {
		return
	}

	labels := prometheus.Labels{
		"service": service,
		"code":    sanitizeCode(code),
		"method":  sanitizeMethod(method),
		"path":    path,
	}
	httpServerRequestsDurationSummary.With(labels).Observe(duration)
}

package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpClientRequestsTotalCounter      *prometheus.CounterVec
	httpClientRequestsDurationHistogram *prometheus.HistogramVec
	httpClientRequestsDurationSummary   *prometheus.SummaryVec
)

func HttpClientInit() {
	httpClientRequestsTotalCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_client_request_total",
			Help: "The total number of http requests",
		},
		[]string{"service", "peer", "code", "method", "path"},
	)

	httpClientRequestsDurationHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_client_request_duration_seconds_histogram",
			Help:    "A histogram of latencies for requests.",
			Buckets: promecfg.HTTPClientRequestsDurationHistogramBuckets,
		},
		[]string{"service", "peer", "code", "method", "path"},
	)

	httpClientRequestsDurationSummary = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_client_request_duration_seconds_summary",
			Help:       "A summary of latencies for requests.",
			Objectives: promecfg.HTTPClientRequestsDurationSummaryObjectives,
		},
		[]string{"service", "peer", "code", "method", "path"},
	)
}

// HTTPClientRequestsTotalCounter ...
func HTTPClientRequestsTotalCounter(peer string, code int, method string, path string) {
	if promecfg.DisableHTTPClientRequestsTotalCounter {
		return
	}

	labels := prometheus.Labels{
		"service": service,
		"peer":    peer,
		"code":    sanitizeCode(code),
		"method":  sanitizeMethod(method),
		"path":    path,
	}
	httpClientRequestsTotalCounter.With(labels).Inc()
}

// HTTPClientRequestsDurationHistogram ...
func HTTPClientRequestsDurationHistogram(peer string, code int, method string, path string, duration float64) {
	if promecfg.DisableHTTPClientRequestsDurationHistogram {
		return
	}

	labels := prometheus.Labels{
		"service": service,
		"peer":    peer,
		"code":    sanitizeCode(code),
		"method":  sanitizeMethod(method),
		"path":    path,
	}
	httpClientRequestsDurationHistogram.With(labels).Observe(duration)
}

// HTTPClientRequestsDurationSummary ...
func HTTPClientRequestsDurationSummary(peer string, code int, method string, path string, duration float64) {
	if promecfg.DisableHTTPClientRequestsDurationSummary {
		return
	}

	labels := prometheus.Labels{
		"service": service,
		"peer":    peer,
		"code":    sanitizeCode(code),
		"method":  sanitizeMethod(method),
		"path":    path,
	}
	httpClientRequestsDurationSummary.With(labels).Observe(duration)
}

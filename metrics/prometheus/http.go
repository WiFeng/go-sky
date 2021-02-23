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

	httpRequestsDurationHistogram = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds_histogram",
			Help:    "A histogram of latencies for requests.",
			Buckets: promecfg.HTTPRequestsDurationHistogramBuckets,
		},
	)

	httpRequestsDurationSummary = promauto.NewSummary(prometheus.SummaryOpts{
		Name:       "http_request_duration_seconds_summary",
		Help:       "A summary of latencies for requests.",
		Objectives: promecfg.HTTPRequestsDurationSummaryObjectives,
	})
)

// HTTPRequestsTotalCounter ...
func HTTPRequestsTotalCounter(code int, method string, path string) {
	if promecfg.DisableHTTPRequestsTotalCounter {
		return
	}

	labels := prometheus.Labels{
		"code":   sanitizeCode(code),
		"method": sanitizeMethod(method),
		"path":   path,
	}
	httpRequestsTotalCounter.With(labels).Inc()
}

// HTTPRequestsDurationHistogram ...
func HTTPRequestsDurationHistogram(duration float64) {
	if promecfg.DisableHTTPRequestsDurationHistogram {
		return
	}

	httpRequestsDurationHistogram.Observe(duration)
}

// HTTPRequestsDurationSummary ...
func HTTPRequestsDurationSummary(duration float64) {
	if promecfg.DisableHTTPRequestsDurationSummary {
		return
	}

	httpRequestsDurationSummary.Observe(duration)
}

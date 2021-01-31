package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotalCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "The total number of http requests",
	},
		[]string{"code", "method", "path"},
	)
)

// HTTPRequestsTotalCounterInc ...
func HTTPRequestsTotalCounterInc(code int, method string, path string) {
	labels := prometheus.Labels{
		"code":   sanitizeCode(code),
		"method": sanitizeMethod(method),
		"path":   path,
	}
	httpRequestsTotalCounter.With(labels).Inc()
}

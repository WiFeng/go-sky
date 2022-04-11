package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	_logCounter *prometheus.CounterVec
	logCounter  *prometheus.CounterVec
)

func LogInit() {
	_logCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "srv_log_total",
			Help: "The total number of server log",
		},
		[]string{"service", "level"},
	)
	logCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "svc_log_total",
			Help: "The total number of server log",
		},
		[]string{"service", "level"},
	)
}

// LogCounter ...
func LogCounter(level string) {
	if promecfg.DisableLogTotalCounter {
		return
	}

	if logCounter == nil {
		return
	}

	labels := prometheus.Labels{
		"service": service,
		"level":   level,
	}
	_logCounter.With(labels).Inc()
	logCounter.With(labels).Inc()
}

package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	logCounter *prometheus.CounterVec
)

func LogInit() {
	logCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "srv_log_total",
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
	logCounter.With(labels).Inc()
}

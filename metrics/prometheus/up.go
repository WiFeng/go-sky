package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func UpInit() {
	upCounter := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "svc_up",
			Help: "The server up state",
		},
		[]string{"service"},
	)
	labels := prometheus.Labels{
		"service": service,
	}
	upCounter.With(labels).Set(1)
}

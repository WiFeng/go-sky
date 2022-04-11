package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func UpInit() {
	_upCounter := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "srv_up",
			Help: "The server up state",
		},
		[]string{"service"},
	)
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
	_upCounter.With(labels).Set(1)
	upCounter.With(labels).Set(1)
}

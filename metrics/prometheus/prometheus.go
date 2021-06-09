package prometheus

import "github.com/WiFeng/go-sky/config"

var (
	DefaultBuckets    = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
	DefaultObjectives = map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001, 0.999: 0.0001, 0.9999: 0.00001}
)

var (
	service  string
	promecfg config.Prometheus
)

func SetPromeCfg(cfg config.Prometheus) {
	promecfg = cfg
}

func SetPromeService(s string) {
	service = s
}

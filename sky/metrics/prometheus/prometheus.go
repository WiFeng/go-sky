package prometheus

import (
	"context"
	"net/http"

	"github.com/WiFeng/go-sky/sky/config"
	"github.com/WiFeng/go-sky/sky/log"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	promecfg          = config.Prometheus{}
	defaultBuckets    = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
	defaultObjectives = map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001, 0.999: 0.0001, 0.9999: 0.00001}
)

// Init ...
func Init(ctx context.Context, cfg config.Prometheus) {
	if cfg.Addr == "nil" {
		return
	}

	go func() {
		log.Infof(ctx, "Start HTTP Prometheus metrics. http://%s", cfg.Addr)
		log.Fatal(ctx, http.ListenAndServe(cfg.Addr, promhttp.Handler()))
	}()

	if len(cfg.HTTPRequestsDurationHistogramBuckets) < 1 {
		cfg.HTTPRequestsDurationHistogramBuckets = defaultBuckets
	}

	if len(cfg.HTTPRequestsDurationSummaryObjectives) < 1 {
		cfg.HTTPRequestsDurationSummaryObjectives = defaultObjectives
	}

	promecfg = cfg
}

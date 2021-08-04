package metrics

import (
	"context"
	"net/http"

	"github.com/WiFeng/go-sky/config"
	"github.com/WiFeng/go-sky/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	skyprome "github.com/WiFeng/go-sky/metrics/prometheus"
)

// Init ...
func Init(ctx context.Context, serviceName string, cfg config.Metrics) {
	initPrometheus(ctx, serviceName, cfg.Prometheus)
}

// Init ...
func initPrometheus(ctx context.Context, serviceName string, cfg config.Prometheus) {

	if cfg.Addr == "" {
		return
	}

	if len(cfg.HTTPServerRequestsDurationHistogramBuckets) < 1 {
		cfg.HTTPServerRequestsDurationHistogramBuckets = skyprome.DefaultBuckets
	}

	if len(cfg.HTTPServerRequestsDurationSummaryObjectives) < 1 {
		cfg.HTTPServerRequestsDurationSummaryObjectives = skyprome.DefaultObjectives
	}

	if len(cfg.HTTPClientRequestsDurationHistogramBuckets) < 1 {
		cfg.HTTPClientRequestsDurationHistogramBuckets = skyprome.DefaultBuckets
	}

	if len(cfg.HTTPClientRequestsDurationSummaryObjectives) < 1 {
		cfg.HTTPClientRequestsDurationSummaryObjectives = skyprome.DefaultObjectives
	}

	skyprome.SetPromeCfg(cfg)
	skyprome.SetPromeService(serviceName)

	skyprome.UpInit()
	skyprome.LogInit()
	skyprome.HttpServerInit()
	skyprome.HttpClientInit()

	go func() {
		log.Infof(ctx, "Start HTTP Prometheus metrics. http://%s", cfg.Addr)
		log.Fatal(ctx, http.ListenAndServe(cfg.Addr, promhttp.Handler()))
	}()

}

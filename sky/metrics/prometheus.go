package metrics

import (
	"context"
	"net/http"

	"github.com/WiFeng/go-sky/sky/config"
	"github.com/WiFeng/go-sky/sky/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func initPrometheus(ctx context.Context, cfg config.Prometheus) {
	if cfg.Addr == "nil" {
		return
	}

	go func() {
		log.Infof(ctx, "Start HTTP Prometheus metrics. http://%s", cfg.Addr)
		log.Fatal(ctx, http.ListenAndServe(cfg.Addr, promhttp.Handler()))
	}()
}

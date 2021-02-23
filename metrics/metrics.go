package metrics

import (
	"context"

	"github.com/WiFeng/go-sky/config"
	"github.com/WiFeng/go-sky/metrics/prometheus"
)

// Init ...
func Init(ctx context.Context, cfg config.Metrics) {
	prometheus.Init(ctx, cfg.Prometheus)
}

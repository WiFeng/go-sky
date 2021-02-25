package metrics

import (
	"context"

	"github.com/WiFeng/go-sky/config"
	"github.com/WiFeng/go-sky/metrics/prometheus"
)

// Init ...
func Init(ctx context.Context, serviceName string, cfg config.Metrics) {
	prometheus.Init(ctx, serviceName, cfg.Prometheus)
}

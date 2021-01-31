package metrics

import (
	"context"

	"github.com/WiFeng/go-sky/sky/config"
)

// Init ...
func Init(ctx context.Context, cfg config.Metrics) {
	initPrometheus(ctx, cfg.Prometheus)
}

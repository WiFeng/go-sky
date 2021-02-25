package trace

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/WiFeng/go-sky/config"
	"github.com/WiFeng/go-sky/helper"
	jaegerconfig "github.com/uber/jaeger-client-go/config"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics/prometheus"
)

// Init ...
func Init(ctx context.Context, serviceName string, cfg config.Trace) {
	metricsFactory := prometheus.New()
	tracer, tracerCloser, err := jaegerconfig.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegerconfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegerconfig.ReporterConfig{
			CollectorEndpoint:   cfg.Reporter.CollectorEndpoint,
			LocalAgentHostPort:  cfg.Reporter.LocalAgentHostPort,
			BufferFlushInterval: cfg.Reporter.BufferFlushInterval * time.Second,
		},
	}.NewTracer(
		jaegerconfig.Metrics(metricsFactory),
	)

	if err != nil {
		fmt.Println("Init trace error. ", err)
		os.Exit(1)
		return
	}

	opentracing.InitGlobalTracer(tracer)
	helper.AddDeferFunc(func() {
		tracerCloser.Close()
	})

}

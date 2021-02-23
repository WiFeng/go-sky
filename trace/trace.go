package trace

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/WiFeng/go-sky/config"
	"github.com/WiFeng/go-sky/helper"
	"github.com/uber/jaeger-client-go"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics/prometheus"
)

// Init ...
func Init(ctx context.Context, serviceName string, cfg config.Trace) {
	metricsFactory := prometheus.New()

	logger := jaeger.StdLogger
	// logger := log.GetDefaultLogger()
	loggerOption := jaegerconfig.Logger(jaegerlog.DebugLogAdapter(logger))
	tracer, tracerCloser, err := jaegerconfig.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegerconfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegerconfig.ReporterConfig{
			// LogSpans:            true,
			// LocalAgentHostPort:  "localhost:6831",
			// BufferFlushInterval: time.Second * 1,
			CollectorEndpoint:   cfg.Reporter.CollectorEndpoint,
			LocalAgentHostPort:  cfg.Reporter.LocalAgentHostPort,
			BufferFlushInterval: cfg.Reporter.BufferFlushInterval * time.Second,
		},
	}.NewTracer(
		jaegerconfig.Metrics(metricsFactory),
		loggerOption,
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

package trace

import (
	"context"
	"io"
	"time"

	"github.com/WiFeng/go-sky/sky/config"
	"github.com/uber/jaeger-client-go"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics/prometheus"
)

// Init ...
func Init(ctx context.Context, serviceName string, cfg config.Trace) (opentracing.Tracer, io.Closer, error) {
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
	opentracing.InitGlobalTracer(tracer)
	return tracer, tracerCloser, err
}

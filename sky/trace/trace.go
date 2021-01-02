package trace

import (
	"context"
	"io"

	"github.com/opentracing/opentracing-go"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
)

// Init ...
func Init(serviceName string) (opentracing.Tracer, io.Closer, error) {
	metricsFactory := prometheus.New()
	tracer, tracerCloser, err = jaegerconfig.Configuration{
		ServiceName: serviceName,
	}.NewTracer(
		jaegerconfig.Metrics(metricsFactory),
	)

	defer opentracing.InitGlobalTracer(tracer)

	return tracer, tracerCloser, nil
}

// GetTraceID Get trace id from the context.
func GetTraceID(ctx context.Context) string {
	var traceID string
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		spanContext := span.Context()
		jeagerSpanContext, ok := spanContext.(jaegerclient.SpanContext)
		if ok {
			traceID = jeagerSpanContext.TraceID().String()
		}
	}

	return traceID
}

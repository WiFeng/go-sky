package trace

import (
	"context"
	"fmt"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// StartSpan ...
func StartSpan(ctx context.Context, r *http.Request) context.Context {
	var serverSpan opentracing.Span
	var appSpecificOperationName = fmt.Sprintf("[%s]%s", r.Method, r.URL)

	wireContext, err := opentracing.GlobalTracer().Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header))
	if err != nil {
		// Optionally record something about err here
	}

	// Create the span referring to the RPC client if available.
	// If wireContext == nil, a root span will be created.
	serverSpan = opentracing.StartSpan(
		appSpecificOperationName,
		ext.RPCServerOption(wireContext))

	// defer serverSpan.Finish()

	// ctx = opentracing.ContextWithSpan(context.Background(), serverSpan)
	newCtx := opentracing.ContextWithSpan(ctx, serverSpan)

	return newCtx
}

// FinishSpan ...
func FinishSpan(ctx context.Context, w http.ResponseWriter) context.Context {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span.Finish()
	}
	return ctx
}

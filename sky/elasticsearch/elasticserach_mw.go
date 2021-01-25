package elasticsearch

import (
	"fmt"
	"net/http"

	skyhttp "github.com/WiFeng/go-sky/sky/http"
	"github.com/WiFeng/go-sky/sky/log"
	kitopentracing "github.com/go-kit/kit/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	opentracingext "github.com/opentracing/opentracing-go/ext"
)

// ==========================================
// RoundTripper Middleware
// ==========================================

// RoundTripperTracingMiddleware ...
func RoundTripperTracingMiddleware(next http.RoundTripper) http.RoundTripper {
	return skyhttp.RoundTripperFunc(func(req *http.Request) (resp *http.Response, err error) {
		var ctx = req.Context()
		var logger = log.LoggerFromContext(ctx)
		var tracer = opentracing.GlobalTracer()

		var parentSpan opentracing.Span
		var childSpan opentracing.Span

		defer func() {
			if childSpan == nil {
				return
			}
			if err != nil {
				opentracingext.Error.Set(childSpan, true)
				childSpan.SetTag("http.error", err.Error())
				childSpan.Finish()
				return
			}
			if resp.StatusCode >= 400 {
				opentracingext.Error.Set(childSpan, true)
			}

			opentracingext.HTTPStatusCode.Set(childSpan, uint16(resp.StatusCode))
			childSpan.Finish()
		}()

		if parentSpan = opentracing.SpanFromContext(ctx); parentSpan != nil {
			childSpan = parentSpan.Tracer().StartSpan(
				fmt.Sprintf("[%s] %s", req.Method, req.URL.Path),
				opentracing.ChildOf(parentSpan.Context()),
				opentracing.Tag{Key: string(opentracingext.Component), Value: "elasticsearch"},
				opentracingext.SpanKindRPCClient,
			)
			ctx = opentracing.ContextWithSpan(ctx, childSpan)
		}

		kitopentracing.ContextToHTTP(tracer, logger)(ctx, req)
		resp, err = next.RoundTrip(req)
		return
	})
}

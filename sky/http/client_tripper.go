package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/WiFeng/go-sky/sky/config"
	"github.com/WiFeng/go-sky/sky/log"
	kitopentracing "github.com/go-kit/kit/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	opentracingext "github.com/opentracing/opentracing-go/ext"
)

// RoundTripperFunc ...
type RoundTripperFunc func(*http.Request) (*http.Response, error)

// RoundTrip ...
func (r RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req)
}

// RoundTripper ...
type RoundTripper struct {
	base        http.RoundTripper
	middlewares []RoundTripperMiddleware
}

// NewTransport ...
func NewTransport(cf config.HTTPTransport) *http.Transport {
	tr := http.DefaultTransport.(*http.Transport).Clone()

	if !cf.Customized {
		return tr
	}

	unit := time.Second
	if cf.MillSecUnit {
		unit = time.Millisecond
	}

	tr = &http.Transport{
		MaxConnsPerHost:     cf.MaxConnsPerHost,
		MaxIdleConns:        cf.MaxIdleConns,
		MaxIdleConnsPerHost: cf.MaxIdleConnsPerHost,

		IdleConnTimeout:       cf.IdleConnTimeout * unit,
		TLSHandshakeTimeout:   cf.TLSHandshakeTimeout * unit,
		ExpectContinueTimeout: cf.ExpectContinueTimeout * unit,
		ResponseHeaderTimeout: cf.ResponseHeaderTimeout * unit,

		DisableKeepAlives:  cf.DisableKeepAlives,
		DisableCompression: cf.DisableCompression,
	}

	return tr
}

// NewRoundTripper ...
func NewRoundTripper(base http.RoundTripper, mwf ...RoundTripperMiddlewareFunc) *RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	roundTripper := &RoundTripper{
		base: base,
	}
	roundTripper.Use(mwf...)
	return roundTripper
}

// NewRoundTripperFromConfig ...
func NewRoundTripperFromConfig(cf config.HTTPTransport) *RoundTripper {
	return NewRoundTripper(NewTransport(cf))
}

// Use ...
func (r *RoundTripper) Use(mwf ...RoundTripperMiddlewareFunc) {
	for _, fn := range mwf {
		r.middlewares = append(r.middlewares, fn)
	}
}

// RoundTrip ...
func (r RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var rr = r.base
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		rr = r.middlewares[i].Middleware(rr)
	}
	return rr.RoundTrip(req)
}

// RoundTripperMiddleware ...
type RoundTripperMiddleware interface {
	Middleware(http.RoundTripper) http.RoundTripper
}

// RoundTripperMiddlewareFunc ...
type RoundTripperMiddlewareFunc func(http.RoundTripper) http.RoundTripper

// Middleware allows MiddlewareFunc to implement the middleware interface.
func (mw RoundTripperMiddlewareFunc) Middleware(roundTripper http.RoundTripper) http.RoundTripper {
	return mw(roundTripper)
}

// ==========================================
// RoundTripper Middleware
// ==========================================

// RoundTripperTracingMiddleware ...
func RoundTripperTracingMiddleware(next http.RoundTripper) http.RoundTripper {
	return RoundTripperFunc(func(req *http.Request) (resp *http.Response, err error) {
		var ctx = req.Context()
		var logger = log.LoggerFromContext(ctx)
		var tracer = opentracing.GlobalTracer()

		var parentSpan opentracing.Span
		var childSpan opentracing.Span

		defer func() {
			if childSpan != nil {
				opentracingext.HTTPStatusCode.Set(childSpan, uint16(resp.StatusCode))
				childSpan.Finish()
			}
		}()

		if parentSpan = opentracing.SpanFromContext(ctx); parentSpan != nil {
			childSpan = parentSpan.Tracer().StartSpan(
				fmt.Sprintf("[%s] %s", req.Method, req.URL.Path),
				opentracing.ChildOf(parentSpan.Context()),
				opentracingext.SpanKindRPCClient,
			)
			ctx = opentracing.ContextWithSpan(ctx, childSpan)
		}

		kitopentracing.ContextToHTTP(tracer, logger)(ctx, req)
		resp, err = next.RoundTrip(req)
		return
	})
}

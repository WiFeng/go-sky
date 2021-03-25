package http

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/WiFeng/go-sky/config"
	"github.com/WiFeng/go-sky/log"
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

type wapperBody struct {
	base  io.ReadCloser
	buff  io.ReadCloser
	bytes []byte
}

func newWapperBody(base io.ReadCloser) *wapperBody {
	bodyBytes, _ := ioutil.ReadAll(base)
	bodyBuffer := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return &wapperBody{
		base:  base,
		buff:  bodyBuffer,
		bytes: bodyBytes,
	}
}

func (b *wapperBody) Read(p []byte) (n int, err error) {
	return b.buff.Read(p)
}

func (b *wapperBody) Close() error {
	return b.base.Close()
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
				opentracingext.SpanKindRPCClient,
			)
			ctx = opentracing.ContextWithSpan(ctx, childSpan)
		}

		kitopentracing.ContextToHTTP(tracer, logger)(ctx, req)
		resp, err = next.RoundTrip(req)
		return
	})
}

// RoundTripperLoggingMiddleware ...
func RoundTripperLoggingMiddleware(next http.RoundTripper) http.RoundTripper {
	return RoundTripperFunc(func(req *http.Request) (resp *http.Response, err error) {
		var ctx = req.Context()
		var reqBodyBytes = make([]byte, 0)
		var respBodyBytes = make([]byte, 0)

		if req.Body != nil {
			body := newWapperBody(req.Body)
			reqBodyBytes = body.bytes
			req.Body = body
		}

		defer func(begin time.Time) {
			var reqBody string
			var respBody string
			var respStatus int

			if resp != nil {
				respStatus = resp.StatusCode
				if resp.Body != nil {
					// respBodyBytes, _ = ioutil.ReadAll(resp.Body)
					// resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBodyBytes))

					body := newWapperBody(resp.Body)
					respBodyBytes = body.bytes
					resp.Body = body
				}
			}

			reqBody = string(reqBodyBytes)
			respBody = string(respBodyBytes)
			if len(reqBody) > 800 {
				reqBody = reqBody[0:800]
			}
			if len(respBody) > 500 {
				respBody = respBody[0:500]
			}

			if err != nil || respStatus >= 500 {
				log.Errorw(ctx, fmt.Sprintf("%s %s?%s", req.Method, req.URL.Path, req.URL.RawQuery), log.TypeKey, log.TypeValRPC, "host", req.Host, "req", reqBody,
					"resp", respBody, "status", respStatus, "request_time", fmt.Sprintf("%.3f", float32(time.Since(begin).Microseconds())/1000), "err", err)
				return
			} else if respStatus >= 400 {
				log.Warnw(ctx, fmt.Sprintf("%s %s?%s", req.Method, req.URL.Path, req.URL.RawQuery), log.TypeKey, log.TypeValRPC, "host", req.Host, "req", reqBody,
					"resp", respBody, "status", respStatus, "request_time", fmt.Sprintf("%.3f", float32(time.Since(begin).Microseconds())/1000))
				return
			}

			log.Infow(ctx, fmt.Sprintf("%s %s?%s", req.Method, req.URL.Path, req.URL.RawQuery), log.TypeKey, log.TypeValRPC, "host", req.Host, "req", reqBody,
				"resp", respBody, "status", respStatus, "request_time", fmt.Sprintf("%.3f", float32(time.Since(begin).Microseconds())/1000))
		}(time.Now())

		resp, err = next.RoundTrip(req)
		return
	})
}

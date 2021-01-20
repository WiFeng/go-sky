package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/WiFeng/go-sky/sky/log"
	kitopentracing "github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/opentracing/opentracing-go"
	opentracingext "github.com/opentracing/opentracing-go/ext"
)

// HTTPResponseWriter ...
type HTTPResponseWriter struct {
	http.ResponseWriter
	statusCode int
	reqBody    []byte
	respBody   []byte
}

func (w *HTTPResponseWriter) Write(b []byte) (int, error) {
	w.respBody = append(w.respBody, b...)
	return w.ResponseWriter.Write(b)
}

// WriteHeader ...
func (w *HTTPResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// HTTPClientDoFunc ...
type HTTPClientDoFunc func(*http.Request) (*http.Response, error)

// Do ...
func (c HTTPClientDoFunc) Do(req *http.Request) (*http.Response, error) {
	return c(req)
}

// HTTPClient ...
type HTTPClient struct {
	*http.Client
	middlewares []HTTPClientMiddleware
}

// Use ...
func (c *HTTPClient) Use(mwf ...HTTPClientMiddlewareFunc) {
	for _, fn := range mwf {
		c.middlewares = append(c.middlewares, fn)
	}
}

// Do ...
func (c HTTPClient) Do(req *http.Request) (*http.Response, error) {
	var cl = kithttp.HTTPClient(c.Client)
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		cl = c.middlewares[i].Middleware(cl)
	}
	return cl.Do(req)
}

// HTTPClientMiddleware ...
type HTTPClientMiddleware interface {
	Middleware(kithttp.HTTPClient) kithttp.HTTPClient
}

// HTTPClientMiddlewareFunc ...
type HTTPClientMiddlewareFunc func(kithttp.HTTPClient) kithttp.HTTPClient

// Middleware allows MiddlewareFunc to implement the middleware interface.
func (mw HTTPClientMiddlewareFunc) Middleware(httpClient kithttp.HTTPClient) kithttp.HTTPClient {
	return mw(httpClient)
}

// HTTPClientDoMiddleware ...
func HTTPClientDoMiddleware(next kithttp.HTTPClient) kithttp.HTTPClient {
	return HTTPClientDoFunc(func(req *http.Request) (*http.Response, error) {
		ctx := req.Context()

		defer func(begin time.Time) {
			log.Infow(ctx, "", "request_time", time.Since(begin).Microseconds())
		}(time.Now())
		return next.Do(req)
	})
}

// ==========================================
// HTTPServer Middleware
// ==========================================

// HTTPServerLoggingMiddleware ...
func HTTPServerLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var reqBodyBytes = make([]byte, 0)
		var respBodyBytes = make([]byte, 0)
		{
			if r.Body != nil {
				reqBodyBytes, _ = ioutil.ReadAll(r.Body)
			}
			r.Body = ioutil.NopCloser(bytes.NewBuffer(respBodyBytes))
		}

		iw := &HTTPResponseWriter{
			w,
			http.StatusOK,
			reqBodyBytes,
			respBodyBytes,
		}

		defer func(begin time.Time) {
			var reqBody string
			var respBody string
			{
				reqBody = string(iw.reqBody)
				respBody = string(iw.respBody)
				if len(reqBody) > 800 {
					reqBody = reqBody[0:800]
				}
				if len(respBody) > 500 {
					respBody = respBody[0:500]
				}
			}

			if span := opentracing.SpanFromContext(ctx); span != nil {
				opentracingext.HTTPStatusCode.Set(span, uint16(iw.statusCode))
			}

			log.Infow(ctx, fmt.Sprintf("%s %s", r.Method, r.RequestURI), log.TypeKey, log.TypeValAccess, "host", r.Host, "req", reqBody,
				"resp", respBody, "status", iw.statusCode, "request_time", fmt.Sprintf("%.3f", float32(time.Since(begin).Microseconds())/1000))
		}(time.Now())

		next.ServeHTTP(iw, r)
	})
}

// HTTPServerTracingMiddleware ...
func HTTPServerTracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx = r.Context()
		var logger = log.LoggerFromContext(ctx)
		var tracer = opentracing.GlobalTracer()
		var operationName = fmt.Sprintf("[%s] %s", r.Method, r.URL.Path)

		ctx = kitopentracing.HTTPToContext(tracer, operationName, logger)(ctx, r)
		ctx = log.BuildLogger(ctx)
		r = r.WithContext(ctx)

		defer func() {
			if span := opentracing.SpanFromContext(ctx); span != nil {
				span.Finish()
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// ==========================================
// HTTPClient Middleware
// ==========================================

// HTTPClientLoggingMiddleware ...
func HTTPClientLoggingMiddleware(next kithttp.HTTPClient) kithttp.HTTPClient {
	return HTTPClientDoFunc(func(req *http.Request) (resp *http.Response, err error) {
		ctx := req.Context()

		var reqBodyBytes = make([]byte, 0)
		var respBodyBytes = make([]byte, 0)
		{
			if req.Body != nil {
				reqBodyBytes, _ = ioutil.ReadAll(req.Body)
				req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBodyBytes))
			}
		}

		defer func(begin time.Time) {
			var reqBody string
			var respBody string
			{
				if resp.Body != nil {
					respBodyBytes, _ = ioutil.ReadAll(resp.Body)
					resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBodyBytes))
				}

				reqBody = string(reqBodyBytes)
				respBody = string(respBodyBytes)
				if len(reqBody) > 800 {
					reqBody = reqBody[0:800]
				}
				if len(respBody) > 500 {
					respBody = respBody[0:500]
				}
			}

			log.Infow(ctx, fmt.Sprintf("%s %s?%s", req.Method, req.URL.Path, req.URL.RawQuery), log.TypeKey, log.TypeValRPC, "host", req.Host, "req", reqBody,
				"resp", respBody, "status", resp.StatusCode, "request_time", fmt.Sprintf("%.3f", float32(time.Since(begin).Microseconds())/1000))
		}(time.Now())

		resp, err = next.Do(req)
		return
	})
}

// HTTPClientTracingMiddleware ...
func HTTPClientTracingMiddleware(next kithttp.HTTPClient) kithttp.HTTPClient {
	return HTTPClientDoFunc(func(req *http.Request) (resp *http.Response, err error) {
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
		resp, err = next.Do(req)
		return
	})
}

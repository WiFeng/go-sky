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
)

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
		var operationName = fmt.Sprintf("[%s]%s", r.Method, r.URL.Path)

		ctx = kitopentracing.HTTPToContext(tracer, operationName, logger)(ctx, r)
		ctx = log.BuildLogger(ctx)
		r = r.WithContext(ctx)

		defer func() {
			span := opentracing.SpanFromContext(ctx)
			span.Finish()
		}()

		next.ServeHTTP(w, r)
	})
}

// HTTPClientLoggingMiddleware ...
func HTTPClientLoggingMiddleware(next kithttp.HTTPClient) kithttp.HTTPClient {
	return ClientDoFunc(func(req *http.Request) (resp *http.Response, err error) {
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
	return ClientDoFunc(func(req *http.Request) (*http.Response, error) {
		var ctx = req.Context()
		var logger = log.LoggerFromContext(ctx)
		var tracer = opentracing.GlobalTracer()
		kitopentracing.ContextToHTTP(tracer, logger)(ctx, req)
		return next.Do(req)
	})
}

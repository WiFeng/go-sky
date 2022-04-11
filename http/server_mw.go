package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/WiFeng/go-sky/log"
	skyprome "github.com/WiFeng/go-sky/metrics/prometheus"
	kitopentracing "github.com/go-kit/kit/tracing/opentracing"
	opentracing "github.com/opentracing/opentracing-go"
	opentracingext "github.com/opentracing/opentracing-go/ext"
)

// ResponseWriter ...
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
	reqBody    []byte
	respBody   []byte
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	w.respBody = append(w.respBody, b...)
	return w.ResponseWriter.Write(b)
}

// WriteHeader ...
func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// ==========================================
// Server Middleware
// ==========================================

// ServerLoggingMiddleware ...
func ServerLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var reqBodyBytes = make([]byte, 0)
		var respBodyBytes = make([]byte, 0)
		{
			if r.Body != nil {
				reqBodyBytes, _ = ioutil.ReadAll(r.Body)
			}
			r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBodyBytes))
		}

		iw := &ResponseWriter{
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

// ServerTracingMiddleware ...
func ServerTracingMiddleware(next http.Handler) http.Handler {
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

// ServerMetricsMiddleware ...
func ServerMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var iw *ResponseWriter
		{
			reqBodyBytes := make([]byte, 0)
			respBodyBytes := make([]byte, 0)

			if r.Body != nil {
				reqBodyBytes, _ = ioutil.ReadAll(r.Body)
			}
			r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBodyBytes))

			iw = &ResponseWriter{
				w,
				http.StatusOK,
				reqBodyBytes,
				respBodyBytes,
			}
		}

		defer func(begin time.Time) {
			duration := float64(time.Since(begin).Microseconds()) / 1000000

			skyprome.HTTPServerRequestsTotalCounter(iw.statusCode, r.Method, r.URL.Path)
			skyprome.HTTPServerRequestsDurationHistogram(iw.statusCode, r.Method, r.URL.Path, duration)
			skyprome.HTTPServerRequestsDurationSummary(iw.statusCode, r.Method, r.URL.Path, duration)
		}(time.Now())

		next.ServeHTTP(iw, r)
	})
}

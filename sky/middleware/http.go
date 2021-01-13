package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/WiFeng/go-sky/sky/log"
	kitopentracing "github.com/go-kit/kit/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
)

type responseWriter struct {
	http.ResponseWriter
	buffer     []byte
	statusCode int
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.buffer = append(w.buffer, b...)
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// HTTPServerLoggingMiddleware ...
func HTTPServerLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		iw := &responseWriter{
			w,
			nil,
			http.StatusOK,
		}

		defer func(begin time.Time) {
			var resp string
			if iw.buffer != nil {
				resp = string(iw.buffer)
				if len(resp) > 500 {
					resp = resp[0:500]
				}
			}
			log.Infow(ctx, fmt.Sprintf("%s %s", r.Method, r.RequestURI),
				"resp", resp, "status", iw.statusCode, "header", iw.Header(),
				"request_time", time.Since(begin).Microseconds())
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
		var operationName = fmt.Sprintf("[%s]%s", r.Method, r.URL)

		ctx = kitopentracing.HTTPToContext(tracer, operationName, logger)(ctx, r)
		ctx = log.BuildLogger(ctx)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

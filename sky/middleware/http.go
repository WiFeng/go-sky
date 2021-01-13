package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/WiFeng/go-sky/sky/log"
	kitopentracing "github.com/go-kit/kit/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	reqBody    []byte
	respBody   []byte
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.respBody = append(w.respBody, b...)
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

		var reqBody = make([]byte, 0)
		var respBody = make([]byte, 0)
		{
			if r.Body != nil {
				reqBody, _ = ioutil.ReadAll(r.Body)
			}
			r.Body = ioutil.NopCloser(bytes.NewBuffer(respBody))
		}

		iw := &responseWriter{
			w,
			http.StatusOK,
			reqBody,
			respBody,
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
				"resp", respBody, "status", iw.statusCode, "request_time", time.Since(begin).Microseconds())
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

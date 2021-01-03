package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/WiFeng/go-sky/sky/config"
	"github.com/WiFeng/go-sky/sky/log"
	"github.com/WiFeng/go-sky/sky/middleware"
	"github.com/oklog/oklog/pkg/group"
	"github.com/opentracing/opentracing-go"

	kitendpoint "github.com/go-kit/kit/endpoint"
	kitopentracing "github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
)

// Server ...
type Server struct {
	*kithttp.Server
}

// NewServer ...
func NewServer(
	e kitendpoint.Endpoint,
	dec kithttp.DecodeRequestFunc,
	enc kithttp.EncodeResponseFunc,
	opt ...kithttp.ServerOption,
) *Server {

	logger := log.GetDefaultLogger()
	options := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(errorEncoder),
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerBefore(beforeHandler),
		kithttp.ServerAfter(afterHandler),
	}

	if opt != nil {
		options = append(options, opt...)
	}

	e = middleware.PanicMiddleware()(e)
	e = middleware.LoggingMiddleware()(e)

	ks := kithttp.NewServer(
		e,
		dec,
		enc,
		options...,
	)

	s := &Server{
		ks,
	}

	return s
}

// ListenAndServe ...
func ListenAndServe(ctx context.Context, conf config.HTTP, httpHandler http.Handler) {

	var g group.Group
	{
		// httpAddr is configurable
		httpAddr := &conf.Addr

		// The HTTP listener mounts the Go kit HTTP handler we created.
		httpListener, err := net.Listen("tcp", *httpAddr)
		if err != nil {
			log.Fatalw(ctx, "listen error", "transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			log.Infow(ctx, "serve start", "transport", "HTTP", "addr", *httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			httpListener.Close()
		})
	}

	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}

	log.Info(ctx, "serve exit. ", g.Run())
}

func errorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	log.Error(ctx, err)
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

type errorWrapper struct {
	Error string `json:"error"`
}

// StartSpan ...
func StartSpan(ctx context.Context, r *http.Request) context.Context {

	var logger = log.LoggerFromContext(ctx)
	var tracer = opentracing.GlobalTracer()
	var operationName = fmt.Sprintf("[%s]%s", r.Method, r.URL)

	return kitopentracing.HTTPToContext(tracer, operationName, logger)(ctx, r)
}

// FinishSpan ...
func FinishSpan(ctx context.Context, w http.ResponseWriter) context.Context {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span.Finish()
	}
	return ctx
}

func beforeHandler(ctx context.Context, r *http.Request) context.Context {
	//ctx = trace.StartSpan(ctx, r)
	ctx = StartSpan(ctx, r)
	ctx = log.BuildLogger(ctx)
	return ctx
}

func afterHandler(ctx context.Context, w http.ResponseWriter) context.Context {
	// ctx = trace.FinishSpan(ctx, w)
	// ctx = syncLogger(ctx)
	ctx = FinishSpan(ctx, w)
	return ctx
}

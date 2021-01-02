package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/WiFeng/go-sky/sky/config"
	"github.com/WiFeng/go-sky/sky/middleware"
	"github.com/WiFeng/go-sky/sky/log"
	"github.com/WiFeng/go-sky/sky/trace"
	"github.com/oklog/oklog/pkg/group"

	kitendpoint "github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
)

type Server struct {
	*kithttp.Server
}

func NewServer(
	e kitendpoint.Endpoint,
	dec kithttp.DecodeRequestFunc,
	enc kithttp.EncodeResponseFunc,
	opt ...kithttp.ServerOption,
) *Server {

	options := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(errorEncoder),
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerBefore(beforeHandler),
		kithttp.ServerAfter(afterHandler),
	}

	if opt != nil {
		options = append(options, opt...)
	}

	e = middleware.LoggingMiddleware()(e)

	ks := kithttp.NewServer(
		e,
		dec,
		enc,
		options...,
	)

	s := &Server{
		ks
	}

	return s
}

// ListenAndServe ...
func ListenAndServe(conf config.HTTP, httpHandler http.Handler) {

	var g group.Group
	{
		// httpAddr is configurable
		httpAddr := &conf.Addr

		// The HTTP listener mounts the Go kit HTTP handler we created.
		httpListener, err := net.Listen("tcp", *httpAddr)
		if err != nil {
			log.Fatalw("listen error", "transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			log.Infow("serve start", "transport", "HTTP", "addr", *httpAddr)
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

	log.Info("serve exit. ", g.Run())
}

func beforeHandler(ctx context.Context, r *http.Request) context.Context {
	ctx = trace.StartSpan(ctx, r)
	ctx = log.BuildLogger(ctx)
	return ctx
}

func afterHandler(ctx context.Context, w http.ResponseWriter) context.Context {
	ctx = trace.FinishSpan(ctx, w)
	// ctx = syncLogger(ctx)
	return ctx
}



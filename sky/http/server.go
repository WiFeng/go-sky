package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WiFeng/go-sky/sky/config"
	"github.com/WiFeng/go-sky/sky/log"
	"github.com/gorilla/mux"
	"github.com/oklog/oklog/pkg/group"

	kitendpoint "github.com/go-kit/kit/endpoint"
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
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
	}

	if opt != nil {
		options = append(options, opt...)
	}

	e = PanicMiddleware()(e)

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

// NewRouter ...
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(HTTPServerTracingMiddleware)
	r.Use(HTTPServerLoggingMiddleware)
	return r
}

// ListenAndServe ...
func ListenAndServe(ctx context.Context, conf config.HTTP, httpHandler http.Handler) {

	var g group.Group
	var s *http.Server
	{

		httpAddr := conf.Addr
		s = &http.Server{
			Addr: httpAddr,
			// Good practice to set timeouts to avoid Slowloris attacks.
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      httpHandler, // Pass our instance of gorilla/mux in.
		}

		g.Add(func() (err error) {
			defer func(ctx context.Context) {
				if err == http.ErrServerClosed {
					return
				}
				log.Fatalw(ctx, "listen error", "transport", "HTTP", "during", "Listen", "err", err)
			}(ctx)

			log.Infow(ctx, "serve start", "transport", "HTTP", "addr", httpAddr)
			err = s.ListenAndServe()
			return
		}, func(err error) {
			log.Info(ctx, "serve prepare shutdown. ", err)
			s.Shutdown(ctx)
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

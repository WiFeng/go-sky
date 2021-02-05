package http

import (
	"net/http"

	"github.com/WiFeng/go-sky/sky-example/pkg/endpoint"
	skyhttp "github.com/WiFeng/go-sky/sky/http"
)

// NewHandler returns an HTTP handler that makes a set of endpoints
// available on predefined paths.
func NewHandler(endpoints endpoint.Endpoints) http.Handler {
	r := skyhttp.NewRouter()

	r.Methods(http.MethodGet).Path("/example/echo").Handler(skyhttp.NewServer(
		endpoints.Example.Echo,
		decodeHTTPExampleEchoRequest,
		encodeHTTPGenericResponse,
	))

	r.Methods(http.MethodGet).Path("/example/rpc").Handler(skyhttp.NewServer(
		endpoints.Example.RPC,
		decodeHTTPExampleEchoRequest,
		encodeHTTPGenericResponse,
	))

	r.Methods(http.MethodGet).Path("/example/trace").Handler(skyhttp.NewServer(
		endpoints.Example.Trace,
		decodeHTTPExampleEchoRequest,
		encodeHTTPGenericResponse,
	))

	return r
}

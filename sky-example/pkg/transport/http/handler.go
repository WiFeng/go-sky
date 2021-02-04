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

	r.Methods(http.MethodGet).Path("/hello/say").Handler(skyhttp.NewServer(
		endpoints.Hello.Say,
		decodeHTTPExampleEchoRequest,
		encodeHTTPGenericResponse,
	))

	return r
}

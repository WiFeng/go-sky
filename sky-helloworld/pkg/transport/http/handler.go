package http

import (
	"net/http"

	"github.com/WiFeng/go-sky/sky-helloworld/pkg/endpoint"
	skyhttp "github.com/WiFeng/go-sky/sky/http"
)

// NewHandler returns an HTTP handler that makes a set of endpoints
// available on predefined paths.
func NewHandler(endpoints endpoint.Endpoints) http.Handler {
	r := skyhttp.NewRouter()

	r.Methods(http.MethodGet).Path("/hello/say").Handler(skyhttp.NewServer(
		endpoints.Hello.Say,
		decodeHTTPHelloSayRequest,
		encodeHTTPGenericResponse,
	))

	r.Methods(http.MethodGet).Path("/hello/say2").Handler(skyhttp.NewServer(
		endpoints.Hello.Say2,
		decodeHTTPHelloSayRequest,
		encodeHTTPGenericResponse,
	))

	return r
}

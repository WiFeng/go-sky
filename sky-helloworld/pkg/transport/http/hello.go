package http

import (
	"context"
	"net/http"

	"github.com/WiFeng/go-sky/sky-helloworld/pkg/service"
)

func decodeHTTPHelloSayRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := service.HelloSayRequest{
		Words: r.URL.Query().Get("words"),
	}
	return req, nil
}

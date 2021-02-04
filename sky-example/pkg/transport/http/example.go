package http

import (
	"context"
	"net/http"

	"github.com/WiFeng/go-sky/sky-example/pkg/service"
)

func decodeHTTPExampleEchoRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := service.ExampleEchoRequest{
		Msg: r.URL.Query().Get("msg"),
	}
	return req, nil
}

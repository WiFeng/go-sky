package helloworld

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	skyhttp "github.com/WiFeng/go-sky/sky/http"
	"github.com/WiFeng/go-sky/sky/log"
)

const (
	helloSayURI  = "/hello/say"
	helloSay2URI = "/hello/say2"
)

// HelloSayRequest ...
type HelloSayRequest struct {
	Words string `json:"words"`
}

// HelloSayResponse ...
type HelloSayResponse struct {
	Words string `json:"words"`
}

// Hello ...
type Hello struct {
}

// NewHello ...
func NewHello() *Hello {
	return &Hello{}
}

func decodeHelloSayResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp HelloSayResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// Say ...
func (*Hello) Say(ctx context.Context, req HelloSayRequest) (HelloSayResponse, error) {
	var resp HelloSayResponse

	url := fmt.Sprintf("%s?words=%s", helloSayURI, req.Words)
	cli, err := skyhttp.NewClient(ctx, serviceName, http.MethodGet, url,
		encodeHTTPGenericRequest, decodeHelloSayResponse)
	if err != nil {
		return resp, err
	}

	result, err := cli.Endpoint()(ctx, req)

	if err != nil {
		log.Errorw(ctx, "hello.say.endpoint error", "err", err)
		return resp, err
	}

	resp = result.(HelloSayResponse)
	return resp, nil

}

// Say2 ...
func (*Hello) Say2(ctx context.Context, req HelloSayRequest) (HelloSayResponse, error) {
	var resp HelloSayResponse

	url := fmt.Sprintf("%s?words=%s", helloSay2URI, req.Words)
	cli, err := skyhttp.NewClient(ctx, serviceName, http.MethodGet, url,
		encodeHTTPGenericRequest, decodeHelloSayResponse)
	if err != nil {
		return resp, err
	}

	result, err := cli.Endpoint()(ctx, req)

	if err != nil {
		log.Errorw(ctx, "hello.say.endpoint error", "err", err)
		return resp, err
	}

	resp = result.(HelloSayResponse)
	return resp, nil

}

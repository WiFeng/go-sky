package endpoint

import (
	"context"

	"github.com/WiFeng/go-sky/sky-helloworld/pkg/service"
	kitendpoint "github.com/go-kit/kit/endpoint"
)

// HelloEndpoints ...
type HelloEndpoints struct {
	Say  kitendpoint.Endpoint
	Say2 kitendpoint.Endpoint
}

// NewHelloEndpoints ...
func NewHelloEndpoints(s service.Service) HelloEndpoints {
	return HelloEndpoints{
		Say:  MakeHelloSayEndpoint(s),
		Say2: MakeHelloSay2Endpoint(s),
	}
}

// MakeHelloSayEndpoint ...
func MakeHelloSayEndpoint(s service.Service) kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(service.HelloSayRequest)
		return s.Hello.Say(ctx, req)
	}
}

// MakeHelloSay2Endpoint ...
func MakeHelloSay2Endpoint(s service.Service) kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(service.HelloSayRequest)
		return s.Hello.Say2(ctx, req)
	}
}

package endpoint

import (
	"context"

	"github.com/WiFeng/go-sky/sky-example/pkg/service"
	kitendpoint "github.com/go-kit/kit/endpoint"
)

// ExampleEndpoints ...
type ExampleEndpoints struct {
	Echo  kitendpoint.Endpoint
	RPC   kitendpoint.Endpoint
	Trace kitendpoint.Endpoint
}

// NewExampleEndpoints ...
func NewExampleEndpoints(s service.Service) ExampleEndpoints {
	return ExampleEndpoints{
		Echo:  MakeExampleEchoEndpoint(s),
		RPC:   MakeExampleRPCEndpoint(s),
		Trace: MakeExampleTraceEndpoint(s),
	}
}

// MakeExampleEchoEndpoint ...
func MakeExampleEchoEndpoint(s service.Service) kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(service.ExampleEchoRequest)
		return s.Example.Echo(ctx, req)
	}
}

// MakeExampleRPCEndpoint ...
func MakeExampleRPCEndpoint(s service.Service) kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(service.ExampleEchoRequest)
		return s.Example.RPC(ctx, req)
	}
}

// MakeExampleTraceEndpoint ...
func MakeExampleTraceEndpoint(s service.Service) kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(service.ExampleEchoRequest)
		return s.Example.Trace(ctx, req)
	}
}

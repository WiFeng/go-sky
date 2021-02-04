package endpoint

import (
	"context"

	kitendpoint "github.com/go-kit/kit/endpoint"
)

// ExampleEndpoints ...
type ExampleEndpoints struct {
	Echo kitendpoint.Endpoint
}

// NewExampleEndpoints ...
func NewExampleEndpoints(s service.Service) ExampleEndpoints {
	return ExampleEndpoints{
		Echo: MakeExampleEchoEndpoint(s),
	}
}

// MakeExampleEchoEndpoint ...
func MakeExampleEchoEndpoint(s service.Service) kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(service.ExampleEchoRequest)
		return s.Example.Echo(ctx, req)
	}
}

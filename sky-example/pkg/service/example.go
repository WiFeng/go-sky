package service

import "context"

// ExampleEchoRequest ...
type ExampleEchoRequest struct {
	Msg string
}

// ExampleEchoResponse ...
type ExampleEchoResponse struct {
	Msg string `json:"msg"`
}

// ExampleService ...
type ExampleService struct {
}

// Echo ...
func (s ExampleService) Echo(ctx context.Context, req ExampleEchoRequest) (interface{}, error) {
	resp := ExampleEchoResponse{
		Msg: req.Msg,
	}
	return resp, nil
}

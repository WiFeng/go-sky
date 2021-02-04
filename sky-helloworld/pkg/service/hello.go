package service

import "context"

// HelloSayRequest ...
type HelloSayRequest struct {
	Words string
}

// HelloSayResponse ...
type HelloSayResponse struct {
	Words string `json:"words"`
}

// HelloService ...
type HelloService struct {
}

// Say ...
func (s HelloService) Say(ctx context.Context, req HelloSayRequest) (interface{}, error) {
	resp := HelloSayResponse{
		Words: req.Words,
	}
	return resp, nil
}

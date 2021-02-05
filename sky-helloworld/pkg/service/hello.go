package service

import (
	"context"
	"time"

	skyredis "github.com/WiFeng/go-sky/sky/redis"
)

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

// Say2 ...
func (s HelloService) Say2(ctx context.Context, req HelloSayRequest) (interface{}, error) {

	redisCli, err := skyredis.GetInstance(ctx, "rdb")
	if err != nil {
		return nil, err
	}

	redisCli.Set(ctx, "__helloworld__:key", "nihao", 30*time.Minute)
	if _, err = redisCli.Get(ctx, "__helloworld__:key").Result(); err != nil {
		return nil, err
	}

	resp := HelloSayResponse{
		Words: req.Words,
	}
	return resp, nil
}

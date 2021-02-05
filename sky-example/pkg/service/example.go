package service

import (
	"context"
	"time"

	"github.com/WiFeng/go-sky/sky/log"

	skydb "github.com/WiFeng/go-sky/sky/database"
	skyes "github.com/WiFeng/go-sky/sky/elasticsearch"
	skyredis "github.com/WiFeng/go-sky/sky/redis"

	"github.com/WiFeng/go-sky/sky-example/pkg/rpc/helloworld"
)

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

// RPC ...
func (s ExampleService) RPC(ctx context.Context, req ExampleEchoRequest) (interface{}, error) {

	rpcReq := helloworld.HelloSayRequest{
		Words: "Ememem",
	}
	rpcResp, err := helloworld.NewHello().Say(ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp := ExampleEchoResponse{
		Msg: req.Msg + " " + rpcResp.Words,
	}
	return resp, nil
}

// Trace ...
func (s ExampleService) Trace(ctx context.Context, req ExampleEchoRequest) (interface{}, error) {

	// ==============================
	// redis
	// ==============================
	redisCli, err := skyredis.GetInstance(ctx, "rdb")
	if err != nil {
		return nil, err
	}

	redisCli.Set(ctx, "__example__:key", "nihao", 30*time.Minute)
	if _, err = redisCli.Get(ctx, "__example__:key").Result(); err != nil {
		return nil, err
	}

	// ==============================
	// elasticsearch
	// ==============================
	es, err := skyes.GetInstance(ctx, "es")
	if err != nil {
		return nil, err
	}
	_, err = es.Search(
		es.Search.WithIndex("test"),
		es.Search.WithSort("timestamp:desc"),
		es.Search.WithSize(1),
		// Annotate this search request; https://www.elastic.co/guide/en/elasticsearch/reference/current/search.html#stats-groups
		es.Search.WithStats("foo"),
		es.Search.WithContext(ctx),
	)
	if err != nil {
		return nil, err
	}

	// ==============================
	// database
	// ==============================
	db, err := skydb.GetInstance(ctx, "db")
	if err != nil {
		return nil, err
	}

	var got string
	err = db.QueryRowContext(ctx, "SELECT 'OK'").Scan(&got)
	log.Errorf(ctx, "db.QueryRowContext error. err:%s", err)

	// ==============================
	// RPC
	// ==============================
	rpcReq := helloworld.HelloSayRequest{
		Words: "Ememem",
	}
	rpcResp, err := helloworld.NewHello().Say2(ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp := ExampleEchoResponse{
		Msg: req.Msg + " " + rpcResp.Words,
	}
	return resp, nil
}

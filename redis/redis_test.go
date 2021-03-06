package redis

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/WiFeng/go-sky/config"
	"github.com/WiFeng/go-sky/log"
)

var (
	testService = "testService"
)

func TestMain(m *testing.M) {
	redisConf := []config.Redis{
		{
			Name: "redis",
			Host: "127.0.0.1",
			Port: 6379,
			Auth: "",
			DB:   0,
		},
	}

	logConf := config.Log{
		Level: "info",
	}
	if _, err := log.Init(context.Background(), testService, logConf); err != nil {
		fmt.Println("Error:", err)
	}

	Init(context.Background(), testService, redisConf)

	os.Exit(m.Run())
}

func TestSet(t *testing.T) {
	var ctx = context.Background()
	redisCli, err := GetInstance(ctx, "redis")
	if err != nil {
		t.Error(err)
	}
	got, err := redisCli.Set(ctx, "__gotest__:set:key1", "val1", 30*time.Minute).Result()
	if err != nil {
		t.Error(err)
	}
	if got != "OK" {
		t.Errorf("redis.Set = %s; want OK", got)
	}
}

func TestGet(t *testing.T) {
	var ctx = context.Background()
	redisCli, err := GetInstance(ctx, "redis")
	if err != nil {
		t.Error(err)
	}

	redisCli.Set(ctx, "__gotest__:get:key2", "val2", 30*time.Minute)
	got, err := redisCli.Get(ctx, "__gotest__:get:key2").Result()
	if err != nil {
		t.Error(err)
	}
	if got != "val2" {
		t.Errorf("redis.Get = %s; want val2", got)
	}
}

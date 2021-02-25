package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/WiFeng/go-sky/config"
	"github.com/WiFeng/go-sky/log"
	"github.com/go-redis/redis/v8"
)

var (
	redisMap    = map[string]*redis.Client{}
	redisConfig = map[string]config.Redis{}
)

var (
	// ErrConfigNotFound ...
	ErrConfigNotFound = errors.New("redis config is not found")
)

// Init ...
func Init(ctx context.Context, serviceName string, cfs []config.Redis) {
	for _, cf := range cfs {
		redisConfig[cf.Name] = cf

		var rdb *redis.Client
		{
			addr := fmt.Sprintf("%s:%d", cf.Host, cf.Port)
			pass := cf.Auth
			db := cf.DB
			rdb = redis.NewClient(&redis.Options{
				Addr:     addr,
				Password: pass,
				DB:       db,
			})

			if _, err := rdb.Ping(ctx).Result(); err != nil {
				log.Fatalw(ctx, "redis ping error", "conf", cf, "err", err)
				continue
			}
		}

		rdb.AddHook(tracingHook{})
		rdb.AddHook(loggingHook{})

		log.Infof(ctx, "Init redis [%s] %+v", cf.Name, cf)
		redisMap[cf.Name] = rdb
	}
}

// GetInstance ...
func GetInstance(ctx context.Context, redisName string) (*redis.Client, error) {
	rdb, ok := redisMap[redisName]
	if !ok {
		err := ErrConfigNotFound
		log.Errorw(ctx, "redis.GetInstance, redisName is not in redisMap map", "redis_name", redisName, "err", err)
		return nil, err
	}
	return rdb, nil
}

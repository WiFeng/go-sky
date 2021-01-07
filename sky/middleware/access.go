package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/WiFeng/go-sky/sky/log"
	kitendpoint "github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
)

// LoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func LoggingMiddleware() kitendpoint.Middleware {
	return func(next kitendpoint.Endpoint) kitendpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				requestMethod := ctx.Value(kithttp.ContextKeyRequestMethod).(string)
				requestPath := ctx.Value(kithttp.ContextKeyRequestPath).(string)

				log.Infow(ctx, fmt.Sprintf("%s %s", requestMethod, requestPath),
					"request_time", time.Since(begin).Microseconds(), "err", err)
			}(time.Now())
			return next(ctx, request)
		}
	}
}

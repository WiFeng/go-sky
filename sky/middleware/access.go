package middleware

import (
	"context"
	"time"

	"github.com/WiFeng/go-sky/sky/log"
	kitendpoint "github.com/go-kit/kit/endpoint"
)

// LoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func LoggingMiddleware() kitendpoint.Middleware {
	return func(next kitendpoint.Endpoint) kitendpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				log.Infow(ctx, "defer caller", "transport_error", err, "took", time.Since(begin).Microseconds())
			}(time.Now())
			return next(ctx, request)
		}
	}
}

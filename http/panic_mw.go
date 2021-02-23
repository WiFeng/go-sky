package http

import (
	"context"
	"errors"

	"github.com/WiFeng/go-sky/log"
	kitendpoint "github.com/go-kit/kit/endpoint"
)

var (
	// ErrPanic ...
	ErrPanic = errors.New("panic error")
)

// PanicMiddleware ...
func PanicMiddleware() kitendpoint.Middleware {
	return func(next kitendpoint.Endpoint) kitendpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func() {
				if panicErr := recover(); panicErr != nil {
					log.Error(ctx, panicErr)
					err = ErrPanic
				}
			}()
			return next(ctx, request)
		}
	}
}

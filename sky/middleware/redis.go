package middleware

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	opentracingext "github.com/opentracing/opentracing-go/ext"
)

// RedisLoggingHook ...
type RedisLoggingHook struct {
}

// BeforeProcess ....
func (r RedisLoggingHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return ctx, nil
}

// AfterProcess ...
func (r RedisLoggingHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	return nil
}

// BeforeProcessPipeline ...
func (r RedisLoggingHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

// AfterProcessPipeline ...
func (r RedisLoggingHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil

}

// RedisTracingHook ...
type RedisTracingHook struct {
}

// BeforeProcess ...
func (r RedisTracingHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {

	var parentSpan opentracing.Span
	var childSpan opentracing.Span

	if parentSpan = opentracing.SpanFromContext(ctx); parentSpan != nil {
		childSpan = parentSpan.Tracer().StartSpan(
			fmt.Sprintf("redis.%s", cmd.Name()),
			opentracing.ChildOf(parentSpan.Context()),

			opentracing.Tag{Key: "cmd.name", Value: cmd.Name()},
			// opentracing.Tag{Key: "cmd.fullname", Value: cmd.FullName()},
			opentracing.Tag{Key: "cmd.string", Value: cmd.String()},
			opentracing.Tag{Key: string(opentracingext.Component), Value: "redis"},
			opentracingext.SpanKindRPCClient,
		)
		ctx = opentracing.ContextWithSpan(ctx, childSpan)
	}

	return ctx, nil
}

// AfterProcess ...
func (r RedisTracingHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if childSpan := opentracing.SpanFromContext(ctx); childSpan != nil {
		childSpan.Finish()
	}

	return nil
}

// BeforeProcessPipeline ...
func (r RedisTracingHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

// AfterProcessPipeline ...
func (r RedisTracingHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil

}

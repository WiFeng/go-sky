package kafka

import (
	"context"
	"errors"

	kafka "github.com/Shopify/sarama"
	"github.com/WiFeng/go-sky/log"
	"github.com/opentracing/opentracing-go"
	opentracingext "github.com/opentracing/opentracing-go/ext"
)

var (
	// ErrDeprecatedMethod ...
	ErrDeprecatedMethod = errors.New("kafka deprecated method")
)

// SyncProducerSendMessage ...
type SyncProducerSendMessage interface {
	Do(ctx context.Context, msg *kafka.ProducerMessage) (partition int32, offset int64, err error)
}

// SyncProducerSendMessageFunc ...
type SyncProducerSendMessageFunc func(ctx context.Context, msg *kafka.ProducerMessage) (partition int32, offset int64, err error)

// Do ...
func (s SyncProducerSendMessageFunc) Do(ctx context.Context, msg *kafka.ProducerMessage) (partition int32, offset int64, err error) {
	return s(ctx, msg)
}

// SyncProducerSendMessageCoreFunc ...
func SyncProducerSendMessageCoreFunc(s kafka.SyncProducer) SyncProducerSendMessage {
	return SyncProducerSendMessageFunc(func(ctx context.Context, msg *kafka.ProducerMessage) (partition int32, offset int64, err error) {
		return s.SendMessage(msg)
	})
}

// SyncProducer ...
type SyncProducer interface {
	kafka.SyncProducer
	Use(ctx context.Context, mwf ...interface{})
	SendMessageContext(ctx context.Context, msg *kafka.ProducerMessage) (partition int32, offset int64, err error)
}

// SyncProducer ...
type syncProducer struct {
	kafka.SyncProducer
	sendMessageMiddlewares []SyncProducerSendMessageMiddleware
}

// Use ...
func (s *syncProducer) Use(ctx context.Context, mwf ...interface{}) {
	for _, fn := range mwf {
		switch fn := fn.(type) {
		case SyncProducerSendMessageMiddlewareFunc:
			s.sendMessageMiddlewares = append(s.sendMessageMiddlewares, fn)
		default:
			log.Errorf(ctx, "syncProducer.Use error. ", "type is not found.")
		}
	}
}

// SendMessage ...
func (s syncProducer) SendMessage(msg *kafka.ProducerMessage) (partition int32, offset int64, err error) {
	return 0, 0, ErrDeprecatedMethod
}

// SendMessage ...
func (s syncProducer) SendMessageContext(ctx context.Context, msg *kafka.ProducerMessage) (partition int32, offset int64, err error) {
	var sp = SyncProducerSendMessageCoreFunc(s.SyncProducer)
	for i := len(s.sendMessageMiddlewares) - 1; i >= 0; i-- {
		sp = s.sendMessageMiddlewares[i].Middleware(sp)
	}
	return sp.Do(ctx, msg)
}

// SyncProducerSendMessageMiddleware ...
type SyncProducerSendMessageMiddleware interface {
	Middleware(SyncProducerSendMessage) SyncProducerSendMessage
}

// SyncProducerSendMessageMiddlewareFunc ...
type SyncProducerSendMessageMiddlewareFunc func(SyncProducerSendMessage) SyncProducerSendMessage

// Middleware allows MiddlewareFunc to implement the middleware interface.
func (mw SyncProducerSendMessageMiddlewareFunc) Middleware(sp SyncProducerSendMessage) SyncProducerSendMessage {
	return mw(sp)
}

// ==========================================
// SyncProducer Middleware
// ==========================================

// SyncProducerSendMessageTracingMiddleware ...
func SyncProducerSendMessageTracingMiddleware(next SyncProducerSendMessage) SyncProducerSendMessage {
	return SyncProducerSendMessageFunc(func(ctx context.Context, msg *kafka.ProducerMessage) (partition int32, offset int64, err error) {
		var parentSpan opentracing.Span
		var childSpan opentracing.Span

		defer func() {
			if childSpan != nil {
				childSpan.Finish()
			}
		}()

		if parentSpan = opentracing.SpanFromContext(ctx); parentSpan != nil {
			childSpan = parentSpan.Tracer().StartSpan(
				"kafka.SyncProducer.SendMessage",
				opentracing.ChildOf(parentSpan.Context()),
				opentracing.Tag{Key: "message.topic", Value: msg.Topic},
				opentracing.Tag{Key: "message.offset", Value: msg.Offset},
				opentracing.Tag{Key: "message.partition", Value: msg.Partition},
				opentracing.Tag{Key: string(opentracingext.Component), Value: "kafka"},
				opentracingext.SpanKindRPCClient,
			)
			ctx = opentracing.ContextWithSpan(ctx, childSpan)
		}

		return next.Do(ctx, msg)
	})
}

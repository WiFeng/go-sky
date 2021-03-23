package kafka

import (
	"context"
	"errors"

	kafka "github.com/Shopify/sarama"
	"github.com/WiFeng/go-sky/config"
	"github.com/WiFeng/go-sky/log"
)

var (
	kafkaMap    = map[string]kafka.Client{}
	kafkaConfig = map[string]config.Kafka{}
)

var (
	// ErrConfigNotFound ...
	ErrConfigNotFound = errors.New("kafka config is not found")
)

// Init ...
func Init(ctx context.Context, serviceName string, cfs []config.Kafka) {
	for _, cf := range cfs {
		kafkaConfig[cf.Name] = cf

		var kcl kafka.Client
		var err error
		{
			kConfig := kafka.NewConfig()
			kcl, err = kafka.NewClient(cf.Addrs, kConfig)

			if err != nil {
				log.Fatalw(ctx, "kafka.NewClient error", "conf", cf, "err", err)
				continue
			}
		}

		log.Infof(ctx, "Init kafka [%s] %+v", cf.Name, cf)
		kafkaMap[cf.Name] = kcl
	}
}

// NewConsumer ...
func NewConsumer(ctx context.Context, name string) (kafka.Consumer, error) {
	kcl, ok := kafkaMap[name]
	if !ok {
		return nil, ErrConfigNotFound
	}

	return kafka.NewConsumerFromClient(kcl)
}

// NewConsumerGroup ...
func NewConsumerGroup(ctx context.Context, name string) (kafka.ConsumerGroup, error) {
	kcl, ok := kafkaMap[name]
	if !ok {
		return nil, ErrConfigNotFound
	}

	kcf, ok := kafkaConfig[name]
	if !ok {
		return nil, ErrConfigNotFound
	}

	return kafka.NewConsumerGroupFromClient(kcf.Consumer.GroupID, kcl)
}

// NewAsyncProducer ...
func NewAsyncProducer(ctx context.Context, name string) (kafka.AsyncProducer, error) {
	kcl, ok := kafkaMap[name]
	if !ok {
		return nil, ErrConfigNotFound
	}

	return kafka.NewAsyncProducerFromClient(kcl)
}

// NewSyncProducer ...
func NewSyncProducer(ctx context.Context, name string) (SyncProducer, error) {
	kcl, ok := kafkaMap[name]
	if !ok {
		return nil, ErrConfigNotFound
	}

	kcl.Config().Producer.Return.Successes = true
	sp, err := kafka.NewSyncProducerFromClient(kcl)
	if err != nil {
		return nil, err
	}

	spp := &syncProducer{
		SyncProducer: sp,
	}

	spp.Use(ctx, SyncProducerSendMessageMiddlewareFunc(SyncProducerSendMessageTracingMiddleware))

	return spp, nil
}

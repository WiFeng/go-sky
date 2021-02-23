package kafka

import (
	"context"
	"fmt"
	"os"
	"testing"

	kafka "github.com/Shopify/sarama"
	"github.com/WiFeng/go-sky/config"
	"github.com/WiFeng/go-sky/log"
)

var (
	testName = "testKafka"
)

func TestMain(m *testing.M) {
	kafkaConf := []config.Kafka{
		{
			Name:  testName,
			Addrs: []string{"localhost:9092"},
		},
	}

	logConf := config.Log{
		Level: "info",
	}
	if _, err := log.Init(context.Background(), logConf); err != nil {
		fmt.Println("Error:", err)
	}

	Init(context.Background(), kafkaConf)

	os.Exit(m.Run())
}

func TestNewAsyncProducer(t *testing.T) {
	_, err := NewAsyncProducer(context.Background(), testName)
	if err != nil {
		t.Error(err)
	}
}

func TestSyncProducerSendMessage(t *testing.T) {
	producer, err := NewSyncProducer(context.Background(), testName)
	if err != nil {
		t.Error(err)
	}

	msg := &kafka.ProducerMessage{
		Topic: "go-test-topic",
		Value: kafka.StringEncoder("vvvvvvvvvvvvv,vvvvvvvvvv"),
	}
	_, _, err = producer.SendMessageContext(context.Background(), msg)
	if err != nil {
		t.Error(err)
	}
}

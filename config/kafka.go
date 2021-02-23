package config

// Kafka ...
type Kafka struct {
	Name         string
	Addrs        []string
	CustomConfig bool
	Producer     KafkaProducer
	Consumer     KafkaConsumer
}

// KafkaConsumer ...
type KafkaConsumer struct {
	GroupID string
}

// KafkaProducer ...
type KafkaProducer struct {
	MaxMessageBytes int
}

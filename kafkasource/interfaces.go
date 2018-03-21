package kafkasource

import "github.com/confluentinc/confluent-kafka-go/kafka"

// KafkaProducer exposes the Produce operation for a Kafka Producer
type KafkaProducer interface {
	Produce(message *kafka.Message, deliveryChannel chan kafka.Event) error
}

// KafkaConsumer describes the consumer behavior
type KafkaConsumer interface {
	SubscribeTopics(topics []string, rebalance kafka.RebalanceCb) error
	Poll(int) kafka.Event
	Assign([]kafka.TopicPartition) error
	Close() error
	Unassign() error
	Events() chan kafka.Event
}

// CreateKafkaConsumer factory method for creating a KafkaConsumer
type CreateKafkaConsumer func(groupID string) (KafkaConsumer, error)

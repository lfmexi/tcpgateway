package config

import (
	"bitbucket.org/challengerdevs/gpsdriver/events"
	"bitbucket.org/challengerdevs/gpsdriver/kafkasource"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func sessionKafkaConsumerConfig() *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers": "kafka1",
		"group.id":          "driver.sessions",
		"auto.offset.reset": "earliest",
	}
}

func sessionKafkaProducerConfig() *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers": "kafka1",
	}
}

var sessionEvtSource events.EventSource

func sessionEventSource() events.EventSource {
	if sessionEvtSource == nil {
		sessionEvtSource = kafkasource.CreateKafkaEventSource(sessionKafkaConsumerConfig(), sessionKafkaProducerConfig())
	}

	return sessionEvtSource
}

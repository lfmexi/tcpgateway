package config

import (
	"bitbucket.org/challengerdevs/gpsdriver/events"
	"bitbucket.org/challengerdevs/gpsdriver/kafkasource"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func packetKafkaConsumerConfig() *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers": "broker",
		"group.id":          "driver.responses",
		"auto.offset.reset": "earliest",
	}
}

func packetKafkaProducerConfig() *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers": "broker",
	}
}

var packetEvtSource events.EventSource

func packetEventSource() events.EventSource {
	if packetEvtSource == nil {
		packetEvtSource = kafkasource.CreateKafkaEventSource(packetKafkaConsumerConfig(), packetKafkaProducerConfig())
	}

	return packetEvtSource
}

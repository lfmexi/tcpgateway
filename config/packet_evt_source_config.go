package config

import (
	"bitbucket.org/challengerdevs/gpsdriver/events"
	"bitbucket.org/challengerdevs/gpsdriver/kafkasource"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func packetKafkaConsumerConfig() *kafka.ConfigMap {
	responsesConfig := configuration.KafkaConsumers["responses"]
	return &kafka.ConfigMap{
		"bootstrap.servers": responsesConfig.Broker,
		"group.id":          responsesConfig.GroupID,
		"auto.offset.reset": responsesConfig.AutoOffsetReset,
	}
}

func packetKafkaProducerConfig() *kafka.ConfigMap {
	packetsConfig := configuration.KafkaProducers["packets"]
	return &kafka.ConfigMap{
		"bootstrap.servers": packetsConfig.Broker,
	}
}

var packetEvtSource events.EventSource

func packetEventSource() events.EventSource {
	if packetEvtSource == nil {
		packetEvtSource = kafkasource.CreateKafkaEventSource(packetKafkaConsumerConfig(), packetKafkaProducerConfig())
	}

	return packetEvtSource
}

package config

import (
	"bitbucket.org/challengerdevs/gpsdriver/events"
	"bitbucket.org/challengerdevs/gpsdriver/kafkasource"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func packetKafkaConsumerConfig() kafkasource.CreateKafkaConsumer {
	responsesConfig := configuration.KafkaConsumers["responses"]
	configMap := &kafka.ConfigMap{
		"bootstrap.servers":        responsesConfig.Broker,
		"go.events.channel.enable": true,
		"default.topic.config":     kafka.ConfigMap{"auto.offset.reset": responsesConfig.AutoOffsetReset},
	}

	factory := func(groupID string) (kafkasource.KafkaConsumer, error) {
		configMap.SetKey("group.id", groupID)
		return kafka.NewConsumer(configMap)
	}

	return factory
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
		producer, err := kafka.NewProducer(packetKafkaProducerConfig())

		if err != nil {
			panic(err)
		}

		packetEvtSource = kafkasource.CreateKafkaEventSource(packetKafkaConsumerConfig(), producer)
	}

	return packetEvtSource
}

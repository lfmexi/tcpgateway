package config

import (
	"github.com/lfmexi/tcpgateway/events"
	"github.com/lfmexi/tcpgateway/kafkasource"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func createConsumerFactory() kafkasource.CreateKafkaConsumer {
	return func(groupID string) (kafkasource.KafkaConsumer, error) {
		responsesConfig := configuration.KafkaConsumers["responses"]
		configMap := &kafka.ConfigMap{
			"bootstrap.servers":        responsesConfig.Broker,
			"go.events.channel.enable": true,
			"group.id":                 groupID,
			"default.topic.config":     kafka.ConfigMap{"auto.offset.reset": responsesConfig.AutoOffsetReset},
		}

		return kafka.NewConsumer(configMap)
	}
}

func packetKafkaConsumerConfig() kafkasource.CreateKafkaConsumer {
	return createConsumerFactory()
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

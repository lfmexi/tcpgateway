package config

import (
	"bitbucket.org/challengerdevs/tcpgateway/events"
	"bitbucket.org/challengerdevs/tcpgateway/kafkasource"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jinzhu/copier"
)

func createConsumerFactory(consumerConfig *kafka.ConfigMap) kafkasource.CreateKafkaConsumer {
	return func(groupID string) (kafkasource.KafkaConsumer, error) {
		var config kafka.ConfigMap

		groupIDConfig := kafka.ConfigMap{
			"group.id": groupID,
		}

		copier.Copy(&config, consumerConfig)
		copier.Copy(&config, &groupIDConfig)
		return kafka.NewConsumer(&config)
	}
}

func packetKafkaConsumerConfig() kafkasource.CreateKafkaConsumer {
	responsesConfig := configuration.KafkaConsumers["responses"]
	configMap := &kafka.ConfigMap{
		"bootstrap.servers":        responsesConfig.Broker,
		"go.events.channel.enable": true,
		"default.topic.config":     kafka.ConfigMap{"auto.offset.reset": responsesConfig.AutoOffsetReset},
	}

	return createConsumerFactory(configMap)
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

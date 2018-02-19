package config

import (
	"bitbucket.org/challengerdevs/gpsdriver/amqpsource"
	"bitbucket.org/challengerdevs/gpsdriver/events"
)

func packetConsumerExchangeConfig() *amqpsource.AmqpExchangeConfig {
	return &amqpsource.AmqpExchangeConfig{
		Name:        "packet.responses",
		Type:        "direct",
		Durable:     true,
		AutoDeleted: false,
		Internal:    false,
		NoWait:      false,
		Arguments:   nil,
	}
}

func packetConsumerQueueOptions() *amqpsource.AmqpQueueOptions {
	return &amqpsource.AmqpQueueOptions{
		Name:            "",
		Durable:         false,
		DeletedWhenUsed: false,
		Exclusive:       false,
		NoWait:          false,
		Arguments:       nil,
	}
}

func packetConsumerConfig() *amqpsource.AmqpConsumerConfig {
	return &amqpsource.AmqpConsumerConfig{
		ConsumerExchangeConfig: packetConsumerExchangeConfig(),
		ConsumerQueueOptions:   packetConsumerQueueOptions(),
		Name:                   "",
		Exclusive:              false,
		NoLocal:                false,
		NoWait:                 false,
		Arguments:              nil,
	}
}

func packetPublisherExchangeConfig() *amqpsource.AmqpExchangeConfig {
	return &amqpsource.AmqpExchangeConfig{
		Name:        "packet.events",
		Type:        "topic",
		Durable:     true,
		AutoDeleted: false,
		Internal:    false,
		NoWait:      false,
		Arguments:   nil,
	}
}

func packetPublisherConfig() *amqpsource.AmqpPublisherConfig {
	return &amqpsource.AmqpPublisherConfig{
		PublisherExchangeConfig: packetPublisherExchangeConfig(),
	}
}

var packetEvtSource events.EventSource

func packetEventSource() events.EventSource {
	if packetEvtSource == nil {
		consumerConfig := packetConsumerConfig()
		publisherConfig := packetPublisherConfig()
		conn := amqpConnection()

		packetEvtSource = amqpsource.CreateEventSource(conn, consumerConfig, publisherConfig)
	}

	return packetEvtSource
}

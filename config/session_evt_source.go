package config

import (
	"bitbucket.org/challengerdevs/gpsdriver/amqpsource"
	"bitbucket.org/challengerdevs/gpsdriver/events"
)

func sessionConsumerExchangeConfig() *amqpsource.AmqpExchangeConfig {
	return &amqpsource.AmqpExchangeConfig{
		Name:        "sessions",
		Type:        "direct",
		Durable:     true,
		AutoDeleted: false,
		Internal:    false,
		NoWait:      false,
		Arguments:   nil,
	}
}

func sessionConsumerQueueOptions() *amqpsource.AmqpQueueOptions {
	return &amqpsource.AmqpQueueOptions{
		Name:            "",
		Durable:         false,
		DeletedWhenUsed: false,
		Exclusive:       false,
		NoWait:          false,
		Arguments:       nil,
	}
}

func sessionConsumerConfig() *amqpsource.AmqpConsumerConfig {
	return &amqpsource.AmqpConsumerConfig{
		ConsumerExchangeConfig: sessionConsumerExchangeConfig(),
		ConsumerQueueOptions:   sessionConsumerQueueOptions(),
		Name:                   "",
		Exclusive:              false,
		NoLocal:                false,
		NoWait:                 false,
		Arguments:              nil,
	}
}

func sessionPublisherExchangeConfig() *amqpsource.AmqpExchangeConfig {
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

func sessionPublisherConfig() *amqpsource.AmqpPublisherConfig {
	return &amqpsource.AmqpPublisherConfig{
		PublisherExchangeConfig: sessionPublisherExchangeConfig(),
	}
}

var sessionEvtSource events.EventSource

func sessionEventSource() events.EventSource {
	if sessionEvtSource == nil {
		consumerConfig := sessionConsumerConfig()
		publisherConfig := sessionPublisherConfig()
		conn := amqpConnection()

		sessionEvtSource = amqpsource.CreateEventSource(conn, consumerConfig, publisherConfig)
	}

	return sessionEvtSource
}

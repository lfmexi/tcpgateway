package amqpsource

import "github.com/streadway/amqp"

type AmqpExchangeConfig struct {
	Name        string
	Type        string
	Durable     bool
	AutoDeleted bool
	Internal    bool
	NoWait      bool
	Arguments   amqp.Table
}

type AmqpQueueOptions struct {
	Name            string
	Durable         bool
	DeletedWhenUsed bool
	Exclusive       bool
	NoWait          bool
	Arguments       amqp.Table
	RoutingKey      string
}

type AmqpConsumerConfig struct {
	ConsumerExchangeConfig *AmqpExchangeConfig
	ConsumerQueueOptions   *AmqpQueueOptions
	Name                   string
	Exclusive              bool
	NoLocal                bool
	NoWait                 bool
	Arguments              amqp.Table
}

type AmqpPublisherConfig struct {
	PublisherExchangeConfig *AmqpExchangeConfig
}

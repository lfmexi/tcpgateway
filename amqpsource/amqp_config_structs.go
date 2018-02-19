package amqpsource

import "github.com/streadway/amqp"

// AmqpExchangeConfig provides a configuration schema for amqp exchanges
type AmqpExchangeConfig struct {
	Name        string     // the name of the exchange
	Type        string     // the type of the exchange
	Durable     bool       // sets if the exchange is durable or not
	AutoDeleted bool       // sets if the exchange should be autodeleted
	Internal    bool       // sets if the exchange should be of internal use only
	NoWait      bool       // no wait
	Arguments   amqp.Table // Additional arguments for the exchange
}

// AmqpQueueOptions provides a configuration schema for amqp queues
type AmqpQueueOptions struct {
	Name            string     // the name of the queue
	Durable         bool       // sets if the queue should be durable
	DeletedWhenUsed bool       // sets if the queue should be deleted when used
	Exclusive       bool       // sets if the queue should be exclusive
	NoWait          bool       // no wait
	Arguments       amqp.Table // additional arguments for the queue
	RoutingKey      string     // the routing key of the queue
}

// AmqpConsumerConfig provides a configuration schema for consumers
type AmqpConsumerConfig struct {
	ConsumerExchangeConfig *AmqpExchangeConfig // config of the exchange to consume
	ConsumerQueueOptions   *AmqpQueueOptions   // config of the queue to consume
	Name                   string              // name of the consumer
	Exclusive              bool                // sets if the consumer should be exclusive or not
	NoLocal                bool                // sets if the consumer should be local or not
	NoWait                 bool                // no wait
	Arguments              amqp.Table          // additional arguments for the consumer
}

// AmqpPublisherConfig provides a configuration schema for publishers
type AmqpPublisherConfig struct {
	PublisherExchangeConfig *AmqpExchangeConfig // config of the exchange where to publish
}

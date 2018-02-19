package amqpsource

import (
	"bitbucket.org/challengerdevs/gpsdriver/events"
	"github.com/streadway/amqp"
)

func CreateEventSource(conn *amqp.Connection, cc *AmqpConsumerConfig, pc *AmqpPublisherConfig) events.EventSource {
	return &amqpEventSource{
		conn,
		cc,
		pc,
		make(map[string]chan bool),
	}
}

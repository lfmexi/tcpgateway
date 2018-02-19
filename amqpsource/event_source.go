package amqpsource

import (
	"fmt"

	"bitbucket.org/challengerdevs/gpsdriver/events"
	"github.com/streadway/amqp"
)

type amqpEvent struct {
	data []byte
}

func (a amqpEvent) Data() []byte {
	return a.data
}

type amqpEventSource struct {
	conn             *amqp.Connection
	consumerConfig   *AmqpConsumerConfig
	publisherConfig  *AmqpPublisherConfig
	consumersControl map[string]chan bool
}

func (a *amqpEventSource) Publish(key string, data []byte) error {
	ch, err := a.setUpPublisher()

	defer ch.Close()

	if err != nil {
		return err
	}

	destination := a.publisherConfig.PublisherExchangeConfig.Name

	return ch.Publish(
		destination,
		key,
		false,
		false,
		amqp.Publishing{
			Body: data,
		},
	)
}

func (a *amqpEventSource) setUpPublisher() (*amqp.Channel, error) {
	ch, err := a.conn.Channel()

	if err != nil {
		return nil, err
	}

	if err := a.declarePublisherExchange(ch); err != nil {
		return nil, err
	}

	return ch, nil
}

func (a *amqpEventSource) declarePublisherExchange(ch *amqp.Channel) error {
	config := a.publisherConfig.PublisherExchangeConfig
	return ch.ExchangeDeclare(
		config.Name,
		config.Type,
		config.Durable,
		config.AutoDeleted,
		config.Internal,
		config.NoWait,
		config.Arguments,
	)
}

func (a *amqpEventSource) Consume(key string) (<-chan events.Event, error) {
	consumerChannel, amqpChannel, err := a.setUpConsumer(key)

	if err != nil {
		return nil, err
	}

	resultChannel := make(chan events.Event)

	stopChannel := make(chan bool)

	a.consumersControl[key] = stopChannel

	// This will generate a new goroutine
	go consumeForEvent(amqpChannel, consumerChannel, resultChannel, stopChannel)

	return resultChannel, nil
}

func (a *amqpEventSource) Stop(key string) error {
	consumerControl, ok := a.consumersControl[key]

	if !ok {
		return fmt.Errorf("Consumer %s does not exist", key)
	}

	close(consumerControl)

	delete(a.consumersControl, key)

	return nil
}

func consumeForEvent(ch *amqp.Channel, origin <-chan amqp.Delivery, destination chan events.Event, stop <-chan bool) {
loop:
	for {
		select {
		case message := <-origin:
			{
				event := amqpEvent{
					message.Body,
				}

				destination <- event
			}
		case <-stop:
			{
				ch.Close()
				close(destination)
				break loop
			}
		}
	}
}

func (a *amqpEventSource) setUpConsumer(key string) (<-chan amqp.Delivery, *amqp.Channel, error) {
	ch, err := a.conn.Channel()

	if err != nil {
		return nil, nil, err
	}

	if err := a.declareConsumerExchange(ch); err != nil {
		ch.Close()
		return nil, nil, err
	}

	q, err := a.declareQueue(ch, key)

	if err != nil {
		ch.Close()
		return nil, nil, err
	}

	if err := a.bindQueue(ch, q, key); err != nil {
		ch.Close()
		return nil, nil, err
	}

	consumer, err := a.createConsumer(ch, q)

	if err != nil {
		ch.Close()
		return nil, nil, err
	}

	return consumer, ch, nil
}

func (a *amqpEventSource) createConsumer(ch *amqp.Channel, q amqp.Queue) (<-chan amqp.Delivery, error) {
	consumerConfig := a.consumerConfig
	return ch.Consume(
		q.Name,
		consumerConfig.Name,
		true,
		consumerConfig.Exclusive,
		consumerConfig.NoLocal,
		consumerConfig.NoWait,
		consumerConfig.Arguments,
	)
}

func (a *amqpEventSource) declareConsumerExchange(ch *amqp.Channel) error {
	config := a.consumerConfig.ConsumerExchangeConfig
	return ch.ExchangeDeclare(
		config.Name,
		config.Type,
		config.Durable,
		config.AutoDeleted,
		config.Internal,
		config.NoWait,
		config.Arguments,
	)
}

func (a *amqpEventSource) declareQueue(ch *amqp.Channel, key string) (amqp.Queue, error) {
	config := a.consumerConfig.ConsumerQueueOptions

	name := config.Name

	if name == "" {
		name = key
	}

	return ch.QueueDeclare(
		name,
		config.Durable,
		config.DeletedWhenUsed,
		config.Exclusive,
		config.NoWait,
		config.Arguments,
	)
}

func (a *amqpEventSource) bindQueue(ch *amqp.Channel, q amqp.Queue, key string) error {
	return ch.QueueBind(q.Name, key, a.consumerConfig.ConsumerExchangeConfig.Name, false, nil)
}

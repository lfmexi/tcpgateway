package kafkasource

import (
	"fmt"
	"os"

	"bitbucket.org/challengerdevs/tcpgateway/events"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type kafkaEventsouce struct {
	createConsumer   CreateKafkaConsumer
	producer         KafkaProducer
	consumersControl map[string]chan bool
}

// CreateKafkaEventSource creates an EventSource for kafka
func CreateKafkaEventSource(createConsumer CreateKafkaConsumer, producer KafkaProducer) events.EventSource {
	return &kafkaEventsouce{
		createConsumer,
		producer,
		make(map[string]chan bool),
	}
}

func (es *kafkaEventsouce) Publish(destination string, key string, data []byte) error {
	deliveryChannel := make(chan kafka.Event)

	defer func() {
		close(deliveryChannel)
	}()

	kafkaMessage := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &destination,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(key),
		Value: data,
	}

	err := es.producer.Produce(kafkaMessage, deliveryChannel)

	e := <-deliveryChannel
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	return err
}

func (es *kafkaEventsouce) Consume(key string) (<-chan events.Event, error) {
	consumer, err := es.createConsumer(key)

	if err != nil {
		return nil, err
	}

	consumer.SubscribeTopics([]string{key}, nil)

	stopChannel := make(chan bool)

	es.consumersControl[key] = stopChannel

	eventChannel := make(chan events.Event, 10)

	go func() {
	loop:
		for {
			select {
			case <-stopChannel:
				break loop
			case ev := <-consumer.Events():
				switch e := ev.(type) {
				case kafka.AssignedPartitions:
					fmt.Fprintf(os.Stderr, "%% %v\n", e)
					consumer.Assign(e.Partitions)
				case kafka.RevokedPartitions:
					fmt.Fprintf(os.Stderr, "%% %v\n", e)
					consumer.Unassign()
				case *kafka.Message:
					eventChannel <- &kafkaEvent{e.Value}
				case kafka.PartitionEOF:
					fmt.Printf("%% Reached %v\n", e)
				case kafka.Error:
					fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
					es.Stop(key)
				}
			}
		}

		consumer.Close()
		close(stopChannel)
		close(eventChannel)
	}()

	return eventChannel, nil
}

func (es *kafkaEventsouce) Stop(key string) error {
	consumerControl, ok := es.consumersControl[key]

	if !ok {
		return fmt.Errorf("Consumer %s does not exist", key)
	}

	consumerControl <- true

	delete(es.consumersControl, key)
	return nil
}

type kafkaEvent struct {
	data []byte
}

func (k kafkaEvent) Data() []byte {
	return k.data
}

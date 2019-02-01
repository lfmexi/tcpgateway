package kafkasource

import (
	"fmt"
	"log"
	"os"

	"github.com/lfmexi/tcpgateway/events"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type kafkaEventsouce struct {
	createConsumer       CreateKafkaConsumer
	producer             KafkaProducer
	consumersControlChan chan string
}

// CreateKafkaEventSource creates an EventSource for kafka
func CreateKafkaEventSource(createConsumer CreateKafkaConsumer, producer KafkaProducer) events.EventSource {
	return &kafkaEventsouce{
		createConsumer,
		producer,
		make(chan string),
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

	eventChannel := make(chan events.Event, 10)

	go func() {
		stopChannel := make(chan bool)

	loop:

		for {
			select {
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
				case kafka.Error:
					fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				}
			case consumerKey := <-es.consumersControlChan:
				log.Printf("Removing consumer for %s", consumerKey)
				break loop
			case <-stopChannel:
				break loop
			}
		}

		consumer.Close()
		close(stopChannel)
		close(eventChannel)
	}()

	return eventChannel, nil
}

func (es *kafkaEventsouce) Stop(key string) error {
	es.consumersControlChan <- key
	return nil
}

type kafkaEvent struct {
	data []byte
}

func (k kafkaEvent) Data() []byte {
	return k.data
}

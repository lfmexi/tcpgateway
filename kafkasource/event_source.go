package kafkasource

import (
	"fmt"
	"os"

	"bitbucket.org/challengerdevs/gpsdriver/events"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type kafkaEventsouce struct {
	consumerConfig   *kafka.ConfigMap
	producerConfig   *kafka.ConfigMap
	consumersControl map[string]chan bool
}

// CreateKafkaEventSource creates an EventSource for kafka
func CreateKafkaEventSource(consumerConfig *kafka.ConfigMap, producerConfig *kafka.ConfigMap) events.EventSource {
	return &kafkaEventsouce{
		consumerConfig,
		producerConfig,
		make(map[string]chan bool),
	}
}

func (es *kafkaEventsouce) Publish(destination string, key string, data []byte) error {
	producer, err := kafka.NewProducer(es.producerConfig)

	if err != nil {
		return err
	}

	deliveryChannel := make(chan kafka.Event)

	defer func() {
		close(deliveryChannel)
	}()

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &destination,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(key),
		Value: data,
	}, deliveryChannel)

	e := <-deliveryChannel
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	return nil
}

func (es *kafkaEventsouce) Consume(key string) (<-chan events.Event, error) {
	consumer, err := kafka.NewConsumer(es.consumerConfig)

	if err != nil {
		return nil, err
	}

	consumer.SubscribeTopics([]string{key}, nil)

	stopChannel := make(chan bool)

	es.consumersControl[key] = stopChannel

	eventChannel := make(chan events.Event)

	go func() {
	loop:
		for {
			select {
			case <-stopChannel:
				consumer.Close()
				close(eventChannel)
				break loop
			default:
				event := consumer.Poll(100)

				if event == nil {
					continue
				}

				switch e := event.(type) {
				case *kafka.Message:
					fmt.Printf("%% Message on %s\n", e.TopicPartition)
					eventChannel <- &kafkaEvent{e.Value}
				case kafka.PartitionEOF:
					fmt.Printf("%% Reached %v\n", e)
				case kafka.Error:
					fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
					stopChannel <- true
				default:
					fmt.Printf("Ignored %v\n", e)
				}
			}
		}
	}()

	return eventChannel, nil
}

func (es *kafkaEventsouce) Stop(key string) error {
	consumerControl, ok := es.consumersControl[key]

	if !ok {
		return fmt.Errorf("Consumer %s does not exist", key)
	}

	close(consumerControl)

	delete(es.consumersControl, key)
	return nil
}

type kafkaEvent struct {
	data []byte
}

func (k kafkaEvent) Data() []byte {
	return k.data
}

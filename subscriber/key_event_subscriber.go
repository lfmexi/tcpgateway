package subscriber

import (
	"fmt"

	"github.com/lfmexi/tcpgateway/events"
)

// NewKeyBasedEventSubscriberFactory creates a new event subscriber factory based on keys
func NewKeyBasedEventSubscriberFactory(eventSource events.EventSource) events.EventSubscriberFactory {
	return &keyBasedEventSubscriberFactory{
		eventSource,
	}
}

type keyBasedEventSubscriberFactory struct {
	eventSource events.EventSource
}

func (factory *keyBasedEventSubscriberFactory) CreateEventSubscriber() events.EventSubscriber {
	return &keyBasedEventSubscriber{
		factory.eventSource,
		"",
	}
}

type keyBasedEventSubscriber struct {
	eventSource events.EventSource
	key         string
}

func (s *keyBasedEventSubscriber) Observe(args ...string) (<-chan events.Event, error) {
	s.key = args[0]

	if s.key == "" {
		return nil, fmt.Errorf("A key was expected")
	}

	return s.eventSource.Consume(s.key)
}

func (s *keyBasedEventSubscriber) Stop() error {
	return s.eventSource.Stop(s.key)
}

package session

import (
	"log"

	"bitbucket.org/challengerdevs/gpsdriver/events"
)

type Service interface {
	CreateSession(string, []byte) (*Session, error)
}

type EventSessionService struct {
	eventEmitter           events.EventEmitter
	eventSubscriberFactory events.EventSubscriberFactory
}

func NewEventSessionService(eventEmitter events.EventEmitter, eventSubscriberFactory events.EventSubscriberFactory) *EventSessionService {
	return &EventSessionService{
		eventEmitter,
		eventSubscriberFactory,
	}
}

func (s EventSessionService) CreateSession(sessionAddress string, payload []byte) (*Session, error) {

	eventObserver := s.eventSubscriberFactory.CreateEventSubscriber()

	eventChannel, err := eventObserver.Observe(sessionAddress)

	if err != nil {
		return nil, err
	}

	if err := s.eventEmitter.Emit(sessionAddress, payload); err != nil {
		return nil, err
	}

	sessionAckEvent := <-eventChannel

	log.Printf("Session established for %s", sessionAddress)

	if err := eventObserver.Stop(); err != nil {
		return nil, err
	}

	return newSession(sessionAckEvent)
}

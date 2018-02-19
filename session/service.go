package session

import (
	"log"

	"bitbucket.org/challengerdevs/gpsdriver/events"
)

// Service interface that represents the session service
type Service interface {
	CreateSession(string, []byte) (*Session, error)
}

// NewEventSessionService creates a new event based session service
func NewEventSessionService(eventEmitter events.EventEmitter, eventSubscriberFactory events.EventSubscriberFactory) Service {
	return &eventSessionService{
		eventEmitter,
		eventSubscriberFactory,
	}
}

type eventSessionService struct {
	eventEmitter           events.EventEmitter
	eventSubscriberFactory events.EventSubscriberFactory
}

func (s eventSessionService) CreateSession(sessionAddress string, payload []byte) (*Session, error) {

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

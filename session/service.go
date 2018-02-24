package session

import (
	"encoding/json"
	"fmt"
	"log"

	"bitbucket.org/challengerdevs/gpsdriver/events"
)

// Service interface that represents the session service
type Service interface {
	CreateSession(string, string, []byte) (*Session, error)
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

func (s eventSessionService) CreateSession(sessionAddress string, deviceType string, payload []byte) (*Session, error) {
	eventObserver := s.eventSubscriberFactory.CreateEventSubscriber()

	eventChannel, err := eventObserver.Observe(sessionAddress)

	if err != nil {
		return nil, err
	}

	if err := s.eventEmitter.Emit("login."+deviceType, payload); err != nil {
		return nil, err
	}

	sessionAckEvent := <-eventChannel

	log.Printf("Session established for %s", sessionAddress)

	if err := eventObserver.Stop(); err != nil {
		return nil, err
	}

	return s.newSession(sessionAckEvent)
}

func (service *eventSessionService) newSession(event events.Event) (*Session, error) {
	sEvt := sessionEvent{}

	if err := json.Unmarshal(event.Data(), &sEvt); err != nil {
		return nil, err
	}

	if sEvt.EventType != "sessionAck" {
		return nil, fmt.Errorf("Session event was expected")
	}

	s := Session{}

	s.SessionID = sEvt.SessionID
	s.SessionAckPacket = sEvt.SessionAckPacket
	s.Disconnected = make(chan bool)
	s.eventEmitter = service.eventEmitter

	return &s, nil
}

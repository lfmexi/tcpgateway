package session

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"bitbucket.org/challengerdevs/gpsdriver/events"
)

// Service interface that represents the session service
type Service interface {
	CreateSession(string, string) (*Session, error)
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

func (s eventSessionService) CreateSession(sessionAddress string, deviceType string) (*Session, error) {
	eventObserver := s.eventSubscriberFactory.CreateEventSubscriber()

	eventChannel, err := eventObserver.Observe(sessionAddress)

	if err != nil {
		return nil, err
	}

	if err := s.eventEmitter.Emit("devices.login", sessionAddress, []byte(fmt.Sprintf("{\"device_type\":\"%s\", \"address\":\"%s\"}", deviceType, sessionAddress))); err != nil {
		return nil, err
	}

	timeoutChan := make(chan bool)

	go func() {
		time.Sleep(30 * time.Second)
		timeoutChan <- true
	}()

	err = nil
	select {
	case sessionAckEvent := <-eventChannel:
		log.Printf("Session established for %s", sessionAddress)

		if err := eventObserver.Stop(); err != nil {
			return nil, err
		}

		return s.newSession(sessionAckEvent)
	case <-timeoutChan:
		err = fmt.Errorf("Timeout exceeded while wating to establish the session for %s", sessionAddress)
	}

	return nil, err
}

func (s *eventSessionService) newSession(event events.Event) (*Session, error) {
	sEvt := sessionEvent{}

	if err := json.Unmarshal(event.Data(), &sEvt); err != nil {
		return nil, err
	}

	if sEvt.EventType != "sessionAck" {
		return nil, fmt.Errorf("Session event was expected")
	}

	session := Session{}

	session.SessionID = sEvt.SessionID
	session.Disconnected = make(chan bool)
	session.eventEmitter = s.eventEmitter

	return &session, nil
}

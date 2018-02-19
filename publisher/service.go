package publisher

import (
	"bitbucket.org/challengerdevs/gpsdriver/events"
)

type Service interface {
	Publish(string, []byte) error
}

type EventPublisherService struct {
	traceEventEmitter events.EventEmitter
}

func NewEventPublisherService(traceEventEmitter events.EventEmitter) *EventPublisherService {
	return &EventPublisherService{
		traceEventEmitter,
	}
}

func (es EventPublisherService) Publish(sessionID string, data []byte) error {
	return es.traceEventEmitter.Emit(sessionID, data)
}

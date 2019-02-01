package publisher

import (
	"github.com/lfmexi/tcpgateway/events"
)

// Service interface that represents a publisher service
type Service interface {
	Publish(string, string, []byte) error
}

// NewEventPublisherService creates a new event publisher service
func NewEventPublisherService(traceEventEmitter events.EventEmitter) Service {
	return &eventPublisherService{
		traceEventEmitter,
	}
}

type eventPublisherService struct {
	traceEventEmitter events.EventEmitter
}

func (es eventPublisherService) Publish(destination string, sessionID string, data []byte) error {
	return es.traceEventEmitter.Emit(destination, sessionID, data)
}

package session

import (
	"bitbucket.org/challengerdevs/gpsdriver/events"
)

// NewSessionEventEmitter creates a new Event Emitter based on sessions
func NewSessionEventEmitter(s events.EventSource) events.EventEmitter {
	return &sessionEventEmitter{
		s,
	}
}

type sessionEventEmitter struct {
	source events.EventSource
}

func (s *sessionEventEmitter) Emit(key string, data []byte) error {
	return s.source.Publish(key, data)
}

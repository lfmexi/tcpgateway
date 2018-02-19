package session

import (
	"bitbucket.org/challengerdevs/gpsdriver/events"
)

type SessionEventEmitter struct {
	source events.EventSource
}

func (s *SessionEventEmitter) Emit(key string, data []byte) error {
	return s.source.Publish(key, data)
}

func NewSessionEventEmitter(s events.EventSource) events.EventEmitter {
	return &SessionEventEmitter{
		s,
	}
}

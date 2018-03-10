package session

import (
	"fmt"

	"bitbucket.org/challengerdevs/gpsdriver/events"
)

// Session represents an initialized session
type Session struct {
	SessionID    string
	Disconnected chan bool
	eventEmitter events.EventEmitter
}

type sessionEvent struct {
	EventType string `json:"event_type"`
	SessionID string `json:"session_id"`
}

// CloseSession sends a logout event for closing a session
func (s *Session) CloseSession() error {
	return s.eventEmitter.Emit("devices.logout", "", []byte(fmt.Sprintf("{\"id\":\"%s\"}", s.SessionID)))
}

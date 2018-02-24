package session

import (
	"encoding/json"

	"bitbucket.org/challengerdevs/gpsdriver/events"
)

// Session represents an initialized session
type Session struct {
	SessionID        string
	SessionAckPacket []byte
	Disconnected     chan bool
	eventEmitter     events.EventEmitter
}

type sessionEvent struct {
	EventType        string          `json:"event_name"`
	SessionID        string          `json:"session_id"`
	SessionAckPacket json.RawMessage `json:"ack_packet"`
}

func (s *Session) CloseSession() error {
	return s.eventEmitter.Emit("logout."+s.SessionID, nil)
}

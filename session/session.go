package session

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/challengerdevs/gpsdriver/events"
)

// Session represents an initialized session
type Session struct {
	SessionID        string
	SessionAckPacket []byte
	Disconnected     chan bool
}

type sessionEvent struct {
	EventType        string          `json:"event_name"`
	SessionID        string          `json:"device_id"`
	SessionAckPacket json.RawMessage `json:"ack_packet"`
}

func newSession(event events.Event) (*Session, error) {
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

	return &s, nil
}

package session

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/challengerdevs/gpsdriver/events"
)

type Session struct {
	EventType        string          `json:"event_name"`
	SessionID        string          `json:"device_id"`
	SessionAckPacket json.RawMessage `json:"ack_packet"`
	Disconnected     chan bool
}

func newSession(event events.Event) (*Session, error) {
	s := Session{}

	err := json.Unmarshal(event.Data(), &s)

	if s.EventType != "sessionAck" {
		return nil, fmt.Errorf("Session event was expected")
	}

	s.Disconnected = make(chan bool)

	return &s, err
}

package session

import (
	"encoding/json"
	"reflect"
	"testing"

	"bitbucket.org/challengerdevs/gpsdriver/events"
)

func Test_newSession(t *testing.T) {
	type args struct {
		event events.Event
	}
	tests := []struct {
		name    string
		args    args
		want    *Session
		wantErr bool
	}{
		{
			"It should create a valid session",
			args{
				&eventMock{},
			},
			&Session{
				"123",
				[]byte("{\"string\":\"a string\"}"),
				make(chan bool),
			},
			false,
		},
		{
			"It should not create a valid session due to an empty json",
			args{
				&badEventMock{},
			},
			nil,
			true,
		},
		{
			"It should not create a valid session due to a wrong event type",
			args{
				&wrongTypeOfEventMock{},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newSession(tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("newSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil {
				if !reflect.DeepEqual(got.SessionID, tt.want.SessionID) || !reflect.DeepEqual(got.SessionAckPacket, tt.want.SessionAckPacket) {
					t.Errorf("eventSessionService.CreateSession() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

type badEventMock struct{}

func (badEventMock) Data() []byte {
	return nil
}

type wrongTypeOfEventMock struct {
}

func (wrongTypeOfEventMock) Data() []byte {
	type message struct {
		EventName string          `json:"event_name"`
		DeviceID  string          `json:"session_id"`
		AckPacket json.RawMessage `json:"ack_packet"`
	}

	event, _ := json.Marshal(message{
		"anotherMessage",
		"123",
		[]byte("{\"string\":\"a string\"}"),
	})

	return event
}

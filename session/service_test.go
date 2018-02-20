package session

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"bitbucket.org/challengerdevs/gpsdriver/events"
)

func TestNewEventSessionService(t *testing.T) {
	type args struct {
		eventEmitter           events.EventEmitter
		eventSubscriberFactory events.EventSubscriberFactory
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		{
			"It should create an event based session service",
			args{
				&eventEmitterMock{},
				&eventSubscriberFactoryMock{},
			},
			&eventSessionService{
				&eventEmitterMock{},
				&eventSubscriberFactoryMock{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEventSessionService(tt.args.eventEmitter, tt.args.eventSubscriberFactory); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEventSessionService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_eventSessionService_CreateSession(t *testing.T) {
	type args struct {
		sessionAddress string
		payload        []byte
	}
	tests := []struct {
		name    string
		s       eventSessionService
		args    args
		want    *Session
		wantErr bool
	}{
		{
			"It should create a session",
			eventSessionService{
				&eventEmitterMock{},
				&eventSubscriberFactoryMock{},
			},
			args{
				"1.1.1.1",
				[]byte("login"),
			},
			&Session{
				"123",
				[]byte("{\"string\":\"a string\"}"),
				make(chan bool),
			},
			false,
		},
		{
			"It should fail with an error on observing",
			eventSessionService{
				nil,
				&failureEventSubscriberFactoryMock{},
			},
			args{
				"1.1.1.1",
				[]byte("login"),
			},
			nil,
			true,
		},
		{
			"It should fail sending the event to the emitter",
			eventSessionService{
				&failEventEmitterMock{},
				&eventSubscriberFactoryMock{},
			},
			args{
				"1.1.1.1",
				[]byte("login"),
			},
			nil,
			true,
		},
		{
			"It should fail stoping the observer",
			eventSessionService{
				&eventEmitterMock{},
				&nonStopEventSubscriberFactoryMock{},
			},
			args{
				"1.1.1.1",
				[]byte("login"),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.CreateSession(tt.args.sessionAddress, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("eventSessionService.CreateSession() error = %v, wantErr %v", err, tt.wantErr)
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

func createChan() (chan events.Event, error) {
	response := make(chan events.Event)

	go func() {
		response <- &eventMock{}
	}()

	return response, nil
}

type eventEmitterMock struct{}

func (eventEmitterMock) Emit(string, []byte) error {
	return nil
}

type failEventEmitterMock struct{}

func (failEventEmitterMock) Emit(string, []byte) error {
	return fmt.Errorf("An error ocurred")
}

type eventSubscriberFactoryMock struct{}

func (eventSubscriberFactoryMock) CreateEventSubscriber() events.EventSubscriber {
	return &eventSubscriberMock{}
}

type eventSubscriberMock struct{}

func (eventSubscriberMock) Observe(...string) (<-chan events.Event, error) {
	return createChan()
}

func (eventSubscriberMock) Stop() error {
	return nil
}

type failureEventSubscriberFactoryMock struct {
}

func (failureEventSubscriberFactoryMock) CreateEventSubscriber() events.EventSubscriber {
	return &failEventSubscriberMock{}
}

type failEventSubscriberMock struct{}

func (failEventSubscriberMock) Observe(...string) (<-chan events.Event, error) {
	return nil, fmt.Errorf("An unexpected error")
}

func (failEventSubscriberMock) Stop() error {
	panic("not implemented")
}

type nonStopEventSubscriberFactoryMock struct {
}

func (nonStopEventSubscriberFactoryMock) CreateEventSubscriber() events.EventSubscriber {
	return &nonStopEventSubscriber{}
}

type nonStopEventSubscriber struct{}

func (nonStopEventSubscriber) Observe(...string) (<-chan events.Event, error) {
	return createChan()
}

func (nonStopEventSubscriber) Stop() error {
	return fmt.Errorf("An error ocurred")
}

type eventMock struct{}

func (eventMock) Data() []byte {
	type message struct {
		EventName string          `json:"event_name"`
		DeviceID  string          `json:"session_id"`
		AckPacket json.RawMessage `json:"ack_packet"`
	}

	event, _ := json.Marshal(message{
		"sessionAck",
		"123",
		[]byte("{\"string\":\"a string\"}"),
	})

	return event
}

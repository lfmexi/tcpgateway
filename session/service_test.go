package session

import (
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
			},
			&Session{
				"123",
				make(chan bool),
				&eventEmitterMock{},
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
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.CreateSession(tt.args.sessionAddress, "a-type")
			if (err != nil) != tt.wantErr {
				t.Errorf("eventSessionService.CreateSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil && tt.want != nil {
				if !reflect.DeepEqual(got.SessionID, tt.want.SessionID) {
					t.Errorf("eventSessionService.CreateSession() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_newSession(t *testing.T) {
	type args struct {
		event events.Event
	}
	tests := []struct {
		name    string
		s       eventSessionService
		args    args
		want    *Session
		wantErr bool
	}{
		{
			"It should create a valid session",
			eventSessionService{
				&eventEmitterMock{},
				&eventSubscriberFactoryMock{},
			},
			args{
				&eventMock{},
			},
			&Session{
				"123",
				make(chan bool),
				&eventEmitterMock{},
			},
			false,
		},
		{
			"It should not create a valid session due to an empty json",
			eventSessionService{
				&eventEmitterMock{},
				&eventSubscriberFactoryMock{},
			},
			args{
				&badEventMock{},
			},
			nil,
			true,
		},
		{
			"It should not create a valid session due to a wrong event type",
			eventSessionService{
				&eventEmitterMock{},
				&eventSubscriberFactoryMock{},
			},
			args{
				&wrongTypeOfEventMock{},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.newSession(tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("newSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil {
				if !reflect.DeepEqual(got.SessionID, tt.want.SessionID) {
					t.Errorf("eventSessionService.CreateSession() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

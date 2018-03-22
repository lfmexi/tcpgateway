package publisher

import (
	"fmt"
	"reflect"
	"testing"

	"bitbucket.org/challengerdevs/tcpgateway/events"
)

func TestNewEventPublisherService(t *testing.T) {
	type args struct {
		traceEventEmitter events.EventEmitter
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		{
			"It should create an event publisher service",
			args{
				&eventEmitterMock{},
			},
			&eventPublisherService{&eventEmitterMock{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEventPublisherService(tt.args.traceEventEmitter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEventPublisherService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_eventPublisherService_Publish(t *testing.T) {
	type args struct {
		sessionID string
		data      []byte
	}
	tests := []struct {
		name    string
		es      eventPublisherService
		args    args
		wantErr bool
	}{
		{
			"It should publish a message",
			eventPublisherService{&eventEmitterMock{}},
			args{
				"a-session-id-123",
				[]byte("a message to be published"),
			},
			false,
		},
		{
			"It should fail publishing a message",
			eventPublisherService{&onEmitErrorEmitterMock{}},
			args{
				"a-session-id-123",
				[]byte("This message will not be published"),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.es.Publish("dest", tt.args.sessionID, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("eventPublisherService.Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type eventEmitterMock struct{}

func (eventEmitterMock) Emit(string, string, []byte) error {
	return nil
}

type onEmitErrorEmitterMock struct{}

func (onEmitErrorEmitterMock) Emit(destination string, key string, message []byte) error {
	return fmt.Errorf("An error emitting %s", message)
}

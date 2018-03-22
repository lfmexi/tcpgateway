package packets

import (
	"fmt"
	"reflect"
	"testing"

	"bitbucket.org/challengerdevs/tcpgateway/events"
)

func TestNewPacketEventEmitter(t *testing.T) {
	type args struct {
		s events.EventSource
	}
	tests := []struct {
		name string
		args args
		want events.EventEmitter
	}{
		{"it should create a new packetEventEmitter", args{
			&eventSourceMock{},
		}, &packetEventEmitter{&eventSourceMock{}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPacketEventEmitter(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPacketEventEmitter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_packetEventEmitter_Emit(t *testing.T) {
	type args struct {
		key  string
		data []byte
	}
	tests := []struct {
		name    string
		p       *packetEventEmitter
		args    args
		wantErr bool
	}{
		{"It should send a message without an error",
			&packetEventEmitter{&eventSourceMock{}},
			args{
				"a.key.string",
				[]byte("a message"),
			},
			false,
		},
		{"It should send a message with an error",
			&packetEventEmitter{&publishErrorEventSource{}},
			args{
				"a.key.string",
				[]byte("a message"),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.Emit("", tt.args.key, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("packetEventEmitter.Emit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type eventSourceMock struct{}

func (eventSourceMock) Publish(string, key string, data []byte) error {
	return nil
}

func (eventSourceMock) Consume(key string) (<-chan events.Event, error) {
	panic("Not implemented")
}

func (eventSourceMock) Stop(key string) error {
	panic("Not implemented")
}

type publishErrorEventSource struct{}

func (publishErrorEventSource) Publish(string, key string, data []byte) error {
	return fmt.Errorf("An error publishing the message %s", data)
}

func (publishErrorEventSource) Consume(key string) (<-chan events.Event, error) {
	panic("not implemented")
}

func (publishErrorEventSource) Stop(key string) error {
	panic("not implemented")
}

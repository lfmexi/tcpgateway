package readers

import (
	"bufio"
	"log"
	"reflect"
	"strings"
	"testing"

	"bitbucket.org/challengerdevs/gpsdriver/publisher"
	"bitbucket.org/challengerdevs/gpsdriver/session"
)

func TestNewContinuousReaderServiceFactory(t *testing.T) {
	tests := []struct {
		name string
		want ReaderServiceFactory
	}{
		{
			"It should create a reader service factory",
			&continuousReaderServiceFactory{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewContinuousReaderServiceFactory(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewContinuousReaderServiceFactory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_continuousReaderServiceFactory_CreateReaderService(t *testing.T) {
	type args struct {
		reader           *bufio.Reader
		publisherService publisher.Service
	}
	tests := []struct {
		name string
		c    *continuousReaderServiceFactory
		args args
		want ReaderService
	}{
		{
			"it should create a new continuous reader service",
			&continuousReaderServiceFactory{},
			args{
				&bufio.Reader{},
				&publisherServiceMock{},
			},
			&continuousReaderService{
				&bufio.Reader{},
				&publisherServiceMock{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.CreateReaderService(tt.args.reader, tt.args.publisherService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("continuousReaderServiceFactory.CreateReaderService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_continuousReaderService_ReadTraces(t *testing.T) {
	type args struct {
		s *session.Session
	}
	tests := []struct {
		name    string
		c       *continuousReaderService
		args    args
		wantErr bool
	}{
		{
			"It should send an error while reading a buffer",
			&continuousReaderService{
				createReaderFromString(""),
				&publisherServiceMock{},
			},
			args{
				createSession(),
			},
			true,
		},
		{
			"It should publish an event based on a string",
			&continuousReaderService{
				createReaderFromString("a message\n"),
				&publisherServiceMock{},
			},
			args{
				createSession(),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.ReadTraces(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("continuousReaderService.ReadTraces() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type publisherServiceMock struct{}

func (publisherServiceMock) Publish(string, key string, message []byte) error {
	log.Printf("message retrieved: %s", message)
	return nil
}

func createReaderFromString(s string) *bufio.Reader {
	buffer := strings.NewReader(s)
	return bufio.NewReader(buffer)
}

func createSession() *session.Session {
	s := &session.Session{}

	s.Disconnected = make(chan bool)

	go func() {
	loop:
		for {
			select {
			case <-s.Disconnected:
				break loop
			}
		}
	}()

	return s
}

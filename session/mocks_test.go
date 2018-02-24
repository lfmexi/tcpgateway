package session

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/challengerdevs/gpsdriver/events"
)

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

func createChan() (chan events.Event, error) {
	response := make(chan events.Event)

	go func() {
		response <- &eventMock{}
	}()

	return response, nil
}

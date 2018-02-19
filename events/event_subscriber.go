package events

type Event interface {
	Data() []byte
}

type EventSubscriber interface {
	Observe(...string) (<-chan Event, error)
	Stop() error
}

type EventSubscriberFactory interface {
	CreateEventSubscriber() EventSubscriber
}

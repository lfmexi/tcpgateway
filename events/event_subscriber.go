package events

// Event interface that represents the events produced
type Event interface {
	Data() []byte // Data of the event
}

// EventSubscriber interface that represents the event listeners
type EventSubscriber interface {
	Observe(...string) (<-chan Event, error) // Observe the events
	Stop() error                             // Stop the observation
}

// EventSubscriberFactory factory interface for creating event subscribers
type EventSubscriberFactory interface {
	CreateEventSubscriber() EventSubscriber // CreateEventSubscriber creates a subscriber
}

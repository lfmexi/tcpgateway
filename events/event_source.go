package events

// EventSource interface that represents the event sources
type EventSource interface {
	Publish(destination string, key string, data []byte) error // Publish to an event source
	Consume(key string) (<-chan Event, error)                  // Consume an event source
	Stop(key string) error                                     // Stop the consumption of the event source
}

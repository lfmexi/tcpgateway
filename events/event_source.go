package events

type EventSource interface {
	Publish(key string, data []byte) error
	Consume(key string) (<-chan Event, error)
	Stop(key string) error
}

package events

// EventEmitter interface that allows sending events to an event source
type EventEmitter interface {
	Emit(string, []byte) error // Emit events to an envent source
}

package events

type EventEmitter interface {
	Emit(string, []byte) error
}

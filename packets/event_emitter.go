package packets

import (
	"bitbucket.org/challengerdevs/gpsdriver/events"
)

// NewPacketEventEmitter creates a new PacketEventEmitter
func NewPacketEventEmitter(s events.EventSource) events.EventEmitter {
	return &packetEventEmitter{
		s,
	}
}

type packetEventEmitter struct {
	source events.EventSource
}

func (p *packetEventEmitter) Emit(key string, data []byte) error {
	return p.source.Publish(key, data)
}

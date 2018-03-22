package packets

import (
	"bitbucket.org/challengerdevs/tcpgateway/events"
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

func (p *packetEventEmitter) Emit(destination string, key string, data []byte) error {
	return p.source.Publish(destination, key, data)
}

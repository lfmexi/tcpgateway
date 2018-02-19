package packets

import (
	"bitbucket.org/challengerdevs/gpsdriver/events"
)

type PacketEventEmitter struct {
	source events.EventSource
}

func (p *PacketEventEmitter) Emit(key string, data []byte) error {
	return p.source.Publish(key, data)
}

func NewPacketEventEmitter(s events.EventSource) events.EventEmitter {
	return &PacketEventEmitter{
		s,
	}
}

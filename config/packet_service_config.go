package config

import (
	"github.com/lfmexi/tcpgateway/events"
	"github.com/lfmexi/tcpgateway/packets"
	"github.com/lfmexi/tcpgateway/publisher"
	"github.com/lfmexi/tcpgateway/subscriber"
)

func packetEventEmitter() events.EventEmitter {
	return packets.NewPacketEventEmitter(packetEventSource())
}

func publisherService() publisher.Service {
	return publisher.NewEventPublisherService(packetEventEmitter())
}

func packetEventSubscriberFactory() events.EventSubscriberFactory {
	return subscriber.NewKeyBasedEventSubscriberFactory(packetEventSource())
}

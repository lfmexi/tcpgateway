package config

import (
	"bitbucket.org/challengerdevs/tcpgateway/events"
	"bitbucket.org/challengerdevs/tcpgateway/packets"
	"bitbucket.org/challengerdevs/tcpgateway/publisher"
	"bitbucket.org/challengerdevs/tcpgateway/subscriber"
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

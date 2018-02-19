package config

import (
	"bitbucket.org/challengerdevs/gpsdriver/events"
	"bitbucket.org/challengerdevs/gpsdriver/packets"
	"bitbucket.org/challengerdevs/gpsdriver/publisher"
	"bitbucket.org/challengerdevs/gpsdriver/subscriber"
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

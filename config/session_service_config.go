package config

import (
	"bitbucket.org/challengerdevs/gpsdriver/events"
	"bitbucket.org/challengerdevs/gpsdriver/session"
	"bitbucket.org/challengerdevs/gpsdriver/subscriber"
)

func sessionEventEmitter() events.EventEmitter {
	eventSource := sessionEventSource()
	return session.NewSessionEventEmitter(eventSource)
}

func sessionEventSubscriberFactory() events.EventSubscriberFactory {
	eventSource := sessionEventSource()
	return subscriber.NewKeyBasedEventSubscriberFactory(eventSource)
}

func sessionService() session.Service {
	return session.NewEventSessionService(sessionEventEmitter(), sessionEventSubscriberFactory())
}
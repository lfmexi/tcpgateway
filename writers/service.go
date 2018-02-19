package writers

import (
	"bufio"

	"bitbucket.org/challengerdevs/gpsdriver/session"

	"bitbucket.org/challengerdevs/gpsdriver/events"
)

// WriterService interface that represents the writer services
type WriterService interface {
	WriteSinglePacket([]byte) error
	WriteOnEventSubscriber(*session.Session, events.EventSubscriber) error
}

// WriterServiceFactory abstract factory of writers
type WriterServiceFactory interface {
	CreateWriterService(*bufio.Writer) WriterService
}

package writers

import (
	"bufio"

	"bitbucket.org/challengerdevs/gpsdriver/session"

	"bitbucket.org/challengerdevs/gpsdriver/events"
)

type WriterService interface {
	WriteSinglePacket([]byte) error
	WriteOnEventSubscriber(*session.Session, events.EventSubscriber) error
}

type WriterServiceFactory interface {
	CreateWriterService(*bufio.Writer) WriterService
}

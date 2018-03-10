package readers

import (
	"bufio"

	"bitbucket.org/challengerdevs/gpsdriver/session"

	"bitbucket.org/challengerdevs/gpsdriver/publisher"
)

// ReaderService interface that represents a reader service
type ReaderService interface {
	ReadTraces(*session.Session) error
}

// ReaderServiceFactory interface that describes an abstract factory of readers
type ReaderServiceFactory interface {
	CreateReaderService(*bufio.Reader, publisher.Service) ReaderService
}

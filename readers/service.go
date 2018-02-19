package readers

import (
	"bufio"

	"bitbucket.org/challengerdevs/gpsdriver/session"

	"bitbucket.org/challengerdevs/gpsdriver/publisher"
)

type ReaderService interface {
	ReadFirstLine() ([]byte, error)
	ReadTraces(*session.Session) error
}

type ReaderServiceFactory interface {
	CreateReaderService(*bufio.Reader, publisher.Service) ReaderService
}

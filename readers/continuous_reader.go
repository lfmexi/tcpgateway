package readers

import (
	"bufio"
	"fmt"
	"log"

	"bitbucket.org/challengerdevs/gpsdriver/session"

	"bitbucket.org/challengerdevs/gpsdriver/publisher"
)

// NewContinuousReaderServiceFactory creates a new reader factory of continuous readers
func NewContinuousReaderServiceFactory() ReaderServiceFactory {
	return &continuousReaderServiceFactory{}
}

type continuousReaderServiceFactory struct {
}

func (c *continuousReaderServiceFactory) CreateReaderService(reader *bufio.Reader, publisherService publisher.Service) ReaderService {
	return &continuousReaderService{
		reader,
		publisherService,
	}
}

type continuousReaderService struct {
	reader           *bufio.Reader
	publisherService publisher.Service
}

func (c *continuousReaderService) ReadFirstLine() ([]byte, error) {
	line, _, err := c.reader.ReadLine()

	if err == nil && len(line) == 0 {
		return nil, fmt.Errorf("First line of the buffer is empty")
	}

	return line, err
}

func (c *continuousReaderService) ReadTraces(s *session.Session) error {
	for {
		line, _, err := c.reader.ReadLine()

		if err != nil {
			log.Printf("EOF for session %s", s.SessionID)

			s.Disconnected <- true

			return err

		}

		if len(line) > 0 {
			c.publisherService.Publish(s.SessionID, line)
		}
	}
}

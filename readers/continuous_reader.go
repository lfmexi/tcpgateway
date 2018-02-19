package readers

import (
	"bufio"
	"fmt"
	"log"

	"bitbucket.org/challengerdevs/gpsdriver/session"

	"bitbucket.org/challengerdevs/gpsdriver/publisher"
)

type ContinuousReaderService struct {
	reader           *bufio.Reader
	publisherService publisher.Service
}

func (c *ContinuousReaderService) ReadFirstLine() ([]byte, error) {
	line, _, err := c.reader.ReadLine()

	if err == nil && len(line) == 0 {
		return nil, fmt.Errorf("First line of the buffer is empty")
	}

	return line, err
}

func (c *ContinuousReaderService) ReadTraces(s *session.Session) error {
	for {
		line, _, err := c.reader.ReadLine()

		if err != nil {
			if err.Error() == "EOF" {
				log.Printf("EOF for session %s", s.SessionID)
				err = nil
			}

			s.Disconnected <- true

			return err

		}

		if len(line) > 0 {
			c.publisherService.Publish(s.SessionID, line)
		}
	}
}

type ContinuousReaderServiceFactory struct {
}

func (c *ContinuousReaderServiceFactory) CreateReaderService(reader *bufio.Reader, publisherService publisher.Service) ReaderService {
	return &ContinuousReaderService{
		reader,
		publisherService,
	}
}

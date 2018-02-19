package writers

import (
	"bufio"
	"log"

	"bitbucket.org/challengerdevs/gpsdriver/events"
	"bitbucket.org/challengerdevs/gpsdriver/session"
)

// NewAsyncWriterServiceFactory creates a new factory of asynchronous writers
func NewAsyncWriterServiceFactory() WriterServiceFactory {
	return &asyncWriterServiceFactory{}
}

type asyncWriterServiceFactory struct{}

func (asyncWriterServiceFactory) CreateWriterService(writer *bufio.Writer) WriterService {
	return &asyncWriterService{
		writer,
	}
}

type asyncWriterService struct {
	writer *bufio.Writer
}

func (a *asyncWriterService) WriteSinglePacket(packet []byte) error {
	if _, err := a.writer.Write(packet); err != nil {
		return err
	}

	return a.writer.Flush()
}

func (a *asyncWriterService) WriteOnEventSubscriber(s *session.Session, es events.EventSubscriber) error {
	log.Printf("Waiting for messages to %s", s.SessionID)
	eventChannel, err := es.Observe(s.SessionID)

	if err != nil {
		return err
	}

	go func() {
	loop:
		for {
			select {
			case event := <-eventChannel:
				{
					if _, err := a.writer.Write(event.Data()); err != nil {
						log.Println(err)
					}

					if err := a.writer.Flush(); err != nil {
						log.Println(err)
					}
				}
			case <-s.Disconnected:
				{
					log.Printf("Closing connection for device %s, no more events listened", s.SessionID)
					es.Stop()
					break loop
				}
			}
		}
	}()

	return nil
}

package writers

import (
	"bufio"
	"log"

	"bitbucket.org/challengerdevs/gpsdriver/events"
	"bitbucket.org/challengerdevs/gpsdriver/session"
)

type AsyncWriterService struct {
	writer *bufio.Writer
}

func (a *AsyncWriterService) WriteSinglePacket(packet []byte) error {
	if _, err := a.writer.Write(packet); err != nil {
		return err
	}

	return a.writer.Flush()
}

func (a *AsyncWriterService) WriteOnEventSubscriber(s *session.Session, es events.EventSubscriber) error {
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

type AsyncWriterServiceFactory struct{}

func (AsyncWriterServiceFactory) CreateWriterService(writer *bufio.Writer) WriterService {
	return &AsyncWriterService{
		writer,
	}
}

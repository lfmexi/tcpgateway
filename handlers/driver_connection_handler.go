package handlers

import (
	"bufio"
	"net"

	"bitbucket.org/challengerdevs/gpsdriver/events"
	"bitbucket.org/challengerdevs/gpsdriver/publisher"
	"bitbucket.org/challengerdevs/gpsdriver/readers"
	"bitbucket.org/challengerdevs/gpsdriver/session"
	"bitbucket.org/challengerdevs/gpsdriver/writers"
)

// DriverConnectionHandler is the tcp handler that will be used for the driver
type DriverConnectionHandler struct {
	sessionService         session.Service
	publisherService       publisher.Service
	eventSubscriberFactory events.EventSubscriberFactory
	readerServiceFactory   readers.ReaderServiceFactory
	writerServiceFactory   writers.WriterServiceFactory
}

// NewConnectionHandler creates a new driver connection handler
func NewConnectionHandler(ss session.Service, publisher publisher.Service, evSubFactory events.EventSubscriberFactory, rsf readers.ReaderServiceFactory, wsf writers.WriterServiceFactory) *DriverConnectionHandler {
	return &DriverConnectionHandler{
		ss,
		publisher,
		evSubFactory,
		rsf,
		wsf,
	}
}

// ServeConnection handles a tcp connection
func (c *DriverConnectionHandler) ServeConnection(conn net.Conn) error {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	readerService := c.readerServiceFactory.CreateReaderService(reader, c.publisherService)
	writerService := c.writerServiceFactory.CreateWriterService(writer)

	data, err := readerService.ReadFirstLine()

	if err != nil {
		return err
	}

	session, err := c.sessionService.CreateSession(conn.RemoteAddr().String(), data)

	if err != nil {
		return err
	}

	defer func() {
		if session != nil {
			session.Disconnected <- true
		}
	}()

	if session.SessionAckPacket != nil {
		if err := writerService.WriteSinglePacket(session.SessionAckPacket); err != nil {
			return err
		}
	}

	subscriber := c.eventSubscriberFactory.CreateEventSubscriber()

	// This must launch a goroutine
	writerService.WriteOnEventSubscriber(session, subscriber)

	return readerService.ReadTraces(session)
}

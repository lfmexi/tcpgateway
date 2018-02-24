package handlers

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"bitbucket.org/challengerdevs/gpsdriver/events"
	"bitbucket.org/challengerdevs/gpsdriver/publisher"
	"bitbucket.org/challengerdevs/gpsdriver/readers"
	"bitbucket.org/challengerdevs/gpsdriver/server"
	"bitbucket.org/challengerdevs/gpsdriver/session"
	"bitbucket.org/challengerdevs/gpsdriver/writers"
)

// NewConnectionHandler creates a new driver connection handler
func NewConnectionHandler(portsMap map[string]string, ss session.Service, publisher publisher.Service, evSubFactory events.EventSubscriberFactory, rsf readers.ReaderServiceFactory, wsf writers.WriterServiceFactory) server.Handler {
	return &driverConnectionHandler{
		portsMap,
		ss,
		publisher,
		evSubFactory,
		rsf,
		wsf,
	}
}

type driverConnectionHandler struct {
	portsMap               map[string]string
	sessionService         session.Service
	publisherService       publisher.Service
	eventSubscriberFactory events.EventSubscriberFactory
	readerServiceFactory   readers.ReaderServiceFactory
	writerServiceFactory   writers.WriterServiceFactory
}

func (c *driverConnectionHandler) ServeConnection(conn net.Conn) error {
	defer conn.Close()

	deviceType, err := c.getDeviceTypeFromConn(conn)

	if err != nil {
		return err
	}

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	readerService := c.readerServiceFactory.CreateReaderService(reader, c.publisherService)
	writerService := c.writerServiceFactory.CreateWriterService(writer)

	data, err := readerService.ReadFirstLine()

	if err != nil {
		return err
	}

	session, err := c.sessionService.CreateSession(conn.RemoteAddr().String(), deviceType, data)

	if err != nil {
		return err
	}

	defer func() {
		if session != nil {
			session.Disconnected <- true
		}
	}()

	subscriber := c.eventSubscriberFactory.CreateEventSubscriber()

	// This must launch a goroutine
	writerService.WriteOnEventSubscriber(session, subscriber)

	return readerService.ReadTraces(session)
}

func (c *driverConnectionHandler) getDeviceTypeFromConn(conn net.Conn) (string, error) {
	address := conn.LocalAddr().String()
	addressParts := strings.Split(address, ":")

	port := addressParts[1]

	if deviceType, ok := c.portsMap[port]; ok {
		return deviceType, nil
	}

	return "", fmt.Errorf("Not a valid port %s", port)
}

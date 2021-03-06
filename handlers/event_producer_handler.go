package handlers

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/lfmexi/tcpgateway/events"
	"github.com/lfmexi/tcpgateway/publisher"
	"github.com/lfmexi/tcpgateway/server"
	"github.com/lfmexi/tcpgateway/session"
)

// NewConnectionHandler creates a new driver connection handler
func NewConnectionHandler(portsMap map[string]string, ss session.Service, publisher publisher.Service, evSubFactory events.EventSubscriberFactory) server.Handler {
	return &eventProducerHandler{
		portsMap,
		ss,
		publisher,
		evSubFactory,
	}
}

type eventProducerHandler struct {
	portsMap               map[string]string
	sessionService         session.Service
	publisherService       publisher.Service
	eventSubscriberFactory events.EventSubscriberFactory
}

func (c *eventProducerHandler) CreateSession(conn net.Conn, stopWaitGroup *sync.WaitGroup) (chan bool, string, error) {
	deviceType, err := c.getDeviceTypeFromConn(conn)

	if err != nil {
		return nil, "", err
	}

	sess, err := c.sessionService.CreateSession(conn.RemoteAddr().String(), deviceType)

	if err != nil {
		return nil, "", err
	}

	subscriber := c.eventSubscriberFactory.CreateEventSubscriber()
	log.Printf("Waiting for messages on session %s", sess.ID.Hex())
	eventChannel, err := subscriber.Observe(sess.ID.Hex())
	writer := bufio.NewWriter(conn)

	stopChannel := make(chan bool)

	go func() {

		for {
			select {
			case <-stopChannel:
				subscriber.Stop()
				return
			case event := <-eventChannel:
				if event == nil {
					continue
				}

				if _, err := writer.Write(event.Data()); err != nil {
					log.Println(err)
					continue
				}

				if err := writer.Flush(); err != nil {
					log.Println(err)
				}
			}
		}
	}()

	stopWaitGroup.Add(1)
	go func() {
		<-stopChannel
		c.sessionService.DisableSession(sess)
		stopWaitGroup.Done()
	}()

	return stopChannel, sess.ID.Hex(), nil
}

func (c *eventProducerHandler) ServeConnection(conn net.Conn, sessionID string) error {
	reader := bufio.NewReader(conn)

	line, _, err := reader.ReadLine()

	if err != nil {
		log.Printf("EOF for session %s", sessionID)
		return err
	}

	if len(line) > 0 {
		c.publisherService.Publish("packets", sessionID, line)
	}

	return nil
}

func (c *eventProducerHandler) getDeviceTypeFromConn(conn net.Conn) (string, error) {
	address := conn.LocalAddr().String()
	addressParts := strings.Split(address, ":")

	port := addressParts[1]

	if deviceType, ok := c.portsMap[port]; ok {
		return deviceType, nil
	}

	return "", fmt.Errorf("Not a valid port %s", port)
}

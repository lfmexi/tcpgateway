package server

import (
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Handler interface that represents the behavior of a handler
type Handler interface {
	ServeConnection(conn net.Conn, sessionID string) error
	CreateSession(conn net.Conn, stopWaitGroup *sync.WaitGroup) (chan bool, string, error)
}

// Server interface that represents the behavior of a server
type Server interface {
	Listen()
}

// NewTCPServer creates a new Server that listen on an TCP port
func NewTCPServer(address string, ports map[string]string, onConnectionHandler Handler) Server {
	return &tcpServer{
		address,
		ports,
		onConnectionHandler,
	}
}

type tcpServer struct {
	addr    string
	ports   map[string]string
	handler Handler
}

func (s *tcpServer) Listen() {
	exit := make(chan os.Signal, 1)

	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	exitChannels := make([]chan bool, 0)
	waitGroups := make([]*sync.WaitGroup, 0)

	for port := range s.ports {
		stopListen := make(chan bool)

		exitChannels = append(exitChannels, stopListen)
		var wg sync.WaitGroup
		waitGroups = append(waitGroups, &wg)

		go s.listen(port, stopListen, &wg)
	}

	sig := <-exit
	log.Printf("Exiting with signal %s", sig)

	log.Printf("Closing %d servers", len(exitChannels))
	for i, channel := range exitChannels {
		close(channel)
		waitGroups[i].Wait()
		log.Printf("Server closed")
	}
}

func (s *tcpServer) listen(port string, stopListen <-chan bool, globalWG *sync.WaitGroup) {
	listen, err := net.Listen("tcp", s.addr+":"+port)

	defer listen.Close()

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("Server listenning at %s", s.addr+":"+port)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			return
		}

		go func() {
			<-stopListen
			conn.Close()
		}()

		if s.handler != nil {
			globalWG.Add(1)
			go s.handleConnection(conn, stopListen, globalWG)
		}
	}
}

func (s *tcpServer) handleConnection(conn net.Conn, stopListen <-chan bool, globalWG *sync.WaitGroup) {
	defer globalWG.Done()

	var stopWaitGroup sync.WaitGroup
	sessionStopChannel, sessionID, err := s.handler.CreateSession(conn, &stopWaitGroup)

	if err != nil {
		log.Fatal(err)
		return
	}

loop:
	for {
		select {
		case <-stopListen:
			log.Printf("Stoping connection listening for %s", conn.RemoteAddr())
			break loop
		default:
			err := s.handler.ServeConnection(conn, sessionID)
			if err != nil {
				if err.Error() != "EOF" {
					log.Println(err.Error())
				}
				break loop
			}
		}
	}

	close(sessionStopChannel)
	stopWaitGroup.Wait()
	conn.Close()
}

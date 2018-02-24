package server

import (
	"log"
	"net"
)

// Handler interface that represents the behavior of a handler
type Handler interface {
	ServeConnection(conn net.Conn) error
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

func (s *tcpServer) handleConnection(conn net.Conn) {
	log.Printf("Handling connection from %s", conn.RemoteAddr())
	if err := s.handler.ServeConnection(conn); err != nil {
		log.Println(err)
	}
}

func (s *tcpServer) Listen() {
	exit := make(chan bool)

	for port, _ := range s.ports {
		go s.listenPort(port)
	}

	<-exit
}

func (s *tcpServer) listenPort(port string) {
	listen, err := net.Listen("tcp", s.addr+":"+port)

	defer listen.Close()

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server listenning at %s", s.addr+":"+port)

	for {
		conn, err := listen.Accept()

		if err == nil {
			log.Printf("server accepting connection from %s", conn.RemoteAddr())
			go s.handleConnection(conn)
		} else {
			log.Println(err)
		}

	}
}

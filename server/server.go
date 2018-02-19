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
func NewTCPServer(address string, onConnectionHandler Handler) Server {
	return &tcpServer{
		address,
		onConnectionHandler,
	}
}

type tcpServer struct {
	addr    string
	handler Handler
}

func (s *tcpServer) handleConnection(conn net.Conn) {
	log.Printf("Handling connection from %s", conn.RemoteAddr())
	if err := s.handler.ServeConnection(conn); err != nil {
		log.Println(err)
	}
}

func (s *tcpServer) Listen() {
	listen, err := net.Listen("tcp", s.addr)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server listenning at %s", s.addr)

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

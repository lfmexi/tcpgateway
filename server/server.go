package server

import (
	"log"
	"net"
)

type Handler interface {
	ServeConnection(conn net.Conn) error
}

type Server interface {
	Listen()
}

type TCPServer struct {
	addr    string
	handler Handler
}

func NewTCPServer(address string, onConnectionHandler Handler) *TCPServer {
	return &TCPServer{
		address,
		onConnectionHandler,
	}
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	log.Printf("Handling connection from %s", conn.RemoteAddr())
	if err := s.handler.ServeConnection(conn); err != nil {
		log.Println(err)
	}
}

func (s *TCPServer) Listen() {
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

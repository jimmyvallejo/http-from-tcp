package server

import (
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	conn     net.Conn
	running  atomic.Bool
}

func Serve(port int) (*Server, error) {

	portString := strconv.Itoa(port)

	l, err := net.Listen("tcp", ":"+portString)
	if err != nil {
		log.Fatal(err)
	}

	serv := Server{
		listener: l,
	}

	go serv.listen()

	return &serv, nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.running.Load() {
				log.Fatal(err)
			}
			return
		}
		s.conn = conn
		go s.handle(conn)
	}

}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	response := []byte("HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		"Hello World!")

	conn.Write(response)
}

func (s *Server) Close() {
	s.listener.Close()
	s.conn.Close()
}

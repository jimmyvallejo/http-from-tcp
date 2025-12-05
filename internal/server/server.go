package server

import (
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/jimmyvallejo/httpfromtcp/internal/response"
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

	response.WriteStatusLine(conn, 200)
	headers := response.GetDefaultHeaders(0)
	response.WriteHeaders(conn, headers)
}

func (s *Server) Close() {
	s.listener.Close()
	s.conn.Close()
}

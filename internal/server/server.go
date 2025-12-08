package server

import (
	"bytes"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/jimmyvallejo/httpfromtcp/internal/request"
	"github.com/jimmyvallejo/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	running  atomic.Bool
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {

	portString := strconv.Itoa(port)

	l, err := net.Listen("tcp", ":"+portString)
	if err != nil {
		log.Fatal(err)
	}

	serv := Server{
		handler:  handler,
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
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}

}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.StatusCodeBadRequest,
			Message:    err.Error(),
		}
		hErr.Write(conn)
		return
	}
	buf := bytes.NewBuffer([]byte{})
	hErr := s.handler(buf, req)
	if hErr != nil {
		hErr.Write(conn)
		return
	}
	b := buf.Bytes()
	response.WriteStatusLine(conn, response.StatusCodeOk)
	headers := response.GetDefaultHeaders(len(b))
	response.WriteHeaders(conn, headers)
	conn.Write(b)
}

func (s *Server) Close() error {
	s.running.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

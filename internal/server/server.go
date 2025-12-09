package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/jimmyvallejo/httpfromtcp/internal/request"
	"github.com/jimmyvallejo/httpfromtcp/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

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

	w := response.Writer{
		WriterState: 0,
		Dest:        conn,
	}
	req, err := request.RequestFromReader(conn)
	if err != nil {
		w.WriteStatusLine(response.StatusCodeBadRequest)
		body := []byte(fmt.Sprintf("Error parsing request: %v", err))
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.Dest.Write(body)
		return
	}

	s.handler(&w, req)
}

func (s *Server) Close() error {
	s.running.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

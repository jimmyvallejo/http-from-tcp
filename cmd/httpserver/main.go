package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jimmyvallejo/httpfromtcp/internal/request"
	"github.com/jimmyvallejo/httpfromtcp/internal/response"
	"github.com/jimmyvallejo/httpfromtcp/internal/server"
)

const port = 42069

func main() {

	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w response.Writer, req *request.Request) *server.HandlerError {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		return &server.HandlerError{
			StatusCode: response.StatusCodeBadRequest,
			Message:    "Your request honestly kinda sucked.\n",
		}
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		return &server.HandlerError{
			StatusCode: response.StatusCodeInternalServerError,
			Message:    "Okay, you know what? This one is on me.\n",
		}
	}
	w.WriteStatusLine(response.StatusCodeOk)
	messageBytes := []byte("Your request was an absolute banger.\n")
	headers := response.GetDefaultHeaders(len(messageBytes))
	headers.Override("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.Dest.Write(messageBytes)
	return nil
}

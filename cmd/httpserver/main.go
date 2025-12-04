package httpserver

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Server struct{}

const port = 42069

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", ":"+string(port))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	server, err := Server.Serve(port)
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

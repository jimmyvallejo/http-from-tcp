package main

import (
	"fmt"
	"log"
	"net"

	"github.com/jimmyvallejo/httpfromtcp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("failed to open tcp link: %v", err)
	}

	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("failed to accept connection: %v", err)
		}
		fmt.Println("connection accepted")

		returnedReq, err := request.RequestFromReader(connection)
		if err != nil {
			fmt.Printf("failed to read from connection, err: %e", err)
		}

		fmt.Println("Request line:")
		fmt.Println("- Method:", returnedReq.RequestLine.Method)
		fmt.Println("- Target:", returnedReq.RequestLine.RequestTarget)
		fmt.Println("- Version:", returnedReq.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for key, item := range returnedReq.Headers {
			fmt.Println("- " + key + ": " + item)
		}
	}
}

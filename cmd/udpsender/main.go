package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	resolved, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		fmt.Printf("failed to resolve UDP: %v", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, resolved)
	if err != nil {
		fmt.Printf("failed to dial UDP: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(">")
		readString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("failed to read from buffer: %v", err)
		}
		_, err = conn.Write([]byte(readString))
		if err != nil {
			fmt.Printf("failed to write to conn: %v", err)
		}
	}

}

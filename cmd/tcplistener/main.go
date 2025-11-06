package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {

	strChan := make(chan string)

	go func() {
		defer f.Close()
		defer close(strChan)

		var currentLine string

		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			str := string(buffer[:n])
			parts := strings.Split(str, "\n")

			for i := 0; i < len(parts)-1; i++ {
				strChan <- fmt.Sprintf("%s%s", currentLine, parts[i])
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]
		}

		strChan <- currentLine
	}()

	return strChan

}

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

		returnedChan := getLinesChannel(connection)

		for line := range returnedChan {
			fmt.Println(line)
		}
		fmt.Println("connection closed")
	}

}

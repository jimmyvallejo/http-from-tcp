package main

import (
	"fmt"
	"io"
	"log"
	"os"
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
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("failed to open")
	}

	returnedChan := getLinesChannel(file)

	for line := range returnedChan {
		fmt.Println("read:", line)
	}
}

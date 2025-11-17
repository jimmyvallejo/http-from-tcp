package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const supportedVersion = "1.1"

func RequestFromReader(reader io.Reader) (*Request, error) {

	read, err := io.ReadAll(reader)
	if err != nil {
		fmt.Printf("failed to read, err: %e", err)
	}

	s := string(read)

	return parseRequestLine(s)

}

func parseRequestLine(readerString string) (*Request, error) {

	lines := strings.Split(readerString, "\r\n")

	line := lines[0]

	parts := strings.Split(line, " ")

	if len(parts) < 3 {
		return nil, errors.New("line length is incorrect")
	}

	method := parts[0]
	if !isAllUpperAlpha(method) {
		return nil, errors.New("method is incorrect")
	}

	requestTarget := parts[1]

	httpVersion := parts[2]

	versionParts := strings.Split(httpVersion, "/")
	version := versionParts[1]

	if version != supportedVersion {
		return nil, errors.New("version is incorrect")
	}

	requestLine := RequestLine{
		HttpVersion:   version,
		RequestTarget: requestTarget,
		Method:        method,
	}

	request := Request{
		RequestLine: requestLine,
	}

	return &request, nil

}

func isAllUpperAlpha(s string) bool {
	if len(s) == 0 {
		return false
	}

	for _, r := range s {
		if !unicode.IsUpper(r) || !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

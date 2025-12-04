package request

import (
	"errors"
	"io"
	"strings"

	"github.com/jimmyvallejo/httpfromtcp/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	Enum        int
	Headers     headers.Headers
	Body        []byte
	BodyLength  int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const (
	requestStateParsingLine    = 0
	requestStateDone           = 1
	requestStateParsingHeaders = 2
	requestStateParsingBody    = 3
)

const supportedVersion = "1.1"

const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {

	buf := make([]byte, bufferSize)

	readToIndex := 0

	req := Request{Enum: requestStateParsingLine, Headers: headers.NewHeaders(), Body: []byte{}}

	for {
		if readToIndex == len(buf) {
			nb := make([]byte, len(buf)*2)
			copy(nb, buf[:readToIndex])
			buf = nb
		}

		n, err := reader.Read(buf[readToIndex:])
		if n > 0 {
			readToIndex += n

			parsed, perr := req.parse(buf[:readToIndex])
			if perr != nil {
				return nil, perr
			}
			if parsed > 0 {
				copy(buf, buf[parsed:readToIndex])
				readToIndex -= parsed
			}
			if req.Enum == requestStateDone {
				return &req, nil
			}
		}
		if errors.Is(err, io.EOF) {
			parsed, perr := req.parse(buf[:readToIndex])
			if perr != nil {
				return nil, perr
			}
			if parsed > 0 {
				copy(buf, buf[parsed:readToIndex])
				readToIndex -= parsed
			}
			if req.Enum == requestStateDone {
				return &req, nil
			}
			return nil, io.EOF
		}
		if err != nil {
			return nil, err
		}
	}
}

func (r *Request) parse(data []byte) (int, error) {

	totalBytesParsed := 0
	for r.Enum != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			break
		}
		totalBytesParsed += n
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.Enum {
	case requestStateParsingLine:
		reqLine, bytes, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		if bytes == 0 {
			return 0, nil
		}
		r.Enum = requestStateParsingHeaders
		r.RequestLine = reqLine
		return bytes, nil

	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.Enum = requestStateParsingBody
		}
		return n, err

	case requestStateParsingBody:
		return r.parseBody(data)

	case requestStateDone:
		return 0, errors.New("error: trying to read data in a done state")
	}
	return 0, errors.New("error: unknown state")
}

func parseRequestLine(readerString string) (RequestLine, int, error) {

	lines := strings.Split(readerString, "\r\n")
	if len(lines) < 2 {
		return RequestLine{}, 0, nil
	}

	line := lines[0]

	parts := strings.Split(line, " ")

	if len(parts) < 3 {
		return RequestLine{}, 0, errors.New("line length is incorrect")
	}

	method := parts[0]
	if !isAllUpperAlpha(method) {
		return RequestLine{}, 0, errors.New("method is incorrect")
	}

	requestTarget := parts[1]

	httpVersion := parts[2]

	versionParts := strings.Split(httpVersion, "/")
	version := versionParts[1]

	if version != supportedVersion {
		return RequestLine{}, 0, errors.New("version is incorrect")
	}

	requestLine := RequestLine{
		HttpVersion:   version,
		RequestTarget: requestTarget,
		Method:        method,
	}

	bytesConsumed := len(line) + 2

	return requestLine, bytesConsumed, nil

}

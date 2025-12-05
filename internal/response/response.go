package response

import (
	"errors"
	"io"
	"strconv"

	"github.com/jimmyvallejo/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusCodeOk                  = 200
	StatusCodeBadRequest          = 400
	StatusCodeInternalServerError = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case StatusCodeOk:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return errors.New("failed to write to io.Writer")
		}
	case StatusCodeBadRequest:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		if err != nil {
			return errors.New("failed to write to io.Writer")
		}
	case StatusCodeInternalServerError:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		if err != nil {
			return errors.New("failed to write to io.Writer")
		}
	default:
		_, err := w.Write([]byte("HTTP/1.1\r\n"))
		if err != nil {
			return errors.New("failed to write to io.Writer")
		}
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	lenToString := strconv.Itoa(contentLen)

	headers := headers.NewHeaders()
	headers["Content-Length"] = lenToString
	headers["Content-Type"] = "text/plain"
	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, item := range headers {
		formatted := []byte(key + ": " + item + "\r\n")
		_, err := w.Write(formatted)
		if err != nil {
			return errors.New("failed to write to io.Writer")
		}
	}
	w.Write([]byte("\r\n"))
	return nil
}

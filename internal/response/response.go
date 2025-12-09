package response

import (
	"errors"
	"io"
	"strconv"

	"github.com/jimmyvallejo/httpfromtcp/internal/headers"
)

type StatusCode int

type Writer struct {
	WriterState Writerstate
	Dest        io.Writer
}

type Writerstate int

const (
	WriterStateStart Writerstate = iota
	WriterStateLine
	WriterStateHeaders
)

const (
	StatusCodeOk                  = 200
	StatusCodeBadRequest          = 400
	StatusCodeInternalServerError = 500
)

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	switch statusCode {
	case StatusCodeOk:
		_, err := w.Dest.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return errors.New("failed to write to io.Writer")
		}
	case StatusCodeBadRequest:
		_, err := w.Dest.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		if err != nil {
			return errors.New("failed to write to io.Writer")
		}
	case StatusCodeInternalServerError:
		_, err := w.Dest.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		if err != nil {
			return errors.New("failed to write to io.Writer")
		}
	default:
		_, err := w.Dest.Write([]byte("HTTP/1.1\r\n"))
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

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	for key, item := range headers {
		formatted := []byte(key + ": " + item + "\r\n")
		_, err := w.Dest.Write(formatted)
		if err != nil {
			return errors.New("failed to write to io.Writer")
		}
	}
	w.Dest.Write([]byte("\r\n"))
	return nil
}

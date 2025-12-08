package server

import (
	"fmt"
	"io"

	"github.com/jimmyvallejo/httpfromtcp/internal/request"
	"github.com/jimmyvallejo/httpfromtcp/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (e HandlerError) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
}

func (e HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, e.StatusCode)
	messageBytes := []byte(e.Message)
	headers := response.GetDefaultHeaders(len(messageBytes))
	response.WriteHeaders(w, headers)
	w.Write(messageBytes)
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

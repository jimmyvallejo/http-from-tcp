package server

import (
	"fmt"

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

func (e HandlerError) Write(w response.Writer) {
	w.WriteStatusLine(e.StatusCode)
	messageBytes := []byte(e.Message)
	headers := response.GetDefaultHeaders(len(messageBytes))
	headers.Override("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.Dest.Write(messageBytes)
}

type Handler func(w *response.Writer, req *request.Request)

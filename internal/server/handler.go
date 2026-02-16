package server

import (
	"fmt"
	"io"

	"github.com/pavychsasha/httpfromtcp/internal/request"
	"github.com/pavychsasha/httpfromtcp/internal/response"
)

type HandlerError struct {
	Message string
	Status  response.StatusCode
}
type Handler func(w io.Writer, req *request.Request) *HandlerError

func HandleError(e HandlerError, w io.Writer) error {
	_, err := fmt.Fprintf(w, "%s %d%s%s", response.HTTPVersion, e.Status, e.Message, response.CRLF)

	return err
}

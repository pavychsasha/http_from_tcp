package server

import (
	"github.com/pavychsasha/httpfromtcp/internal/request"
	"github.com/pavychsasha/httpfromtcp/internal/response"
)

type HandlerError struct {
	Message string
	Status  response.StatusCode
}
type Handler func(w *response.Writer, req *request.Request)

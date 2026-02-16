package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/pavychsasha/httpfromtcp/internal/request"
	"github.com/pavychsasha/httpfromtcp/internal/response"
)

type Server struct {
	Port     int
	Listener net.Listener
	Closed   atomic.Bool
}

func Serve(port int, handler Handler) (*Server, error) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	srv := &Server{Port: port, Listener: listener, Closed: atomic.Bool{}}
	go srv.listen()
	return srv, nil
}

func (srv *Server) Close() error {
	if srv.Closed.Load() == true {
		return nil
	}
	err := srv.Listener.Close()
	if err != nil {
		return err
	}
	srv.Closed.Store(true)
	return nil
}

func (srv *Server) listen() error {
	for !srv.Closed.Load() {
		conn, err := srv.Listener.Accept()
		defer conn.Close()

		if err != nil {
			log.Fatal(err)
		}
		go func(c net.Conn) {
			srv.handle(c)
		}(conn)
	}
	return nil
}

func (srv *Server) handle(c net.Conn) {
	defer c.Close()
	request, err := request.RequestFromReader(c)
	if err != nil {
		HandleError(HandlerError{Message: err.Error(), Status: response.StatusInternalServerError}, c)
		return
	}

	buff := bytes.NewBuffer([]byte{})
	err = Handler(buff, request)

	response.WriteStatusLine(c, response.StatusResponseOK)
	headers := response.GetDefaultHeaders(0)
	response.WriteHeaders(c, headers)

}

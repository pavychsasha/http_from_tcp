package server

import (
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
	Handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	srv := &Server{Port: port, Listener: listener, Closed: atomic.Bool{}, Handler: handler}
	go srv.listen()
	return srv, nil
}

func (srv *Server) Close() error {
	if srv.Closed.Load() {
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

		if err != nil {
			log.Fatal(err)
		}
		go srv.handle(conn)
	}
	return nil
}

func (srv *Server) handle(c net.Conn) {
	defer c.Close()
	request, err := request.RequestFromReader(c)
	if err != nil {
		log.Printf("Could not parse a request: %v", err)
		return
	}

	w := response.Writer{
		WriterState: response.WriterInitialized,
		Writer:      c,
	}

	srv.Handler(
		&w,
		request,
	)

}

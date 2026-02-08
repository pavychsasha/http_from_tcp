package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/pavychsasha/httpfromtcp/internal/request"
)

const tcpPortListener = ":42069"

func main() {

	listener, err := net.Listen("tcp", tcpPortListener)
	if err != nil {
		log.Fatalf("Could not open a tcp listener: %s\n", err.Error())
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", tcpPortListener)
	fmt.Println("=====================================")
	for {
		conn, err := listener.Accept()

		fmt.Println("Accepted connection from", conn.RemoteAddr())
		if err != nil {
			log.Fatalf("Tcp listener was not accepted: %s\n", conn)
			break
		}

		request, err := request.RequestFromReader(conn)
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", request.RequestLine.Method, request.RequestLine.RequestTarget, request.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for k, value := range request.Headers {
			fmt.Printf("- %s: %s\n", k, value)
		}
		fmt.Printf("Body:\n%s\n", request.Body)

		fmt.Println("Connection to", conn.RemoteAddr(), "closed")
	}
}

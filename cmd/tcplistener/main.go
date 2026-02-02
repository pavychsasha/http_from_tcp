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
		fmt.Println("Connection to", conn.RemoteAddr(), "closed")
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	res := make(chan string)
	go func() {
		defer f.Close()
		defer close(res)

		currentLine := ""
		for {
			b := make([]byte, 8)
			n, err := f.Read(b)
			if err != nil {
				if currentLine != "" {
					res <- currentLine
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("Error: %s\n", err.Error())
				break
			}
			output := string(b[:n])
			parts := strings.Split(output, "\n")

			for i := 0; i < len(parts)-1; i++ {
				res <- currentLine + parts[i]
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]
		}
	}()
	return res
}

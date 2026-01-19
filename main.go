package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const tcpPortListener = ":42069"

func main() {

	listener, err := net.Listen("tcp", tcpPortListener)
	defer listener.Close()

	if err != nil {
		log.Fatalf("Could not open a tcp listener: %s", err)
	}

	// fmt.Printf("Reading data from %s\n", inputPath)
	fmt.Println("=====================================")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Tcp listener was not accepted: %s\n", conn)
			break
		}
		linesCh := getLinesChannel(conn)
		for l := range linesCh {
			fmt.Println(l)
		}
	}
	fmt.Println("Connection has been closed")

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

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const udpPortResolver = ":42069"

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", udpPortResolver)
	if err != nil {
		log.Fatalf("Could not resolve an UDP address: %s\n", err.Error())
	}
	conn, err := net.DialUDP("udp", &net.UDPAddr{}, udpAddr)

	if err != nil {
		log.Fatalf("Could not establish a dial UDP connection: %s\n", err.Error())
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			continue
		}
		_, err = conn.Write([]byte(input))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending message: %v\n", err)
			continue
		}
	}
}

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	port := flag.Int("port", 9000, "specify the listened port")
	address := flag.String("address", "0.0.0.0", "specify the listened address")
	flag.Parse()

	tcpAddress := fmt.Sprintf("%s:%d", *address, *port)
	tcpListener, err := net.Listen("tcp", tcpAddress)
	if err != nil {
		panic(err)
	}

	defer tcpListener.Close()

	log.Printf("[Info] listened tcp in %s\n", tcpAddress)

	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			log.Printf("[Error] error when accept the connections: %v", err)
			continue
		}

		go func(c net.Conn) {
			_, err := io.Copy(c, c)
			if err != nil {
				log.Printf("[Error] error when copy the data: %v", err)
			}

			err = c.Close()
			if err != nil {
				log.Printf("[Error] error when closing the tcp connections: %v", err)
			}
		}(conn)
	}
}

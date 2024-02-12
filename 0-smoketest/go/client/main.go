package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {
	port := flag.Int("port", 9000, "specify the listened port")
	address := flag.String("address", "0.0.0.0", "specify the listened address")
	concurrent := flag.Int("concurrent", 5, "specify the concurrent")
	flag.Parse()

	tcpServer := fmt.Sprintf("%s:%d", *address, *port)

	errChan := make(chan error)
	defer close(errChan)

	for i := 0; i < *concurrent; i++ {
		go func(request int) {
			conn, err := net.Dial("tcp", tcpServer)
			if err != nil {
				log.Printf("[Error] error when dial tcp server: %v", err)
				errChan <- err
			}
			defer conn.Close()

			_, err = fmt.Fprintf(conn, "connection number %d\n", request)
			if err != nil {
				log.Printf("[Error] error when send the request from server: %v", err)
				errChan <- err
			}

			response, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				log.Printf("[Error] error when read the response from server: %v", err)
				errChan <- err
			}

			log.Println("Server's response:", response)
			errChan <- nil
		}(i)
	}

	var finishCount int
	var joinErr error

	for finishCount < *concurrent {
		err := <-errChan
		if err != nil {
			joinErr = errors.Join(joinErr, err)
		}

		finishCount++
	}

	if joinErr != nil {
		log.Printf("[Error] finish with error: %v", joinErr)
	}

	log.Println("[Info] client ended")
}

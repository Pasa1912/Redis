package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	fmt.Println("Listening on port :6379")

	// Create a new server
	l, err := net.Listen("tcp", ":6379")

	if err != nil {
		log.Fatal("unable to start server: ", err)
		return
	}

	conn, err := l.Accept()
	if err != nil {
		log.Fatal(err)
		return
	}

	defer conn.Close()

	for {
		buf := make([]byte, 1024)

		// read message from client
		_, err = conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal("error reading from client: ", err.Error())
			os.Exit(1)
		}

		// ignore request and send back a PONG
		conn.Write([]byte("+OK\r\n"))
	}
}

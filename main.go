package main

import (
	"fmt"
	"net"
	"strings"
)

const MAX_CONNECTIONS = 5

func handleConnection(l *net.Listener, aof *Aof, connChan chan<- bool) {
	conn, err := (*l).Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		connChan <- true
		conn.Close()
	}()

	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		if command == "SET" || command == "HSET" {
			aof.Write(value)
		}

		result := handler(args)

		writer.Write(result)
	}
}

func main() {
	fmt.Println("Listening on port :6379")

	aof, err := NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer aof.Close()

	// Create a new server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	aof.Read(func(value Value) {
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		fmt.Println("Read " + command + " command from AOF")

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			return
		}

		handler(args)
	})

	connChan := make(chan bool)

	// Listen for connections
	for range MAX_CONNECTIONS {
		go handleConnection(&l, aof, connChan)
	}

	for range MAX_CONNECTIONS {
		<-connChan
	}
}

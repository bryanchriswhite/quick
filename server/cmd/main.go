package main

import (
	"bufio"
	"encoding/binary"
	"log"
	"net"
	"sync"
)

var (
	counter uint64
	mutex   sync.Mutex
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		// Read the request
		var value uint64
		err := binary.Read(reader, binary.LittleEndian, &value)
		if err != nil {
			log.Println("Error reading from client:", err)
			return
		}

		// Increment the global counter
		mutex.Lock()
		counter += value
		newCounter := counter
		mutex.Unlock()

		// Send the response
		err = binary.Write(conn, binary.LittleEndian, newCounter)
		if err != nil {
			log.Println("Error writing to client:", err)
			return
		}
	}
}

func main() {
	// Start the server
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	log.Println("Server started on :8080")

	// Accept multiple client connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

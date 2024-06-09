package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Send multiple increment requests
	rand.Seed(time.Now().UnixNano())
	for {
		// Generate a random unsigned integer
		value := uint64(rand.Intn(100))

		// Send the request
		err = binary.Write(writer, binary.LittleEndian, value)
		if err != nil {
			log.Println("Error writing to server:", err)
			return
		}
		writer.Flush()

		// Read the response
		var newCounter uint64
		err = binary.Read(reader, binary.LittleEndian, &newCounter)
		if err != nil {
			log.Println("Error reading from server:", err)
			return
		}

		// Print the new counter value
		fmt.Println("New counter value:", newCounter)

		// Sleep for a random duration to simulate multiple requests
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	}
}

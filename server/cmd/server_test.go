package main

import (
	"bufio"
	"encoding/binary"
	"net"
	"sync"
	"testing"
	"time"
)

func startServer(t *testing.T, wg *sync.WaitGroup) {
	defer wg.Done()
	main() // Start the server
}

func sendIncrementRequest(t *testing.T, conn net.Conn, value uint64) uint64 {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	// Send the increment request
	err := binary.Write(writer, binary.LittleEndian, value)
	if err != nil {
		t.Fatalf("Error writing to server: %v", err)
	}
	writer.Flush()

	// Read the response
	var newCounter uint64
	err = binary.Read(reader, binary.LittleEndian, &newCounter)
	if err != nil {
		t.Fatalf("Error reading from server: %v", err)
	}

	return newCounter
}

func TestSingleClient(t *testing.T) {
	// Start the server in a separate goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go startServer(t, &wg)

	// Give the server a moment to start
	time.Sleep(time.Second)

	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Send increment requests and check the responses
	counter := sendIncrementRequest(t, conn, 1)
	if counter != 1 {
		t.Fatalf("Expected counter to be 1, got %d", counter)
	}

	counter = sendIncrementRequest(t, conn, 2)
	if counter != 3 {
		t.Fatalf("Expected counter to be 3, got %d", counter)
	}
}

func TestMultipleClients(t *testing.T) {
	// Start the server in a separate goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go startServer(t, &wg)

	// Give the server a moment to start
	time.Sleep(time.Second)

	// Number of clients
	numClients := 10

	// Create clients and send increment requests
	var clients []net.Conn
	for i := 0; i < numClients; i++ {
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			t.Fatalf("Failed to connect to server: %v", err)
		}
		defer conn.Close()
		clients = append(clients, conn)
	}

	// Send increment requests from each client
	expectedCounter := uint64(0)
	for i := 0; i < numClients; i++ {
		counter := sendIncrementRequest(t, clients[i], 1)
		expectedCounter++
		if counter != expectedCounter {
			t.Fatalf("Expected counter to be %d, got %d", expectedCounter, counter)
		}
	}

	// Send more increment requests
	for i := 0; i < numClients; i++ {
		counter := sendIncrementRequest(t, clients[i], uint64(i+1))
		expectedCounter += uint64(i + 1)
		if counter != expectedCounter {
			t.Fatalf("Expected counter to be %d, got %d", expectedCounter, counter)
		}
	}
}

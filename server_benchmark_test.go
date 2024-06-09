package quick

//import (
//	"github.com/bryanchriswhite/quick/server/cmd"
//	"net"
//	"sync"
//	"testing"
//	"time"
//)
//
//func startServerBenchmark(b *testing.B, wg *sync.WaitGroup) {
//	defer wg.Done()
//	main.main() // Start the server
//}
//
//func BenchmarkIncrementRequestResponse(b *testing.B) {
//	// Start the server in a separate goroutine
//	var wg sync.WaitGroup
//	wg.Add(1)
//	go startServerBenchmark(b, &wg)
//
//	// Give the server a moment to start
//	time.Sleep(time.Second)
//
//	// Connect to the server
//	conn, err := net.Dial("tcp", "localhost:8080")
//	if err != nil {
//		b.Fatalf("Failed to connect to server: %v", err)
//	}
//	defer conn.Close()
//
//	// Reset the benchmark timer
//	b.ResetTimer()
//
//	// Perform the benchmark
//	for i := 0; i < b.N; i++ {
//		sendIncrementRequest(b, conn, 1)
//	}
//
//	// Stop the benchmark timer
//	b.StopTimer()
//
//	// Clean up
//	wg.Done()
//}

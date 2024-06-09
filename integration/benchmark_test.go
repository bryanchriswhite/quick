package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bryanchriswhite/quick/client"
	"github.com/bryanchriswhite/quick/server"
)

func BenchmarkClient_Increment(b *testing.B) {
	ctx := context.Background()

	// Start the server.
	srv, err := server.NewServer(testServerURL)
	require.NoError(b, err)

	err = srv.Start(ctx)
	require.NoError(b, err)

	b.Cleanup(func() {
		err = srv.Close()
		require.NoError(b, err)
	})

	// Give the server a moment to start.
	time.Sleep(100 * time.Millisecond)

	// Connect to the server.
	client1, err := client.NewClient(testServerURL)
	require.NoError(b, err)

	b.Cleanup(func() {
		err = client1.Close()
		require.NoError(b, err)
	})

	// Reset the benchmark timer
	b.ResetTimer()

	// Perform the benchmark
	for i := 0; i < b.N; i++ {
		_, err := client1.Increment(10)
		require.NoError(b, err)
	}

	// Stop the benchmark timer
	b.StopTimer()
}

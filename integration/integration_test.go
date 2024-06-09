//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bryanchriswhite/quick/client"
	"github.com/bryanchriswhite/quick/server"
)

func TestSingleClient(t *testing.T) {
	ctx := context.Background()

	// Start the server.
	srv, err := server.NewServer(testServerURL)
	require.NoError(t, err)

	err = srv.Start(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		err = srv.Close()
		require.NoError(t, err)
	})

	// Give the server a moment to start.
	time.Sleep(100 * time.Millisecond)

	// Connect to the server.
	client1, err := client.NewClient(testServerURL)
	require.NoError(t, err)

	t.Cleanup(func() {
		err = client1.Close()
		require.NoError(t, err)
	})

	// Send increment requests and check the responses
	newCount, err := client1.Increment(1)
	require.NoError(t, err)
	require.Equal(t, newCount, uint64(1))

	newCount, err = client1.Increment(2)
	require.NoError(t, err)
	require.Equal(t, newCount, uint64(3))
}

func TestMultipleClients(t *testing.T) {
	ctx := context.Background()

	// Start the server.
	srv, err := server.NewServer(testServerURL)
	require.NoError(t, err)

	err = srv.Start(ctx)
	require.NoError(t, err)

	// Ensure the server closes cleanly when the test is finished.
	t.Cleanup(func() {
		err = srv.Close()
		require.NoError(t, err)
	})

	// Give the server a moment to start.
	time.Sleep(100 * time.Millisecond)

	// Connect to the server with client1.
	client1, err := client.NewClient(testServerURL)
	require.NoError(t, err)

	// Ensure the client1 closes cleanly when the test is finished.
	t.Cleanup(func() {
		//time.Sleep(time.Second)
		err = client1.Close()
		require.NoError(t, err)
	})

	// Connect to the server with client2.
	client2, err := client.NewClient(testServerURL)
	require.NoError(t, err)

	// Ensure the client1 closes cleanly when the test is finished.
	t.Cleanup(func() {
		err = client2.Close()
		require.NoError(t, err)
	})

	// Interleave send increment requests from client1 and client2, and check the responses.
	newCount, err := client1.Increment(1)
	require.NoError(t, err)
	require.Equal(t, newCount, uint64(1))

	newCount, err = client2.Increment(5)
	require.NoError(t, err)
	require.Equal(t, newCount, uint64(6))

	newCount, err = client1.Increment(2)
	require.NoError(t, err)
	require.Equal(t, newCount, uint64(8))

	newCount, err = client2.Increment(1)
	require.NoError(t, err)
	require.Equal(t, newCount, uint64(9))
}

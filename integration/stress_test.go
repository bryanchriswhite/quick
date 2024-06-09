package integration

import (
	"context"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bryanchriswhite/quick/client"
	"github.com/bryanchriswhite/quick/server"
)

const (
	// numClient is the number of concurrent clients which will
	// connect during the test.
	numClients = 100

	// numIncrements is the number of times Client#Increment()
	// will be called, with a random value, over the duration
	// of the test.
	numIncrements = 10000
)

func TestStress_Client_Increment(t *testing.T) {
	ctx := context.Background()

	// expectedCount accumulates all increment amounts passed to
	// Client#Increment() so that an assertion can be made about
	// the final count.
	var expectedCount atomic.Uint64

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

	// Create a wait group to wait for all clients to finish.
	wg := sync.WaitGroup{}

	// Set up test clients and kick off concurrent increment requests.
	for clientIdx := 0; clientIdx < numClients; clientIdx++ {
		// Connect to the server with incClient.
		incClient, err := client.NewClient(testServerURL)
		require.NoError(t, err)

		// Ensure the incClient closes cleanly when the test is finished.
		t.Cleanup(func() {
			//time.Sleep(time.Second)
			err = incClient.Close()
			require.NoError(t, err)
		})

		// Add numIncrement to the wait group.
		wg.Add(numIncrements)

		// Send increment requests from all clients concurrently.
		go func() {
			for incIdx := 0; incIdx < numIncrements; incIdx++ {
				// Generate a random increment amount.
				amount := uint64(rand.Intn(100))

				// Increment the expected count for future assertion.
				expectedCount.Add(amount)

				// Send the increment request.
				_, err := incClient.Increment(amount)
				require.NoError(t, err)

				// Done the wait group once for each increment.
				wg.Done()
			}
		}()
	}

	// Wait for all clients to finish.
	wg.Wait()

	incClient, err := client.NewClient(testServerURL)
	require.NoError(t, err)

	// Increment zero to get the final count.
	newCount, err := incClient.Increment(0)
	require.NoError(t, err)

	// Assert the final count matches the expectation.
	// NB: cast to int64 for better error format.
	require.Equal(t, int64(newCount), int64(expectedCount.Load()))
}

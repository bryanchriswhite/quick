package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	testTimeoutDuration = time.Second
	testServerURL       = "tcp://localhost:8080"
)

func TestServer_Start_ClosesWhenCtxIsCancelled(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())

	srv, err := NewServer(testServerURL)
	require.NoError(t, err)

	// srvClosedCh is intende to signal when the server has closed.
	srvClosedCh := make(chan struct{}, 1)

	go func() {
		// Server#Start() should block until the context is cancelled.
		err = srv.Start(ctx)
		require.NoError(t, err)

		for {
			if srv.IsClosed() {
				// Signal that the server has closed and break the loop.
				srvClosedCh <- struct{}{}
				break
			}
			// Sleep for a bit, no need to check more frequently.
			time.Sleep(10 * time.Millisecond)
		}

	}()

	// Cancel server context
	cancelCtx()

	// Wait for testTimeoutDuration OR srvClosedCh to receive.
	select {
	case <-time.After(testTimeoutDuration):
		t.Fatal("timed out waiting for server to close")
	case <-srvClosedCh:
		// Success
	}
}

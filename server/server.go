package server

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
	"sync/atomic"

	"github.com/bryanchriswhite/quick"
)

var _ quick.ServerI = (*Server)(nil)

type Server struct {
	// listenURL is used to start a net.Listener using the scheme as the "network"
	// and the host as the "address".
	listenURL *url.URL

	// listener is assigned in #Start() and its presence (or absence) is used to
	// determine whether the server is closed or not.
	listener net.Listener

	// counter is what is incremented by the server when it handles an increment request.
	counter atomic.Uint64
}

// connectionHandlerFn is a function which handles a connection's I/O.
type connectionHandlerFn func(*bufio.ReadWriter) error

// NewServer returns a new server which will listen on listenURL when #Start()
// is called.
func NewServer(listenURL string) (*Server, error) {
	// Parse the listenURL string into a URL.
	lURL, err := url.Parse(listenURL)
	if err != nil {
		return nil, err
	}

	return &Server{listenURL: lURL}, nil
}

// Start instantiates a listner using s#listenURL's scheme as the "network" (e.g. "tcp")
// and its host as the "address" (i.e. hostname/IP and port).
// It is NOT blocking.
func (s *Server) Start(ctx context.Context) (err error) {
	// Start the server listener
	s.listener, err = net.Listen(s.listenURL.Scheme, s.listenURL.Host)
	if err != nil {
		// Ensure listener is nil if there was an error, so that #IsClosed() will return true.
		s.listener = nil
		return err
	}

	// Ensure the server closes when the context is cancelled.
	go s.goCloseOnCtxDone(ctx)

	// Accept incoming connections concurrently.
	go s.goAcceptConnections(ctx)

	return nil
}

// IsClosed returns true if a listener is assigned to s#listener, which should be
// reset to nil in #Close().
func (s *Server) IsClosed() bool {
	return s.listener == nil
}

// Close closes the underlying listener and resets s#listener to nil.
func (s *Server) Close() error {
	err := s.listener.Close()

	// Reset s.listener to nil.
	s.listener = nil

	return err
}

// goCloseOnCtxDone closes the server when the context is cancelled.
// It is intended to be run in a goroutine.
func (s *Server) goCloseOnCtxDone(ctx context.Context) {
	// Block until the context is cancelled.
	<-ctx.Done()

	err := s.Close()
	if err != nil {
		// TODO_IMPROVE: replace `log` usage with
		// `github.com/pokt-network/poktroll/pkg/polylog`
		// structured logger.
		log.Printf("context done; error while closing server: %v", err)
	}
}

// goHandleConnection reads an increment request, increments s#counter, and responds
// with the new counter value. It is intended to be run in a goroutine.
func (s *Server) goHandleConnection(conn net.Conn, handler connectionHandlerFn) {
	defer func() {
		if err := conn.Close(); err != nil {
			// TODO_IMPROVE: replace `log` usage with
			// `github.com/pokt-network/poktroll/pkg/polylog`
			// structured logger.
			log.Printf("connection handled; error closing connection: %v", err)
		}
	}()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	readWriter := bufio.NewReadWriter(reader, writer)

	// Continuously read from the connection to handle multiple requests, in order,
	// over the same connection.
	for {
		// Return from this go routine if the server is closed.
		if s.IsClosed() {
			return
		}

		// TODO_IMPROVE: ideally, handler() would block until data is received instead
		// of returning an EOF error. We seem to get EOF errors regardless of whether
		// the connection is used directly or wrapped in bufio
		if err := handler(readWriter); err != nil && !errors.Is(err, io.EOF) {
			// TODO_IMPROVE: replace `log` usage with
			// `github.com/pokt-network/poktroll/pkg/polylog`
			// structured logger.
			log.Printf("failed to handle request: %v", err)
			// TODO: remove the else case...
			//} else {
			//	fmt.Printf("xxx EOF error while handling")
		}
	}
}

// goAcceptConnections continuously accepts incoming connections from s#listener and
// handles each concurrently. It is intended to be run in a goroutine.
func (s *Server) goAcceptConnections(ctx context.Context) {
	for {
		// Return from this go routine if the server is closed.
		if s.IsClosed() {
			return
		}

		// Block until a new incoming connection is received.
		conn, err := s.listener.Accept()
		if err != nil {
			// TODO_IMPROVE: figure out why s.IsClosed() can return true above,
			// yet we see this error: "use of closed network connection".
			if !strings.Contains(err.Error(), "use of closed network connection") {
				// TODO_IMPROVE: replace `log` usage with
				// `github.com/pokt-network/poktroll/pkg/polylog`
				// structured logger.
				log.Printf("error accepting connection: %v", err)
			}
			continue
		}

		// Handle each connection concurrently such that none blocks any other.
		//
		// TODO_CONSIDERATION: apply a concurrency limiter to prevent too many simultaneous connections.
		go s.goHandleConnection(conn, s.incrementHandler)
	}
}

// incrementHandler reads an increment request, increments s#counter, and responds
// with the new counter value. It should block until data is received.
func (s *Server) incrementHandler(readWriter *bufio.ReadWriter) error {
	// Read the request
	var incAmount uint64
	if err := binary.Read(readWriter, binary.LittleEndian, &incAmount); err != nil {
		return fmt.Errorf("reading from client connection: %w", err)
	}

	// Increment the global counter
	newCounter := s.counter.Add(incAmount)

	// Send the newCounter as the response
	if err := binary.Write(readWriter, binary.LittleEndian, newCounter); err != nil {
		return fmt.Errorf("writing to client connection: %w", err)
	}

	// Flush the write buffer (send buffered data over the network).
	if err := readWriter.Flush(); err != nil {
		return fmt.Errorf("flushing client connection writer: %w", err)
	}

	return nil
}

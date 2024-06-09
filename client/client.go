package client

import (
	"bufio"
	"encoding/binary"
	"net"
	"net/url"

	"github.com/bryanchriswhite/quick"
)

var _ quick.IncrementClient = (*Client)(nil)

type Client struct {
	conn net.Conn
}

// NewClient returns a new client which will send increment requests to remoteURL
// when #Increment() is called.
//
// TODO_CONSIDERATION: add #Dial() to separate construction from dialing.
func NewClient(remoteURL string) (*Client, error) {
	// Parse the remoteURL string into a URL.
	rURL, err := url.Parse(remoteURL)
	if err != nil {
		return nil, err
	}

	// Connect to the server
	conn, err := net.Dial(rURL.Scheme, rURL.Host)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

// Increment sends an increment request over the connection which was established
// during construction.
func (c *Client) Increment(amount uint64) (uint64, error) {
	reader := bufio.NewReader(c.conn)
	writer := bufio.NewWriter(c.conn)

	// Send the request
	err := binary.Write(writer, binary.LittleEndian, amount)
	if err != nil {
		return 0, err
	}

	// Flush the write buffer (send buffered data over the network).
	if err := writer.Flush(); err != nil {
		return 0, err
	}

	// Read the response
	var newCount uint64
	err = binary.Read(reader, binary.LittleEndian, &newCount)
	if err != nil {
		return 0, err
	}

	return newCount, nil
}

// Close closes the underlying connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	msgpackrpc "github.com/hashicorp/net-rpc-msgpackrpc/v2"
	"github.com/hashicorp/yamux"

	"github.com/hashicorp/hcp-scada-provider/internal/client/dialer"
)

const (
	// ClientPreamble is the preamble to send before upgrading
	// the connection into a SCADA version 1 connection.
	ClientPreamble = "SCADA 1\n"

	// rpcTimeout is how long of a read deadline we provide
	rpcTimeout = 10 * time.Second
)

// Opts is used to parameterize a Dial.
type Opts struct {
	// Dialer is the Dialer used to make a network connection.
	Dialer dialer.Dialer

	// Modifies the log output
	LogOutput io.Writer
}

// Client is a SCADA compatible client. This is a bare bones client that
// only handles the framing and RPC protocol. Higher-level clients should
// be preferred.
type Client struct {
	conn   net.Conn
	client *yamux.Session

	closed     bool
	closedLock sync.Mutex
}

// DialOpts is a parameterized Dial.
func DialOpts(target string, opts *Opts) (*Client, error) {
	conn, err := opts.Dialer.Dial(target)
	if err != nil {
		return nil, err
	}
	return initClient(conn, opts)
}

// DialOptsContext is a parameterized Dial.
func DialOptsContext(ctx context.Context, target string, opts *Opts) (*Client, error) {
	conn, err := opts.Dialer.DialContext(ctx, target)
	if err != nil {
		return nil, err
	}
	return initClient(conn, opts)
}

// initClient does the common initialization.
func initClient(conn net.Conn, opts *Opts) (*Client, error) {
	// Send the preamble
	_, err := conn.Write([]byte(ClientPreamble))
	if err != nil {
		return nil, fmt.Errorf("preamble write failed: %v", err)
	}

	// Wrap the connection in yamux for multiplexing
	ymConf := yamux.DefaultConfig()
	if opts.LogOutput != nil {
		ymConf.LogOutput = opts.LogOutput
	}
	client, err := yamux.Client(conn, ymConf)
	if err != nil {
		return nil, fmt.Errorf("failed to create yamux client: %w", err)
	}
	_, err = client.Ping()
	if err != nil {
		return nil, fmt.Errorf("yamux ping failed: %w", err)
	}

	// Create the client
	c := &Client{
		conn:   conn,
		client: client,
	}
	return c, nil
}

// Close is used to terminate the client connection.
func (c *Client) Close() error {
	c.closedLock.Lock()
	defer c.closedLock.Unlock()

	if c.closed {
		return nil
	}
	c.closed = true
	_ = c.client.GoAway() // Notify the other side of the close
	return c.client.Close()
}

// RPC performs a RPC call.
func (c *Client) RPC(method string, args interface{}, resp interface{}) error {
	// Get a stream
	stream, err := c.Open()
	if err != nil {
		return fmt.Errorf("failed to open stream: %w", err)
	}
	defer stream.Close()
	_ = stream.SetDeadline(time.Now().Add(rpcTimeout))

	// Create the RPC client
	cc := msgpackrpc.NewCodec(true, true, stream)
	return msgpackrpc.CallWithCodec(cc, method, args, resp)
}

// Accept is used to accept an incoming connection.
func (c *Client) Accept() (net.Conn, error) {
	return c.client.Accept()
}

// Open is used to open an outgoing connection.
func (c *Client) Open() (net.Conn, error) {
	return c.client.Open()
}

// Addr is so that client can act like a net.Listener.
func (c *Client) Addr() net.Addr {
	return c.client.LocalAddr()
}

// NumStreams returns the number of open streams on the client.
func (c *Client) NumStreams() int {
	return c.client.NumStreams()
}

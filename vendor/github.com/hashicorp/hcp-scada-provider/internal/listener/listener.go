// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// listener is a capability that pushes the received connection into a
// net.Listener. This capability allows exposing services that wrap a
// listener such as gRPC or HTTP Servers very easily.
package listener

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/hashicorp/hcp-scada-provider/capability"
)

// New creates a new SCADA listener using the given configuration.
// Requests for the HTTP capability are passed off to the listener that is
// returned.
func New(serviceID string) (Provider, net.Listener, error) {
	// Create a listener and handler
	list := newScadaListener(serviceID)
	provider := func(capability string, meta map[string]string,
		conn net.Conn) error {
		return list.Push(conn)
	}

	return provider, list, nil
}

// scadaListener returns a net.Listener for incoming SCADA connections.
type scadaListener struct {
	addr    *capability.Addr
	pending chan net.Conn

	closed   bool
	closedCh chan struct{}
	l        sync.Mutex
}

// newScadaListener returns a new listener.
func newScadaListener(cap string) *scadaListener {
	l := &scadaListener{
		addr:     capability.NewAddr(cap),
		pending:  make(chan net.Conn),
		closedCh: make(chan struct{}),
	}
	return l
}

// Push is used to add a connection to the queue.
func (s *scadaListener) Push(conn net.Conn) error {
	select {
	case s.pending <- conn:
		return nil
	case <-time.After(time.Second):
		return fmt.Errorf("accept timed out")
	case <-s.closedCh:
		return fmt.Errorf("scada listener closed")
	}
}

func (s *scadaListener) Accept() (net.Conn, error) {
	select {
	case conn := <-s.pending:
		return conn, nil
	case <-s.closedCh:
		return nil, fmt.Errorf("scada listener closed")
	}
}

func (s *scadaListener) Close() error {
	s.l.Lock()
	defer s.l.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true
	close(s.closedCh)
	return nil
}

func (s *scadaListener) Addr() net.Addr {
	return s.addr
}

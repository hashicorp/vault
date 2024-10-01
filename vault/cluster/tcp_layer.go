// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cluster

import (
	"crypto/tls"
	"net"
	"sync"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"go.uber.org/atomic"
)

// TCPLayer implements the NetworkLayer interface and uses TCP as the underlying
// network.
type TCPLayer struct {
	listeners []NetworkListener
	addrs     []*net.TCPAddr
	logger    log.Logger

	l       sync.Mutex
	stopped *atomic.Bool
}

// NewTCPLayer returns a TCPLayer.
func NewTCPLayer(addrs []*net.TCPAddr, logger log.Logger) *TCPLayer {
	return &TCPLayer{
		addrs:   addrs,
		logger:  logger,
		stopped: atomic.NewBool(false),
	}
}

// Addrs implements NetworkLayer.
func (l *TCPLayer) Addrs() []net.Addr {
	l.l.Lock()
	defer l.l.Unlock()

	if len(l.addrs) == 0 {
		return nil
	}

	ret := make([]net.Addr, len(l.addrs))
	for i, a := range l.addrs {
		ret[i] = a
	}

	return ret
}

// Listeners implements NetworkLayer. It starts a new TCP listener for each
// configured address.
func (l *TCPLayer) Listeners() []NetworkListener {
	l.l.Lock()
	defer l.l.Unlock()

	if l.listeners != nil {
		return l.listeners
	}

	listeners := []NetworkListener{}
	for i, laddr := range l.addrs {
		if l.logger.IsInfo() {
			l.logger.Info("starting listener", "listener_address", laddr)
		}

		tcpLn, err := net.ListenTCP("tcp", laddr)
		if err != nil {
			l.logger.Error("error starting listener", "error", err)
			continue
		}
		if laddr.String() != tcpLn.Addr().String() {
			// If we listened on port 0, record the port the OS gave us.
			l.addrs[i] = tcpLn.Addr().(*net.TCPAddr)
		}

		listeners = append(listeners, tcpLn)
	}

	l.listeners = listeners

	return listeners
}

// Dial implements the NetworkLayer interface.
func (l *TCPLayer) Dial(address string, timeout time.Duration, tlsConfig *tls.Config) (*tls.Conn, error) {
	dialer := &net.Dialer{
		Timeout: timeout,
	}
	return tls.DialWithDialer(dialer, "tcp", address, tlsConfig)
}

// Close implements the NetworkLayer interface.
func (l *TCPLayer) Close() error {
	if l.stopped.Swap(true) {
		return nil
	}
	l.l.Lock()
	defer l.l.Unlock()

	var retErr *multierror.Error
	for _, ln := range l.listeners {
		if err := ln.Close(); err != nil {
			retErr = multierror.Append(retErr, err)
		}
	}

	l.listeners = nil

	return retErr.ErrorOrNil()
}

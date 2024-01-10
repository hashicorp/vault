// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cluster

import (
	"io"
	"net"
	"time"

	uberAtomic "go.uber.org/atomic"
)

type delayedConn struct {
	net.Conn
	dr *delayedReader
}

func newDelayedConn(conn net.Conn, delay time.Duration) net.Conn {
	dr := &delayedReader{
		r:     conn,
		delay: uberAtomic.NewDuration(delay),
	}
	return &delayedConn{
		dr:   dr,
		Conn: conn,
	}
}

func (conn *delayedConn) Read(data []byte) (int, error) {
	return conn.dr.Read(data)
}

func (conn *delayedConn) SetDelay(delay time.Duration) {
	conn.dr.delay.Store(delay)
}

type delayedReader struct {
	r     io.Reader
	delay *uberAtomic.Duration
}

func (dr *delayedReader) Read(data []byte) (int, error) {
	// Sleep for the delay period prior to reading
	if delay := dr.delay.Load(); delay != 0 {
		time.Sleep(delay)
	}

	return dr.r.Read(data)
}

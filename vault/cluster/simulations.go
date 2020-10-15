package cluster

import (
	"io"
	"net"
	"time"
)

type delayedConn struct {
	net.Conn
	dr *delayedReader
}

func newDelayedConn(conn net.Conn, delay time.Duration) net.Conn {
	return &delayedConn{
		Conn: conn,
		dr: &delayedReader{
			r:     conn,
			delay: delay,
		},
	}
}

func (conn *delayedConn) Read(data []byte) (int, error) {
	return conn.dr.Read(data)
}

func (conn *delayedConn) SetDelay(delay time.Duration) {
	conn.dr.delay = delay
}

type delayedReader struct {
	r     io.Reader
	delay time.Duration
}

func (dr *delayedReader) Read(data []byte) (int, error) {
	// Sleep for the delay period prior to reading
	if dr.delay > 0 {
		time.Sleep(dr.delay)
	}

	return dr.r.Read(data)
}

package driver

import (
	"context"
	"crypto/tls"
	"database/sql/driver"
	"fmt"
	"io"
	"log/slog"
	"net"
	"runtime/pprof"
	"time"
)

var (
	cpuProfile = false
)

type dbConn interface {
	io.ReadWriteCloser
	lastRead() time.Time
	lastWrite() time.Time
}

type profileDBConn struct {
	dbConn
}

func (c *profileDBConn) Read(b []byte) (n int, err error) {
	pprof.Do(context.Background(), pprof.Labels("db", "read"), func(ctx context.Context) {
		n, err = c.dbConn.Read(b)
	})
	return
}

func (c *profileDBConn) Write(b []byte) (n int, err error) {
	pprof.Do(context.Background(), pprof.Labels("db", "write"), func(ctx context.Context) {
		n, err = c.dbConn.Write(b)
	})
	return
}

// stdDBConn wraps the database tcp connection. It sets timeouts and handles driver ErrBadConn behavior.
type stdDBConn struct {
	metrics    *metrics
	conn       net.Conn
	timeout    time.Duration
	logger     *slog.Logger
	_lastRead  time.Time
	_lastWrite time.Time
}

func newDBConn(ctx context.Context, logger *slog.Logger, host string, metrics *metrics, attrs *connAttrs) (dbConn, error) {
	conn, err := attrs.dialContext(ctx, host)
	if err != nil {
		return nil, err
	}
	// is TLS connection requested?
	if attrs._tlsConfig != nil {
		conn = tls.Client(conn, attrs._tlsConfig)
	}

	dbConn := &stdDBConn{metrics: metrics, conn: conn, timeout: attrs._timeout, logger: logger}
	if cpuProfile {
		return &profileDBConn{dbConn: dbConn}, nil
	}
	return dbConn, nil
}

func (c *stdDBConn) lastRead() time.Time  { return c._lastRead }
func (c *stdDBConn) lastWrite() time.Time { return c._lastWrite }

func (c *stdDBConn) deadline() (deadline time.Time) {
	if c.timeout == 0 {
		return
	}
	return time.Now().Add(c.timeout)
}

func (c *stdDBConn) Close() error { return c.conn.Close() }

// Read implements the io.Reader interface.
func (c *stdDBConn) Read(b []byte) (int, error) {
	// set timeout
	if err := c.conn.SetReadDeadline(c.deadline()); err != nil {
		return 0, fmt.Errorf("%w: %w", driver.ErrBadConn, err)
	}
	c._lastRead = time.Now()
	n, err := c.conn.Read(b)
	c.metrics.msgCh <- timeMsg{idx: timeRead, d: time.Since(c._lastRead)}
	c.metrics.msgCh <- counterMsg{idx: counterBytesRead, v: uint64(n)} //nolint:gosec
	if err != nil {
		c.logger.LogAttrs(context.Background(), slog.LevelError, "DB conn read error", slog.String("error", err.Error()), slog.String("local address", c.conn.LocalAddr().String()), slog.String("remote address", c.conn.RemoteAddr().String()))
		// wrap error in driver.ErrBadConn
		return n, fmt.Errorf("%w: %w", driver.ErrBadConn, err)
	}
	return n, nil
}

// Write implements the io.Writer interface.
func (c *stdDBConn) Write(b []byte) (int, error) {
	// set timeout
	if err := c.conn.SetWriteDeadline(c.deadline()); err != nil {
		return 0, fmt.Errorf("%w: %w", driver.ErrBadConn, err)
	}
	c._lastWrite = time.Now()
	n, err := c.conn.Write(b)
	c.metrics.msgCh <- timeMsg{idx: timeWrite, d: time.Since(c._lastWrite)}
	c.metrics.msgCh <- counterMsg{idx: counterBytesWritten, v: uint64(n)} //nolint:gosec
	if err != nil {
		c.logger.LogAttrs(context.Background(), slog.LevelError, "DB conn write error", slog.String("error", err.Error()), slog.String("local address", c.conn.LocalAddr().String()), slog.String("remote address", c.conn.RemoteAddr().String()))
		// wrap error in driver.ErrBadConn
		return n, fmt.Errorf("%w: %w", driver.ErrBadConn, err)
	}
	return n, nil
}

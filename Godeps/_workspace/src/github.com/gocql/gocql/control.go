package gocql

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

type controlConn struct {
	session *Session

	conn       atomic.Value
	connecting uint64

	retry RetryPolicy

	quit chan struct{}
}

func createControlConn(session *Session) *controlConn {
	control := &controlConn{
		session: session,
		quit:    make(chan struct{}),
		retry:   &SimpleRetryPolicy{NumRetries: 3},
	}

	control.conn.Store((*Conn)(nil))
	go control.heartBeat()

	return control
}

func (c *controlConn) heartBeat() {
	for {
		select {
		case <-c.quit:
			return
		case <-time.After(5 * time.Second):
		}

		resp, err := c.writeFrame(&writeOptionsFrame{})
		if err != nil {
			goto reconn
		}

		switch resp.(type) {
		case *supportedFrame:
			continue
		case error:
			goto reconn
		default:
			panic(fmt.Sprintf("gocql: unknown frame in response to options: %T", resp))
		}

	reconn:
		c.reconnect(true)
		// time.Sleep(5 * time.Second)
		continue

	}
}

func (c *controlConn) reconnect(refreshring bool) {
	if !atomic.CompareAndSwapUint64(&c.connecting, 0, 1) {
		return
	}

	success := false
	defer func() {
		// debounce reconnect a little
		if success {
			go func() {
				time.Sleep(500 * time.Millisecond)
				atomic.StoreUint64(&c.connecting, 0)
			}()
		} else {
			atomic.StoreUint64(&c.connecting, 0)
		}
	}()

	oldConn := c.conn.Load().(*Conn)

	// TODO: should have our own roundrobbin for hosts so that we can try each
	// in succession and guantee that we get a different host each time.
	host, conn := c.session.pool.Pick(nil)
	if conn == nil {
		return
	}

	newConn, err := Connect(conn.addr, conn.cfg, c, c.session)
	if err != nil {
		host.Mark(err)
		// TODO: add log handler for things like this
		return
	}

	host.Mark(nil)
	c.conn.Store(newConn)
	success = true

	if oldConn != nil {
		oldConn.Close()
	}

	if refreshring && c.session.cfg.DiscoverHosts {
		c.session.hostSource.refreshRing()
	}
}

func (c *controlConn) HandleError(conn *Conn, err error, closed bool) {
	if !closed {
		return
	}

	oldConn := c.conn.Load().(*Conn)
	if oldConn != conn {
		return
	}

	c.reconnect(true)
}

func (c *controlConn) writeFrame(w frameWriter) (frame, error) {
	conn := c.conn.Load().(*Conn)
	if conn == nil {
		return nil, errNoControl
	}

	framer, err := conn.exec(w, nil)
	if err != nil {
		return nil, err
	}

	return framer.parseFrame()
}

func (c *controlConn) withConn(fn func(*Conn) *Iter) *Iter {
	const maxConnectAttempts = 5
	connectAttempts := 0

	for i := 0; i < maxConnectAttempts; i++ {
		conn := c.conn.Load().(*Conn)
		if conn == nil {
			if connectAttempts > maxConnectAttempts {
				break
			}

			connectAttempts++

			c.reconnect(false)
			continue
		}

		return fn(conn)
	}

	return &Iter{err: errNoControl}
}

// query will return nil if the connection is closed or nil
func (c *controlConn) query(statement string, values ...interface{}) (iter *Iter) {
	q := c.session.Query(statement, values...).Consistency(One)

	for {
		iter = c.withConn(func(conn *Conn) *Iter {
			return conn.executeQuery(q)
		})

		q.attempts++
		if iter.err == nil || !c.retry.Attempt(q) {
			break
		}
	}

	return
}

func (c *controlConn) awaitSchemaAgreement() error {
	return c.withConn(func(conn *Conn) *Iter {
		return &Iter{err: conn.awaitSchemaAgreement()}
	}).err
}

func (c *controlConn) addr() string {
	conn := c.conn.Load().(*Conn)
	if conn == nil {
		return ""
	}
	return conn.addr
}

func (c *controlConn) close() {
	// TODO: handle more gracefully
	close(c.quit)
}

var errNoControl = errors.New("gocql: no control connection available")

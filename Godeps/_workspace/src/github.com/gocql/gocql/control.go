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
	control.reconnect()
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
		c.reconnect()
		time.Sleep(5 * time.Second)
		continue

	}
}

func (c *controlConn) reconnect() {
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
	conn := c.session.pool.Pick(nil)
	if conn == nil {
		return
	}

	newConn, err := Connect(conn.addr, conn.cfg, c)
	if err != nil {
		// TODO: add log handler for things like this
		return
	}

	c.conn.Store(newConn)
	success = true

	if oldConn != nil {
		oldConn.Close()
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

	c.reconnect()
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

// query will return nil if the connection is closed or nil
func (c *controlConn) query(statement string, values ...interface{}) (iter *Iter) {
	q := c.session.Query(statement, values...).Consistency(One)

	const maxConnectAttempts = 5
	connectAttempts := 0

	for {
		conn := c.conn.Load().(*Conn)
		if conn == nil {
			if connectAttempts > maxConnectAttempts {
				return &Iter{err: errNoControl}
			}

			connectAttempts++

			c.reconnect()
			continue
		}

		iter = conn.executeQuery(q)
		q.attempts++
		if iter.err == nil || !c.retry.Attempt(q) {
			break
		}
	}

	return
}

func (c *controlConn) awaitSchemaAgreement() (err error) {

	const (
		// TODO(zariel): if we export this make this configurable
		maxWaitTime = 60 * time.Second

		peerSchemas  = "SELECT schema_version FROM system.peers"
		localSchemas = "SELECT schema_version FROM system.local WHERE key='local'"
	)

	endDeadline := time.Now().Add(maxWaitTime)

	for time.Now().Before(endDeadline) {
		iter := c.query(peerSchemas)

		versions := make(map[string]struct{})

		var schemaVersion string
		for iter.Scan(&schemaVersion) {
			versions[schemaVersion] = struct{}{}
			schemaVersion = ""
		}

		if err = iter.Close(); err != nil {
			goto cont
		}

		iter = c.query(localSchemas)
		for iter.Scan(&schemaVersion) {
			versions[schemaVersion] = struct{}{}
			schemaVersion = ""
		}

		if err = iter.Close(); err != nil {
			goto cont
		}

		if len(versions) <= 1 {
			return nil
		}

	cont:
		time.Sleep(200 * time.Millisecond)
	}

	if err != nil {
		return
	}

	// not exported
	return errors.New("gocql: cluster schema versions not consistent")
}
func (c *controlConn) close() {
	// TODO: handle more gracefully
	close(c.quit)
}

var errNoControl = errors.New("gocql: no controll connection available")

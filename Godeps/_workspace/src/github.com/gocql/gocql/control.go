package gocql

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// Ensure that the atomic variable is aligned to a 64bit boundary
// so that atomic operations can be applied on 32bit architectures.
type controlConn struct {
	connecting int64

	session *Session
	conn    atomic.Value

	retry RetryPolicy

	closeWg sync.WaitGroup
	quit    chan struct{}
}

func createControlConn(session *Session) *controlConn {
	control := &controlConn{
		session: session,
		quit:    make(chan struct{}),
		retry:   &SimpleRetryPolicy{NumRetries: 3},
	}

	control.conn.Store((*Conn)(nil))

	return control
}

func (c *controlConn) heartBeat() {
	defer c.closeWg.Done()

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

func (c *controlConn) connect(endpoints []string) error {
	// intial connection attmept, try to connect to each endpoint to get an initial
	// list of nodes.

	// shuffle endpoints so not all drivers will connect to the same initial
	// node.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	perm := r.Perm(len(endpoints))
	shuffled := make([]string, len(endpoints))

	for i, endpoint := range endpoints {
		shuffled[perm[i]] = endpoint
	}

	// store that we are not connected so that reconnect wont happen if we error
	atomic.StoreInt64(&c.connecting, -1)

	var (
		conn *Conn
		err  error
	)

	for _, addr := range shuffled {
		conn, err = c.session.connect(JoinHostPort(addr, c.session.cfg.Port), c)
		if err != nil {
			log.Printf("gocql: unable to control conn dial %v: %v\n", addr, err)
			continue
		}

		if err = c.registerEvents(conn); err != nil {
			conn.Close()
			continue
		}

		// we should fetch the initial ring here and update initial host data. So that
		// when we return from here we have a ring topology ready to go.
		break
	}

	if conn == nil {
		// this is fatal, not going to connect a session
		return err
	}

	c.conn.Store(conn)
	atomic.StoreInt64(&c.connecting, 0)

	c.closeWg.Add(1)
	go c.heartBeat()

	return nil
}

func (c *controlConn) registerEvents(conn *Conn) error {
	framer, err := conn.exec(&writeRegisterFrame{
		events: []string{"TOPOLOGY_CHANGE", "STATUS_CHANGE", "STATUS_CHANGE"},
	}, nil)
	if err != nil {
		return err
	}

	frame, err := framer.parseFrame()
	if err != nil {
		return err
	} else if _, ok := frame.(*readyFrame); !ok {
		return fmt.Errorf("unexpected frame in response to register: got %T: %v\n", frame, frame)
	}

	return nil
}

func (c *controlConn) reconnect(refreshring bool) {
	// TODO: simplify this function, use session.ring to get hosts instead of the
	// connection pool
	if !atomic.CompareAndSwapInt64(&c.connecting, 0, 1) {
		return
	}

	success := false
	defer func() {
		// debounce reconnect a little
		if success {
			go func() {
				time.Sleep(500 * time.Millisecond)
				atomic.StoreInt64(&c.connecting, 0)
			}()
		} else {
			atomic.StoreInt64(&c.connecting, 0)
		}
	}()

	addr := c.addr()
	oldConn := c.conn.Load().(*Conn)
	if oldConn != nil {
		oldConn.Close()
	}

	var newConn *Conn
	if addr != "" {
		// try to connect to the old host
		conn, err := c.session.connect(addr, c)
		if err != nil {
			// host is dead
			// TODO: this is replicated in a few places
			ip, portStr, _ := net.SplitHostPort(addr)
			port, _ := strconv.Atoi(portStr)
			c.session.handleNodeDown(net.ParseIP(ip), port)
		} else {
			newConn = conn
		}
	}

	// TODO: should have our own roundrobbin for hosts so that we can try each
	// in succession and guantee that we get a different host each time.
	if newConn == nil {
		_, conn := c.session.pool.Pick(nil)
		if conn == nil {
			return
		}

		if conn == nil {
			return
		}

		var err error
		newConn, err = c.session.connect(conn.addr, c)
		if err != nil {
			// TODO: add log handler for things like this
			return
		}
	}

	if err := c.registerEvents(newConn); err != nil {
		// TODO: handle this case better
		newConn.Close()
		log.Printf("gocql: control unable to register events: %v\n", err)
		return
	}

	c.conn.Store(newConn)
	success = true

	if refreshring {
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

func (c *controlConn) fetchHostInfo(addr net.IP, port int) (*HostInfo, error) {
	// TODO(zariel): we should probably move this into host_source or atleast
	// share code with it.
	hostname, _, err := net.SplitHostPort(c.addr())
	if err != nil {
		return nil, fmt.Errorf("unable to fetch host info, invalid conn addr: %q: %v", c.addr(), err)
	}

	isLocal := hostname == addr.String()

	var fn func(*HostInfo) error

	if isLocal {
		fn = func(host *HostInfo) error {
			// TODO(zariel): should we fetch rpc_address from here?
			iter := c.query("SELECT data_center, rack, host_id, tokens, release_version FROM system.local WHERE key='local'")
			iter.Scan(&host.dataCenter, &host.rack, &host.hostId, &host.tokens, &host.version)
			return iter.Close()
		}
	} else {
		fn = func(host *HostInfo) error {
			// TODO(zariel): should we fetch rpc_address from here?
			iter := c.query("SELECT data_center, rack, host_id, tokens, release_version FROM system.peers WHERE peer=?", addr)
			iter.Scan(&host.dataCenter, &host.rack, &host.hostId, &host.tokens, &host.version)
			return iter.Close()
		}
	}

	host := &HostInfo{
		port: port,
	}

	if err := fn(host); err != nil {
		return nil, err
	}
	host.peer = addr.String()

	return host, nil
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
	c.closeWg.Wait()
	conn := c.conn.Load().(*Conn)
	if conn != nil {
		conn.Close()
	}
}

var errNoControl = errors.New("gocql: no control connection available")

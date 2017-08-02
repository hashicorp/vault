package gocql

import (
	"context"
	crand "crypto/rand"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"regexp"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	randr    *rand.Rand
	mutRandr sync.Mutex
)

func init() {
	b := make([]byte, 4)
	if _, err := crand.Read(b); err != nil {
		panic(fmt.Sprintf("unable to seed random number generator: %v", err))
	}

	randr = rand.New(rand.NewSource(int64(readInt(b))))
}

// Ensure that the atomic variable is aligned to a 64bit boundary
// so that atomic operations can be applied on 32bit architectures.
type controlConn struct {
	session *Session
	conn    atomic.Value

	retry RetryPolicy

	started int32
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
	if !atomic.CompareAndSwapInt32(&c.started, 0, 1) {
		return
	}

	sleepTime := 1 * time.Second

	for {
		select {
		case <-c.quit:
			return
		case <-time.After(sleepTime):
		}

		resp, err := c.writeFrame(&writeOptionsFrame{})
		if err != nil {
			goto reconn
		}

		switch resp.(type) {
		case *supportedFrame:
			// Everything ok
			sleepTime = 5 * time.Second
			continue
		case error:
			goto reconn
		default:
			panic(fmt.Sprintf("gocql: unknown frame in response to options: %T", resp))
		}

	reconn:
		// try to connect a bit faster
		sleepTime = 1 * time.Second
		c.reconnect(true)
		// time.Sleep(5 * time.Second)
		continue
	}
}

var hostLookupPreferV4 = false

func hostInfo(addr string, defaultPort int) (*HostInfo, error) {
	var port int
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr
		port = defaultPort
	} else {
		port, err = strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}
	}

	ip := net.ParseIP(host)
	if ip == nil {
		ips, err := net.LookupIP(host)
		if err != nil {
			return nil, err
		} else if len(ips) == 0 {
			return nil, fmt.Errorf("No IP's returned from DNS lookup for %q", addr)
		}

		if hostLookupPreferV4 {
			for _, v := range ips {
				if v4 := v.To4(); v4 != nil {
					ip = v4
					break
				}
			}
			if ip == nil {
				ip = ips[0]
			}
		} else {
			// TODO(zariel): should we check that we can connect to any of the ips?
			ip = ips[0]
		}

	}

	return &HostInfo{connectAddress: ip, port: port}, nil
}

func shuffleHosts(hosts []*HostInfo) []*HostInfo {
	mutRandr.Lock()
	perm := randr.Perm(len(hosts))
	mutRandr.Unlock()
	shuffled := make([]*HostInfo, len(hosts))

	for i, host := range hosts {
		shuffled[perm[i]] = host
	}

	return shuffled
}

func (c *controlConn) shuffleDial(endpoints []*HostInfo) (*Conn, error) {
	// shuffle endpoints so not all drivers will connect to the same initial
	// node.
	shuffled := shuffleHosts(endpoints)

	var err error
	for _, host := range shuffled {
		var conn *Conn
		conn, err = c.session.connect(host, c)
		if err == nil {
			return conn, nil
		}

		Logger.Printf("gocql: unable to dial control conn %v: %v\n", host.ConnectAddress(), err)
	}

	return nil, err
}

// this is going to be version dependant and a nightmare to maintain :(
var protocolSupportRe = regexp.MustCompile(`the lowest supported version is \d+ and the greatest is (\d+)$`)

func parseProtocolFromError(err error) int {
	// I really wish this had the actual info in the error frame...
	matches := protocolSupportRe.FindAllStringSubmatch(err.Error(), -1)
	if len(matches) != 1 || len(matches[0]) != 2 {
		if verr, ok := err.(*protocolError); ok {
			return int(verr.frame.Header().version.version())
		}
		return 0
	}

	max, err := strconv.Atoi(matches[0][1])
	if err != nil {
		return 0
	}

	return max
}

func (c *controlConn) discoverProtocol(hosts []*HostInfo) (int, error) {
	hosts = shuffleHosts(hosts)

	connCfg := *c.session.connCfg
	connCfg.ProtoVersion = 4 // TODO: define maxProtocol

	handler := connErrorHandlerFn(func(c *Conn, err error, closed bool) {
		// we should never get here, but if we do it means we connected to a
		// host successfully which means our attempted protocol version worked
	})

	var err error
	for _, host := range hosts {
		var conn *Conn
		conn, err = Connect(host, &connCfg, handler, c.session)
		if err == nil {
			conn.Close()
			return connCfg.ProtoVersion, nil
		}

		if proto := parseProtocolFromError(err); proto > 0 {
			return proto, nil
		}
	}

	return 0, err
}

func (c *controlConn) connect(hosts []*HostInfo) error {
	if len(hosts) == 0 {
		return errors.New("control: no endpoints specified")
	}

	conn, err := c.shuffleDial(hosts)
	if err != nil {
		return fmt.Errorf("control: unable to connect to initial hosts: %v", err)
	}

	if err := c.setupConn(conn); err != nil {
		conn.Close()
		return fmt.Errorf("control: unable to setup connection: %v", err)
	}

	// we could fetch the initial ring here and update initial host data. So that
	// when we return from here we have a ring topology ready to go.

	go c.heartBeat()

	return nil
}

func (c *controlConn) setupConn(conn *Conn) error {
	if err := c.registerEvents(conn); err != nil {
		conn.Close()
		return err
	}

	c.conn.Store(conn)

	if v, ok := conn.conn.RemoteAddr().(*net.TCPAddr); ok {
		c.session.handleNodeUp(copyBytes(v.IP), v.Port, false)
		return nil
	}

	host, portstr, err := net.SplitHostPort(conn.conn.RemoteAddr().String())
	if err != nil {
		return err
	}

	port, err := strconv.Atoi(portstr)
	if err != nil {
		return err
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return fmt.Errorf("invalid remote addr: addr=%v host=%q", conn.conn.RemoteAddr(), host)
	}

	c.session.handleNodeUp(ip, port, false)

	return nil
}

func (c *controlConn) registerEvents(conn *Conn) error {
	var events []string

	if !c.session.cfg.Events.DisableTopologyEvents {
		events = append(events, "TOPOLOGY_CHANGE")
	}
	if !c.session.cfg.Events.DisableNodeStatusEvents {
		events = append(events, "STATUS_CHANGE")
	}
	if !c.session.cfg.Events.DisableSchemaEvents {
		events = append(events, "SCHEMA_CHANGE")
	}

	if len(events) == 0 {
		return nil
	}

	framer, err := conn.exec(context.Background(),
		&writeRegisterFrame{
			events: events,
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

	var host *HostInfo
	oldConn := c.conn.Load().(*Conn)
	if oldConn != nil {
		host = oldConn.host
		oldConn.Close()
	}

	var newConn *Conn
	if host != nil {
		// try to connect to the old host
		conn, err := c.session.connect(host, c)
		if err != nil {
			// host is dead
			// TODO: this is replicated in a few places
			c.session.handleNodeDown(host.ConnectAddress(), host.Port())
		} else {
			newConn = conn
		}
	}

	// TODO: should have our own round-robin for hosts so that we can try each
	// in succession and guarantee that we get a different host each time.
	if newConn == nil {
		host := c.session.ring.rrHost()
		if host == nil {
			c.connect(c.session.ring.endpoints)
			return
		}

		var err error
		newConn, err = c.session.connect(host, c)
		if err != nil {
			// TODO: add log handler for things like this
			return
		}
	}

	if err := c.setupConn(newConn); err != nil {
		newConn.Close()
		Logger.Printf("gocql: control unable to register events: %v\n", err)
		return
	}

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

	framer, err := conn.exec(context.Background(), w, nil)
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
	q := c.session.Query(statement, values...).Consistency(One).RoutingKey([]byte{}).Trace(nil)

	for {
		iter = c.withConn(func(conn *Conn) *Iter {
			return conn.executeQuery(q)
		})

		if gocqlDebug && iter.err != nil {
			Logger.Printf("control: error executing %q: %v\n", statement, iter.err)
		}

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

func (c *controlConn) GetHostInfo() *HostInfo {
	conn := c.conn.Load().(*Conn)
	if conn == nil {
		return nil
	}
	return conn.host
}

func (c *controlConn) close() {
	if atomic.CompareAndSwapInt32(&c.started, 1, -1) {
		c.quit <- struct{}{}
	}
	conn := c.conn.Load().(*Conn)
	if conn != nil {
		conn.Close()
	}
}

var errNoControl = errors.New("gocql: no control connection available")

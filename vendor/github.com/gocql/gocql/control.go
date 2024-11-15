/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/*
 * Content before git sha 34fdeebefcbf183ed7f916f931aa0586fdaa1b40
 * Copyright (c) 2016, The Gocql authors,
 * provided under the BSD-3-Clause License.
 * See the NOTICE file distributed with this work for additional information.
 */

package gocql

import (
	"context"
	crand "crypto/rand"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
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

const (
	controlConnStarting = 0
	controlConnStarted  = 1
	controlConnClosing  = -1
)

// Ensure that the atomic variable is aligned to a 64bit boundary
// so that atomic operations can be applied on 32bit architectures.
type controlConn struct {
	state        int32
	reconnecting int32

	session *Session
	conn    atomic.Value

	retry RetryPolicy

	quit chan struct{}
}

func createControlConn(session *Session) *controlConn {
	control := &controlConn{
		session: session,
		quit:    make(chan struct{}),
		retry:   &SimpleRetryPolicy{NumRetries: 3},
	}

	control.conn.Store((*connHost)(nil))

	return control
}

func (c *controlConn) heartBeat() {
	if !atomic.CompareAndSwapInt32(&c.state, controlConnStarting, controlConnStarted) {
		return
	}

	sleepTime := 1 * time.Second
	timer := time.NewTimer(sleepTime)
	defer timer.Stop()

	for {
		timer.Reset(sleepTime)

		select {
		case <-c.quit:
			return
		case <-timer.C:
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
		c.reconnect()
		continue
	}
}

var hostLookupPreferV4 = os.Getenv("GOCQL_HOST_LOOKUP_PREFER_V4") == "true"

func hostInfo(addr string, defaultPort int) ([]*HostInfo, error) {
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

	var hosts []*HostInfo

	// Check if host is a literal IP address
	if ip := net.ParseIP(host); ip != nil {
		hosts = append(hosts, &HostInfo{hostname: host, connectAddress: ip, port: port})
		return hosts, nil
	}

	// Look up host in DNS
	ips, err := LookupIP(host)
	if err != nil {
		return nil, err
	} else if len(ips) == 0 {
		return nil, fmt.Errorf("no IP's returned from DNS lookup for %q", addr)
	}

	// Filter to v4 addresses if any present
	if hostLookupPreferV4 {
		var preferredIPs []net.IP
		for _, v := range ips {
			if v4 := v.To4(); v4 != nil {
				preferredIPs = append(preferredIPs, v4)
			}
		}
		if len(preferredIPs) != 0 {
			ips = preferredIPs
		}
	}

	for _, ip := range ips {
		hosts = append(hosts, &HostInfo{hostname: host, connectAddress: ip, port: port})
	}

	return hosts, nil
}

func shuffleHosts(hosts []*HostInfo) []*HostInfo {
	shuffled := make([]*HostInfo, len(hosts))
	copy(shuffled, hosts)

	mutRandr.Lock()
	randr.Shuffle(len(hosts), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	mutRandr.Unlock()

	return shuffled
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
		if !closed {
			c.Close()
		}
	})

	var err error
	for _, host := range hosts {
		var conn *Conn
		conn, err = c.session.dial(c.session.ctx, host, &connCfg, handler)
		if conn != nil {
			conn.Close()
		}

		if err == nil {
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

	// shuffle endpoints so not all drivers will connect to the same initial
	// node.
	hosts = shuffleHosts(hosts)

	cfg := *c.session.connCfg
	cfg.disableCoalesce = true

	var conn *Conn
	var err error
	for _, host := range hosts {
		conn, err = c.session.dial(c.session.ctx, host, &cfg, c)
		if err != nil {
			c.session.logger.Printf("gocql: unable to dial control conn %v:%v: %v\n", host.ConnectAddress(), host.Port(), err)
			continue
		}
		err = c.setupConn(conn)
		if err == nil {
			break
		}
		c.session.logger.Printf("gocql: unable setup control conn %v:%v: %v\n", host.ConnectAddress(), host.Port(), err)
		conn.Close()
		conn = nil
	}
	if conn == nil {
		return fmt.Errorf("unable to connect to initial hosts: %v", err)
	}

	// we could fetch the initial ring here and update initial host data. So that
	// when we return from here we have a ring topology ready to go.

	go c.heartBeat()

	return nil
}

type connHost struct {
	conn *Conn
	host *HostInfo
}

func (c *controlConn) setupConn(conn *Conn) error {
	// we need up-to-date host info for the filterHost call below
	iter := conn.querySystemLocal(context.TODO())
	host, err := c.session.hostInfoFromIter(iter, conn.host.connectAddress, conn.conn.RemoteAddr().(*net.TCPAddr).Port)
	if err != nil {
		return err
	}

	host = c.session.ring.addOrUpdate(host)

	if c.session.cfg.filterHost(host) {
		return fmt.Errorf("host was filtered: %v", host.ConnectAddress())
	}

	if err := c.registerEvents(conn); err != nil {
		return fmt.Errorf("register events: %v", err)
	}

	ch := &connHost{
		conn: conn,
		host: host,
	}

	c.conn.Store(ch)
	if c.session.initialized() {
		// We connected to control conn, so add the connect the host in pool as well.
		// Notify session we can start trying to connect to the node.
		// We can't start the fill before the session is initialized, otherwise the fill would interfere
		// with the fill called by Session.init. Session.init needs to wait for its fill to finish and that
		// would return immediately if we started the fill here.
		// TODO(martin-sucha): Trigger pool refill for all hosts, like in reconnectDownedHosts?
		go c.session.startPoolFill(host)
	}
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

func (c *controlConn) reconnect() {
	if atomic.LoadInt32(&c.state) == controlConnClosing {
		return
	}
	if !atomic.CompareAndSwapInt32(&c.reconnecting, 0, 1) {
		return
	}
	defer atomic.StoreInt32(&c.reconnecting, 0)

	conn, err := c.attemptReconnect()

	if conn == nil {
		c.session.logger.Printf("gocql: unable to reconnect control connection: %v\n", err)
		return
	}

	err = c.session.refreshRing()
	if err != nil {
		c.session.logger.Printf("gocql: unable to refresh ring: %v\n", err)
	}
}

func (c *controlConn) attemptReconnect() (*Conn, error) {
	hosts := c.session.ring.allHosts()
	hosts = shuffleHosts(hosts)

	// keep the old behavior of connecting to the old host first by moving it to
	// the front of the slice
	ch := c.getConn()
	if ch != nil {
		for i := range hosts {
			if hosts[i].Equal(ch.host) {
				hosts[0], hosts[i] = hosts[i], hosts[0]
				break
			}
		}
		ch.conn.Close()
	}

	conn, err := c.attemptReconnectToAnyOfHosts(hosts)

	if conn != nil {
		return conn, err
	}

	c.session.logger.Printf("gocql: unable to connect to any ring node: %v\n", err)
	c.session.logger.Printf("gocql: control falling back to initial contact points.\n")
	// Fallback to initial contact points, as it may be the case that all known initialHosts
	// changed their IPs while keeping the same hostname(s).
	initialHosts, resolvErr := addrsToHosts(c.session.cfg.Hosts, c.session.cfg.Port, c.session.logger)
	if resolvErr != nil {
		return nil, fmt.Errorf("resolve contact points' hostnames: %v", resolvErr)
	}

	return c.attemptReconnectToAnyOfHosts(initialHosts)
}

func (c *controlConn) attemptReconnectToAnyOfHosts(hosts []*HostInfo) (*Conn, error) {
	var conn *Conn
	var err error
	for _, host := range hosts {
		conn, err = c.session.connect(c.session.ctx, host, c)
		if err != nil {
			c.session.logger.Printf("gocql: unable to dial control conn %v:%v: %v\n", host.ConnectAddress(), host.Port(), err)
			continue
		}
		err = c.setupConn(conn)
		if err == nil {
			break
		}
		c.session.logger.Printf("gocql: unable setup control conn %v:%v: %v\n", host.ConnectAddress(), host.Port(), err)
		conn.Close()
		conn = nil
	}
	return conn, err
}

func (c *controlConn) HandleError(conn *Conn, err error, closed bool) {
	if !closed {
		return
	}

	oldConn := c.getConn()

	// If connection has long gone, and not been attempted for awhile,
	// it's possible to have oldConn as nil here (#1297).
	if oldConn != nil && oldConn.conn != conn {
		return
	}

	c.reconnect()
}

func (c *controlConn) getConn() *connHost {
	return c.conn.Load().(*connHost)
}

func (c *controlConn) writeFrame(w frameBuilder) (frame, error) {
	ch := c.getConn()
	if ch == nil {
		return nil, errNoControl
	}

	framer, err := ch.conn.exec(context.Background(), w, nil)
	if err != nil {
		return nil, err
	}

	return framer.parseFrame()
}

func (c *controlConn) withConnHost(fn func(*connHost) *Iter) *Iter {
	const maxConnectAttempts = 5
	connectAttempts := 0

	for i := 0; i < maxConnectAttempts; i++ {
		ch := c.getConn()
		if ch == nil {
			if connectAttempts > maxConnectAttempts {
				break
			}

			connectAttempts++

			c.reconnect()
			continue
		}

		return fn(ch)
	}

	return &Iter{err: errNoControl}
}

func (c *controlConn) withConn(fn func(*Conn) *Iter) *Iter {
	return c.withConnHost(func(ch *connHost) *Iter {
		return fn(ch.conn)
	})
}

// query will return nil if the connection is closed or nil
func (c *controlConn) query(statement string, values ...interface{}) (iter *Iter) {
	q := c.session.Query(statement, values...).Consistency(One).RoutingKey([]byte{}).Trace(nil)

	for {
		iter = c.withConn(func(conn *Conn) *Iter {
			// we want to keep the query on the control connection
			q.conn = conn
			return conn.executeQuery(context.TODO(), q)
		})

		if gocqlDebug && iter.err != nil {
			c.session.logger.Printf("control: error executing %q: %v\n", statement, iter.err)
		}

		q.AddAttempts(1, c.getConn().host)
		if iter.err == nil || !c.retry.Attempt(q) {
			break
		}
	}

	return
}

func (c *controlConn) awaitSchemaAgreement() error {
	return c.withConn(func(conn *Conn) *Iter {
		return &Iter{err: conn.awaitSchemaAgreement(context.TODO())}
	}).err
}

func (c *controlConn) close() {
	if atomic.CompareAndSwapInt32(&c.state, controlConnStarted, controlConnClosing) {
		c.quit <- struct{}{}
	}

	ch := c.getConn()
	if ch != nil {
		ch.conn.Close()
	}
}

var errNoControl = errors.New("gocql: no control connection available")

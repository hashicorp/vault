// Copyright (c) 2012 The gocql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocql

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

/*ConnectionPool represents the interface gocql will use to work with a collection of connections.

Purpose

The connection pool in gocql opens and closes connections as well as selects an available connection
for gocql to execute a query against. The pool is also respnsible for handling connection errors that
are caught by the connection experiencing the error.

A connection pool should make a copy of the variables used from the ClusterConfig provided to the pool
upon creation. ClusterConfig is a pointer and can be modified after the creation of the pool. This can
lead to issues with variables being modified outside the expectations of the ConnectionPool type.

Example of Single Connection Pool:

	type SingleConnection struct {
		conn *Conn
		cfg *ClusterConfig
	}

	func NewSingleConnection(cfg *ClusterConfig) ConnectionPool {
		addr := JoinHostPort(cfg.Hosts[0], cfg.Port)

		connCfg := ConnConfig{
			ProtoVersion:  cfg.ProtoVersion,
			CQLVersion:    cfg.CQLVersion,
			Timeout:       cfg.Timeout,
			NumStreams:    cfg.NumStreams,
			Compressor:    cfg.Compressor,
			Authenticator: cfg.Authenticator,
			Keepalive:     cfg.SocketKeepalive,
		}
		pool := SingleConnection{cfg:cfg}
		pool.conn = Connect(addr,connCfg,pool)
		return &pool
	}

	func (s *SingleConnection) HandleError(conn *Conn, err error, closed bool) {
		if closed {
			connCfg := ConnConfig{
				ProtoVersion:  cfg.ProtoVersion,
				CQLVersion:    cfg.CQLVersion,
				Timeout:       cfg.Timeout,
				NumStreams:    cfg.NumStreams,
				Compressor:    cfg.Compressor,
				Authenticator: cfg.Authenticator,
				Keepalive:     cfg.SocketKeepalive,
			}
			s.conn = Connect(conn.Address(),connCfg,s)
		}
	}

	func (s *SingleConnection) Pick(qry *Query) *Conn {
		if s.conn.isClosed {
			return nil
		}
		return s.conn
	}

	func (s *SingleConnection) Size() int {
		return 1
	}

	func (s *SingleConnection) Close() {
		s.conn.Close()
	}

This is a very simple example of a type that exposes the connection pool interface. To assign
this type as the connection pool to use you would assign it to the ClusterConfig like so:

		cluster := NewCluster("127.0.0.1")
		cluster.ConnPoolType = NewSingleConnection
		...
		session, err := cluster.CreateSession()

To see a more complete example of a ConnectionPool implementation please see the SimplePool type.
*/
type ConnectionPool interface {
	SetHosts
	Pick(*Query) *Conn
	Size() int
	Close()
}

// interface to implement to receive the host information
type SetHosts interface {
	SetHosts(hosts []HostInfo)
}

// interface to implement to receive the partitioner value
type SetPartitioner interface {
	SetPartitioner(partitioner string)
}

//NewPoolFunc is the type used by ClusterConfig to create a pool of a specific type.
type NewPoolFunc func(*ClusterConfig) (ConnectionPool, error)

//SimplePool is the current implementation of the connection pool inside gocql. This
//pool is meant to be a simple default used by gocql so users can get up and running
//quickly.
type SimplePool struct {
	cfg      *ClusterConfig
	hostPool *RoundRobin
	connPool map[string]*RoundRobin
	conns    map[*Conn]struct{}
	keyspace string

	hostMu sync.RWMutex
	// this is the set of current hosts which the pool will attempt to connect to
	hosts map[string]*HostInfo

	// protects hostpool, connPoll, conns, quit
	mu sync.Mutex

	cFillingPool chan int

	quit     bool
	quitWait chan bool
	quitOnce sync.Once

	tlsConfig *tls.Config
}

func setupTLSConfig(sslOpts *SslOptions) (*tls.Config, error) {
	// ca cert is optional
	if sslOpts.CaPath != "" {
		if sslOpts.RootCAs == nil {
			sslOpts.RootCAs = x509.NewCertPool()
		}

		pem, err := ioutil.ReadFile(sslOpts.CaPath)
		if err != nil {
			return nil, fmt.Errorf("connectionpool: unable to open CA certs: %v", err)
		}

		if !sslOpts.RootCAs.AppendCertsFromPEM(pem) {
			return nil, errors.New("connectionpool: failed parsing or CA certs")
		}
	}

	if sslOpts.CertPath != "" || sslOpts.KeyPath != "" {
		mycert, err := tls.LoadX509KeyPair(sslOpts.CertPath, sslOpts.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("connectionpool: unable to load X509 key pair: %v", err)
		}
		sslOpts.Certificates = append(sslOpts.Certificates, mycert)
	}

	sslOpts.InsecureSkipVerify = !sslOpts.EnableHostVerification

	return &sslOpts.Config, nil
}

//NewSimplePool is the function used by gocql to create the simple connection pool.
//This is the default if no other pool type is specified.
func NewSimplePool(cfg *ClusterConfig) (ConnectionPool, error) {
	pool := &SimplePool{
		cfg:          cfg,
		hostPool:     NewRoundRobin(),
		connPool:     make(map[string]*RoundRobin),
		conns:        make(map[*Conn]struct{}),
		quitWait:     make(chan bool),
		cFillingPool: make(chan int, 1),
		keyspace:     cfg.Keyspace,
		hosts:        make(map[string]*HostInfo),
	}

	for _, host := range cfg.Hosts {
		// seed hosts have unknown topology
		// TODO: Handle populating this during SetHosts
		pool.hosts[host] = &HostInfo{Peer: host}
	}

	if cfg.SslOpts != nil {
		config, err := setupTLSConfig(cfg.SslOpts)
		if err != nil {
			return nil, err
		}
		pool.tlsConfig = config
	}

	//Walk through connecting to hosts. As soon as one host connects
	//defer the remaining connections to cluster.fillPool()
	for i := 0; i < len(cfg.Hosts); i++ {
		addr := JoinHostPort(cfg.Hosts[i], cfg.Port)

		if pool.connect(addr) == nil {
			pool.cFillingPool <- 1
			go pool.fillPool()
			break
		}
	}

	return pool, nil
}

func (c *SimplePool) connect(addr string) error {

	cfg := ConnConfig{
		ProtoVersion:  c.cfg.ProtoVersion,
		CQLVersion:    c.cfg.CQLVersion,
		Timeout:       c.cfg.Timeout,
		NumStreams:    c.cfg.NumStreams,
		Compressor:    c.cfg.Compressor,
		Authenticator: c.cfg.Authenticator,
		Keepalive:     c.cfg.SocketKeepalive,
		tlsConfig:     c.tlsConfig,
	}

	conn, err := Connect(addr, cfg, c)
	if err != nil {
		log.Printf("connect: failed to connect to %q: %v", addr, err)
		return err
	}

	return c.addConn(conn)
}

func (c *SimplePool) addConn(conn *Conn) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.quit {
		conn.Close()
		return nil
	}

	//Set the connection's keyspace if any before adding it to the pool
	if c.keyspace != "" {
		if err := conn.UseKeyspace(c.keyspace); err != nil {
			log.Printf("error setting connection keyspace. %v", err)
			conn.Close()
			return err
		}
	}

	connPool := c.connPool[conn.Address()]
	if connPool == nil {
		connPool = NewRoundRobin()
		c.connPool[conn.Address()] = connPool
		c.hostPool.AddNode(connPool)
	}

	connPool.AddNode(conn)
	c.conns[conn] = struct{}{}

	return nil
}

//fillPool manages the pool of connections making sure that each host has the correct
//amount of connections defined. Also the method will test a host with one connection
//instead of flooding the host with number of connections defined in the cluster config
func (c *SimplePool) fillPool() {
	//Debounce large amounts of requests to fill pool
	select {
	case <-time.After(1 * time.Millisecond):
		return
	case <-c.cFillingPool:
		defer func() { c.cFillingPool <- 1 }()
	}

	c.mu.Lock()
	isClosed := c.quit
	c.mu.Unlock()
	//Exit if cluster(session) is closed
	if isClosed {
		return
	}

	c.hostMu.RLock()

	//Walk through list of defined hosts
	var wg sync.WaitGroup
	for host := range c.hosts {
		addr := JoinHostPort(host, c.cfg.Port)

		numConns := 1
		//See if the host already has connections in the pool
		c.mu.Lock()
		conns, ok := c.connPool[addr]
		c.mu.Unlock()

		if ok {
			//if the host has enough connections just exit
			numConns = conns.Size()
			if numConns >= c.cfg.NumConns {
				continue
			}
		} else {
			//See if the host is reachable
			if err := c.connect(addr); err != nil {
				continue
			}
		}

		//This is reached if the host is responsive and needs more connections
		//Create connections for host synchronously to mitigate flooding the host.
		wg.Add(1)
		go func(a string, conns int) {
			defer wg.Done()
			for ; conns < c.cfg.NumConns; conns++ {
				c.connect(a)
			}
		}(addr, numConns)
	}

	c.hostMu.RUnlock()

	//Wait until we're finished connecting to each host before returning
	wg.Wait()
}

// Should only be called if c.mu is locked
func (c *SimplePool) removeConnLocked(conn *Conn) {
	conn.Close()
	connPool := c.connPool[conn.addr]
	if connPool == nil {
		return
	}
	connPool.RemoveNode(conn)
	if connPool.Size() == 0 {
		c.hostPool.RemoveNode(connPool)
		delete(c.connPool, conn.addr)
	}
	delete(c.conns, conn)
}

func (c *SimplePool) removeConn(conn *Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.removeConnLocked(conn)
}

//HandleError is called by a Connection object to report to the pool an error has occured.
//Logic is then executed within the pool to clean up the erroroneous connection and try to
//top off the pool.
func (c *SimplePool) HandleError(conn *Conn, err error, closed bool) {
	if !closed {
		// ignore all non-fatal errors
		return
	}
	c.removeConn(conn)
	c.mu.Lock()
	poolClosed := c.quit
	c.mu.Unlock()
	if !poolClosed {
		go c.fillPool() // top off pool.
	}
}

//Pick selects a connection to be used by the query.
func (c *SimplePool) Pick(qry *Query) *Conn {
	//Check if connections are available
	c.mu.Lock()
	conns := len(c.conns)
	c.mu.Unlock()

	if conns == 0 {
		//try to populate the pool before returning.
		c.fillPool()
	}

	return c.hostPool.Pick(qry)
}

//Size returns the number of connections currently active in the pool
func (p *SimplePool) Size() int {
	p.mu.Lock()
	conns := len(p.conns)
	p.mu.Unlock()
	return conns
}

//Close kills the pool and all associated connections.
func (c *SimplePool) Close() {
	c.quitOnce.Do(func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.quit = true
		close(c.quitWait)
		for conn := range c.conns {
			c.removeConnLocked(conn)
		}
	})
}

func (c *SimplePool) SetHosts(hosts []HostInfo) {

	c.hostMu.Lock()
	toRemove := make(map[string]struct{})
	for k := range c.hosts {
		toRemove[k] = struct{}{}
	}

	for _, host := range hosts {
		host := host
		delete(toRemove, host.Peer)
		// we already have it
		if _, ok := c.hosts[host.Peer]; ok {
			// TODO: Check rack, dc, token range is consistent, trigger topology change
			// update stored host
			continue
		}

		c.hosts[host.Peer] = &host
	}

	// can we hold c.mu whilst iterating this loop?
	for addr := range toRemove {
		c.removeHostLocked(addr)
	}
	c.hostMu.Unlock()

	c.fillPool()
}

func (c *SimplePool) removeHostLocked(addr string) {
	if _, ok := c.hosts[addr]; !ok {
		return
	}
	delete(c.hosts, addr)

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.connPool[addr]; !ok {
		return
	}

	for conn := range c.conns {
		if conn.Address() == addr {
			c.removeConnLocked(conn)
		}
	}
}

//NewRoundRobinConnPool creates a connection pool which selects hosts by
//round-robin, and then selects a connection for that host by round-robin.
func NewRoundRobinConnPool(cfg *ClusterConfig) (ConnectionPool, error) {
	return NewPolicyConnPool(
		cfg,
		NewRoundRobinHostPolicy(),
		NewRoundRobinConnPolicy,
	)
}

//NewTokenAwareConnPool creates a connection pool which selects hosts by
//a token aware policy, and then selects a connection for that host by
//round-robin.
func NewTokenAwareConnPool(cfg *ClusterConfig) (ConnectionPool, error) {
	return NewPolicyConnPool(
		cfg,
		NewTokenAwareHostPolicy(NewRoundRobinHostPolicy()),
		NewRoundRobinConnPolicy,
	)
}

type policyConnPool struct {
	port     int
	numConns int
	connCfg  ConnConfig
	keyspace string

	mu            sync.RWMutex
	hostPolicy    HostSelectionPolicy
	connPolicy    func() ConnSelectionPolicy
	hostConnPools map[string]*hostConnPool
}

//Creates a policy based connection pool. This func isn't meant to be directly
//used as a NewPoolFunc in ClusterConfig, instead a func should be created
//which satisfies the NewPoolFunc type, which calls this func with the desired
//hostPolicy and connPolicy; see NewRoundRobinConnPool or NewTokenAwareConnPool
//for examples.
func NewPolicyConnPool(
	cfg *ClusterConfig,
	hostPolicy HostSelectionPolicy,
	connPolicy func() ConnSelectionPolicy,
) (ConnectionPool, error) {
	var err error
	var tlsConfig *tls.Config

	if cfg.SslOpts != nil {
		tlsConfig, err = setupTLSConfig(cfg.SslOpts)
		if err != nil {
			return nil, err
		}
	}

	// create the pool
	pool := &policyConnPool{
		port:     cfg.Port,
		numConns: cfg.NumConns,
		connCfg: ConnConfig{
			ProtoVersion:  cfg.ProtoVersion,
			CQLVersion:    cfg.CQLVersion,
			Timeout:       cfg.Timeout,
			NumStreams:    cfg.NumStreams,
			Compressor:    cfg.Compressor,
			Authenticator: cfg.Authenticator,
			Keepalive:     cfg.SocketKeepalive,
			tlsConfig:     tlsConfig,
		},
		keyspace:      cfg.Keyspace,
		hostPolicy:    hostPolicy,
		connPolicy:    connPolicy,
		hostConnPools: map[string]*hostConnPool{},
	}

	hosts := make([]HostInfo, len(cfg.Hosts))
	for i, hostAddr := range cfg.Hosts {
		hosts[i].Peer = hostAddr
	}

	pool.SetHosts(hosts)

	return pool, nil
}

func (p *policyConnPool) SetHosts(hosts []HostInfo) {
	p.mu.Lock()

	toRemove := make(map[string]struct{})
	for addr := range p.hostConnPools {
		toRemove[addr] = struct{}{}
	}

	// TODO connect to hosts in parallel, but wait for pools to be
	// created before returning

	for i := range hosts {
		pool, exists := p.hostConnPools[hosts[i].Peer]
		if !exists {
			// create a connection pool for the host
			pool = newHostConnPool(
				hosts[i].Peer,
				p.port,
				p.numConns,
				p.connCfg,
				p.keyspace,
				p.connPolicy(),
			)
			p.hostConnPools[hosts[i].Peer] = pool
		} else {
			// still have this host, so don't remove it
			delete(toRemove, hosts[i].Peer)
		}
	}

	for addr := range toRemove {
		pool := p.hostConnPools[addr]
		delete(p.hostConnPools, addr)
		pool.Close()
	}

	// update the policy
	p.hostPolicy.SetHosts(hosts)

	p.mu.Unlock()
}

func (p *policyConnPool) SetPartitioner(partitioner string) {
	p.hostPolicy.SetPartitioner(partitioner)
}

func (p *policyConnPool) Size() int {
	p.mu.RLock()
	count := 0
	for _, pool := range p.hostConnPools {
		count += pool.Size()
	}
	p.mu.RUnlock()

	return count
}

func (p *policyConnPool) Pick(qry *Query) *Conn {
	nextHost := p.hostPolicy.Pick(qry)

	p.mu.RLock()
	var host *HostInfo
	var conn *Conn
	for conn == nil {
		host = nextHost()
		if host == nil {
			break
		}
		conn = p.hostConnPools[host.Peer].Pick(qry)
	}
	p.mu.RUnlock()
	return conn
}

func (p *policyConnPool) Close() {
	p.mu.Lock()

	// remove the hosts from the policy
	p.hostPolicy.SetHosts([]HostInfo{})

	// close the pools
	for addr, pool := range p.hostConnPools {
		delete(p.hostConnPools, addr)
		pool.Close()
	}
	p.mu.Unlock()
}

// hostConnPool is a connection pool for a single host.
// Connection selection is based on a provided ConnSelectionPolicy
type hostConnPool struct {
	host     string
	port     int
	addr     string
	size     int
	connCfg  ConnConfig
	keyspace string
	policy   ConnSelectionPolicy
	// protection for conns, closed, filling
	mu      sync.RWMutex
	conns   []*Conn
	closed  bool
	filling bool
}

func newHostConnPool(
	host string,
	port int,
	size int,
	connCfg ConnConfig,
	keyspace string,
	policy ConnSelectionPolicy,
) *hostConnPool {

	pool := &hostConnPool{
		host:     host,
		port:     port,
		addr:     JoinHostPort(host, port),
		size:     size,
		connCfg:  connCfg,
		keyspace: keyspace,
		policy:   policy,
		conns:    make([]*Conn, 0, size),
		filling:  false,
		closed:   false,
	}

	// fill the pool with the initial connections before returning
	pool.fill()

	return pool
}

// Pick a connection from this connection pool for the given query.
func (pool *hostConnPool) Pick(qry *Query) *Conn {
	pool.mu.RLock()
	if pool.closed {
		pool.mu.RUnlock()
		return nil
	}

	empty := len(pool.conns) == 0
	pool.mu.RUnlock()

	if empty {
		// try to fill the empty pool
		go pool.fill()
		return nil
	}

	return pool.policy.Pick(qry)
}

//Size returns the number of connections currently active in the pool
func (pool *hostConnPool) Size() int {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	return len(pool.conns)
}

//Close the connection pool
func (pool *hostConnPool) Close() {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if pool.closed {
		return
	}
	pool.closed = true

	// drain, but don't wait
	go pool.drain()
}

// Fill the connection pool
func (pool *hostConnPool) fill() {
	pool.mu.RLock()
	// avoid filling a closed pool, or concurrent filling
	if pool.closed || pool.filling {
		pool.mu.RUnlock()
		return
	}

	// determine the filling work to be done
	startCount := len(pool.conns)
	fillCount := pool.size - startCount

	// avoid filling a full (or overfull) pool
	if fillCount <= 0 {
		pool.mu.RUnlock()
		return
	}

	// switch from read to write lock
	pool.mu.RUnlock()
	pool.mu.Lock()

	// double check everything since the lock was released
	startCount = len(pool.conns)
	fillCount = pool.size - startCount
	if pool.closed || pool.filling || fillCount <= 0 {
		// looks like another goroutine already beat this
		// goroutine to the filling
		pool.mu.Unlock()
		return
	}

	// ok fill the pool
	pool.filling = true

	// allow others to access the pool while filling
	pool.mu.Unlock()
	// only this goroutine should make calls to fill/empty the pool at this
	// point until after this routine or its subordinates calls
	// fillingStopped

	// fill only the first connection synchronously
	if startCount == 0 {
		err := pool.connect()
		pool.logConnectErr(err)

		if err != nil {
			// probably unreachable host
			go pool.fillingStopped()
			return
		}

		// filled one
		fillCount--

		// connect all connections to this host in sync
		for fillCount > 0 {
			err := pool.connect()
			pool.logConnectErr(err)

			// decrement, even on error
			fillCount--
		}

		go pool.fillingStopped()
		return
	}

	// fill the rest of the pool asynchronously
	go func() {
		for fillCount > 0 {
			err := pool.connect()
			pool.logConnectErr(err)

			// decrement, even on error
			fillCount--
		}

		// mark the end of filling
		pool.fillingStopped()
	}()
}

func (pool *hostConnPool) logConnectErr(err error) {
	if opErr, ok := err.(*net.OpError); ok && (opErr.Op == "dial" || opErr.Op == "read") {
		// connection refused
		// these are typical during a node outage so avoid log spam.
	} else if err != nil {
		// unexpected error
		log.Printf("error: failed to connect to %s due to error: %v", pool.addr, err)
	}
}

// transition back to a not-filling state.
func (pool *hostConnPool) fillingStopped() {
	// wait for some time to avoid back-to-back filling
	// this provides some time between failed attempts
	// to fill the pool for the host to recover
	time.Sleep(time.Duration(rand.Int31n(100)+31) * time.Millisecond)

	pool.mu.Lock()
	pool.filling = false
	pool.mu.Unlock()
}

// create a new connection to the host and add it to the pool
func (pool *hostConnPool) connect() error {
	// try to connect
	conn, err := Connect(pool.addr, pool.connCfg, pool)
	if err != nil {
		return err
	}

	if pool.keyspace != "" {
		// set the keyspace
		if err := conn.UseKeyspace(pool.keyspace); err != nil {
			conn.Close()
			return err
		}
	}

	// add the Conn to the pool
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if pool.closed {
		conn.Close()
		return nil
	}

	pool.conns = append(pool.conns, conn)
	pool.policy.SetConns(pool.conns)
	return nil
}

// handle any error from a Conn
func (pool *hostConnPool) HandleError(conn *Conn, err error, closed bool) {
	if !closed {
		// still an open connection, so continue using it
		return
	}

	pool.mu.Lock()
	defer pool.mu.Unlock()

	if pool.closed {
		// pool closed
		return
	}

	// find the connection index
	for i, candidate := range pool.conns {
		if candidate == conn {
			// remove the connection, not preserving order
			pool.conns[i], pool.conns = pool.conns[len(pool.conns)-1], pool.conns[:len(pool.conns)-1]

			// update the policy
			pool.policy.SetConns(pool.conns)

			// lost a connection, so fill the pool
			go pool.fill()
			break
		}
	}
}

// removes and closes all connections from the pool
func (pool *hostConnPool) drain() {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// empty the pool
	conns := pool.conns
	pool.conns = pool.conns[:0]

	// update the policy
	pool.policy.SetConns(pool.conns)

	// close the connections
	for _, conn := range conns {
		conn.Close()
	}
}

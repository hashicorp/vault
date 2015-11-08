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

// interface to implement to receive the host information
type SetHosts interface {
	SetHosts(hosts []HostInfo)
}

// interface to implement to receive the partitioner value
type SetPartitioner interface {
	SetPartitioner(partitioner string)
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

type policyConnPool struct {
	session *Session

	port     int
	numConns int
	connCfg  *ConnConfig
	keyspace string

	mu            sync.RWMutex
	hostPolicy    HostSelectionPolicy
	connPolicy    func() ConnSelectionPolicy
	hostConnPools map[string]*hostConnPool
}

func newPolicyConnPool(session *Session, hostPolicy HostSelectionPolicy,
	connPolicy func() ConnSelectionPolicy) (*policyConnPool, error) {

	var (
		err       error
		tlsConfig *tls.Config
	)

	cfg := session.cfg

	if cfg.SslOpts != nil {
		tlsConfig, err = setupTLSConfig(cfg.SslOpts)
		if err != nil {
			return nil, err
		}
	}

	// create the pool
	pool := &policyConnPool{
		session:  session,
		port:     cfg.Port,
		numConns: cfg.NumConns,
		connCfg: &ConnConfig{
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
	defer p.mu.Unlock()

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
				p.session,
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

func (p *policyConnPool) Pick(qry *Query) (SelectedHost, *Conn) {
	nextHost := p.hostPolicy.Pick(qry)

	var (
		host SelectedHost
		conn *Conn
	)

	p.mu.RLock()
	defer p.mu.RUnlock()
	for conn == nil {
		host = nextHost()
		if host == nil {
			break
		} else if host.Info() == nil {
			panic(fmt.Sprintf("policy %T returned no host info: %+v", p.hostPolicy, host))
		}

		pool, ok := p.hostConnPools[host.Info().Peer]
		if !ok {
			continue
		}

		conn = pool.Pick(qry)
	}
	return host, conn
}

func (p *policyConnPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	// remove the hosts from the policy
	p.hostPolicy.SetHosts([]HostInfo{})

	// close the pools
	for addr, pool := range p.hostConnPools {
		delete(p.hostConnPools, addr)
		pool.Close()
	}
}

// hostConnPool is a connection pool for a single host.
// Connection selection is based on a provided ConnSelectionPolicy
type hostConnPool struct {
	session  *Session
	host     string
	port     int
	addr     string
	size     int
	connCfg  *ConnConfig
	keyspace string
	policy   ConnSelectionPolicy
	// protection for conns, closed, filling
	mu      sync.RWMutex
	conns   []*Conn
	closed  bool
	filling bool
}

func newHostConnPool(session *Session, host string, port, size int, connCfg *ConnConfig,
	keyspace string, policy ConnSelectionPolicy) *hostConnPool {

	pool := &hostConnPool{
		session:  session,
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
	conn, err := Connect(pool.addr, pool.connCfg, pool, pool.session)
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

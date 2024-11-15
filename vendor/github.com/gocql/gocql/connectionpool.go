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
 * Copyright (c) 2012, The Gocql authors,
 * provided under the BSD-3-Clause License.
 * See the NOTICE file distributed with this work for additional information.
 */

package gocql

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// interface to implement to receive the host information
type SetHosts interface {
	SetHosts(hosts []*HostInfo)
}

// interface to implement to receive the partitioner value
type SetPartitioner interface {
	SetPartitioner(partitioner string)
}

func setupTLSConfig(sslOpts *SslOptions) (*tls.Config, error) {
	//  Config.InsecureSkipVerify | EnableHostVerification | Result
	//  Config is nil             | true                   | verify host
	//  Config is nil             | false                  | do not verify host
	//  false                     | false                  | verify host
	//  true                      | false                  | do not verify host
	//  false                     | true                   | verify host
	//  true                      | true                   | verify host
	var tlsConfig *tls.Config
	if sslOpts.Config == nil {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: !sslOpts.EnableHostVerification,
		}
	} else {
		// use clone to avoid race.
		tlsConfig = sslOpts.Config.Clone()
	}

	if tlsConfig.InsecureSkipVerify && sslOpts.EnableHostVerification {
		tlsConfig.InsecureSkipVerify = false
	}

	// ca cert is optional
	if sslOpts.CaPath != "" {
		if tlsConfig.RootCAs == nil {
			tlsConfig.RootCAs = x509.NewCertPool()
		}

		pem, err := ioutil.ReadFile(sslOpts.CaPath)
		if err != nil {
			return nil, fmt.Errorf("connectionpool: unable to open CA certs: %v", err)
		}

		if !tlsConfig.RootCAs.AppendCertsFromPEM(pem) {
			return nil, errors.New("connectionpool: failed parsing or CA certs")
		}
	}

	if sslOpts.CertPath != "" || sslOpts.KeyPath != "" {
		mycert, err := tls.LoadX509KeyPair(sslOpts.CertPath, sslOpts.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("connectionpool: unable to load X509 key pair: %v", err)
		}
		tlsConfig.Certificates = append(tlsConfig.Certificates, mycert)
	}

	return tlsConfig, nil
}

type policyConnPool struct {
	session *Session

	port     int
	numConns int
	keyspace string

	mu            sync.RWMutex
	hostConnPools map[string]*hostConnPool
}

func connConfig(cfg *ClusterConfig) (*ConnConfig, error) {
	var (
		err        error
		hostDialer HostDialer
	)

	hostDialer = cfg.HostDialer
	if hostDialer == nil {
		var tlsConfig *tls.Config

		// TODO(zariel): move tls config setup into session init.
		if cfg.SslOpts != nil {
			tlsConfig, err = setupTLSConfig(cfg.SslOpts)
			if err != nil {
				return nil, err
			}
		}

		dialer := cfg.Dialer
		if dialer == nil {
			d := &net.Dialer{
				Timeout: cfg.ConnectTimeout,
			}
			if cfg.SocketKeepalive > 0 {
				d.KeepAlive = cfg.SocketKeepalive
			}
			dialer = d
		}

		hostDialer = &defaultHostDialer{
			dialer:    dialer,
			tlsConfig: tlsConfig,
		}
	}

	return &ConnConfig{
		ProtoVersion:   cfg.ProtoVersion,
		CQLVersion:     cfg.CQLVersion,
		Timeout:        cfg.Timeout,
		WriteTimeout:   cfg.WriteTimeout,
		ConnectTimeout: cfg.ConnectTimeout,
		Dialer:         cfg.Dialer,
		HostDialer:     hostDialer,
		Compressor:     cfg.Compressor,
		Authenticator:  cfg.Authenticator,
		AuthProvider:   cfg.AuthProvider,
		Keepalive:      cfg.SocketKeepalive,
		Logger:         cfg.logger(),
	}, nil
}

func newPolicyConnPool(session *Session) *policyConnPool {
	// create the pool
	pool := &policyConnPool{
		session:       session,
		port:          session.cfg.Port,
		numConns:      session.cfg.NumConns,
		keyspace:      session.cfg.Keyspace,
		hostConnPools: map[string]*hostConnPool{},
	}

	return pool
}

func (p *policyConnPool) SetHosts(hosts []*HostInfo) {
	p.mu.Lock()
	defer p.mu.Unlock()

	toRemove := make(map[string]struct{})
	for hostID := range p.hostConnPools {
		toRemove[hostID] = struct{}{}
	}

	pools := make(chan *hostConnPool)
	createCount := 0
	for _, host := range hosts {
		if !host.IsUp() {
			// don't create a connection pool for a down host
			continue
		}
		hostID := host.HostID()
		if _, exists := p.hostConnPools[hostID]; exists {
			// still have this host, so don't remove it
			delete(toRemove, hostID)
			continue
		}

		createCount++
		go func(host *HostInfo) {
			// create a connection pool for the host
			pools <- newHostConnPool(
				p.session,
				host,
				p.port,
				p.numConns,
				p.keyspace,
			)
		}(host)
	}

	// add created pools
	for createCount > 0 {
		pool := <-pools
		createCount--
		if pool.Size() > 0 {
			// add pool only if there a connections available
			p.hostConnPools[pool.host.HostID()] = pool
		}
	}

	for addr := range toRemove {
		pool := p.hostConnPools[addr]
		delete(p.hostConnPools, addr)
		go pool.Close()
	}
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

func (p *policyConnPool) getPool(host *HostInfo) (pool *hostConnPool, ok bool) {
	hostID := host.HostID()
	p.mu.RLock()
	pool, ok = p.hostConnPools[hostID]
	p.mu.RUnlock()
	return
}

func (p *policyConnPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	// close the pools
	for addr, pool := range p.hostConnPools {
		delete(p.hostConnPools, addr)
		pool.Close()
	}
}

func (p *policyConnPool) addHost(host *HostInfo) {
	hostID := host.HostID()
	p.mu.Lock()
	pool, ok := p.hostConnPools[hostID]
	if !ok {
		pool = newHostConnPool(
			p.session,
			host,
			host.Port(), // TODO: if port == 0 use pool.port?
			p.numConns,
			p.keyspace,
		)

		p.hostConnPools[hostID] = pool
	}
	p.mu.Unlock()

	pool.fill()
}

func (p *policyConnPool) removeHost(hostID string) {
	p.mu.Lock()
	pool, ok := p.hostConnPools[hostID]
	if !ok {
		p.mu.Unlock()
		return
	}

	delete(p.hostConnPools, hostID)
	p.mu.Unlock()

	go pool.Close()
}

// hostConnPool is a connection pool for a single host.
// Connection selection is based on a provided ConnSelectionPolicy
type hostConnPool struct {
	session  *Session
	host     *HostInfo
	port     int
	size     int
	keyspace string
	// protection for conns, closed, filling
	mu      sync.RWMutex
	conns   []*Conn
	closed  bool
	filling bool

	pos    uint32
	logger StdLogger
}

func (h *hostConnPool) String() string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return fmt.Sprintf("[filling=%v closed=%v conns=%v size=%v host=%v]",
		h.filling, h.closed, len(h.conns), h.size, h.host)
}

func newHostConnPool(session *Session, host *HostInfo, port, size int,
	keyspace string) *hostConnPool {

	pool := &hostConnPool{
		session:  session,
		host:     host,
		port:     port,
		size:     size,
		keyspace: keyspace,
		conns:    make([]*Conn, 0, size),
		filling:  false,
		closed:   false,
		logger:   session.logger,
	}

	// the pool is not filled or connected
	return pool
}

// Pick a connection from this connection pool for the given query.
func (pool *hostConnPool) Pick() *Conn {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	if pool.closed {
		return nil
	}

	size := len(pool.conns)
	if size < pool.size {
		// try to fill the pool
		go pool.fill()

		if size == 0 {
			return nil
		}
	}

	pos := int(atomic.AddUint32(&pool.pos, 1) - 1)

	var (
		leastBusyConn    *Conn
		streamsAvailable int
	)

	// find the conn which has the most available streams, this is racy
	for i := 0; i < size; i++ {
		conn := pool.conns[(pos+i)%size]
		if streams := conn.AvailableStreams(); streams > streamsAvailable {
			leastBusyConn = conn
			streamsAvailable = streams
		}
	}

	return leastBusyConn
}

// Size returns the number of connections currently active in the pool
func (pool *hostConnPool) Size() int {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	return len(pool.conns)
}

// Close the connection pool
func (pool *hostConnPool) Close() {
	pool.mu.Lock()

	if pool.closed {
		pool.mu.Unlock()
		return
	}
	pool.closed = true

	// ensure we dont try to reacquire the lock in handleError
	// TODO: improve this as the following can happen
	// 1) we have locked pool.mu write lock
	// 2) conn.Close calls conn.closeWithError(nil)
	// 3) conn.closeWithError calls conn.Close() which returns an error
	// 4) conn.closeWithError calls pool.HandleError with the error from conn.Close
	// 5) pool.HandleError tries to lock pool.mu
	// deadlock

	// empty the pool
	conns := pool.conns
	pool.conns = nil

	pool.mu.Unlock()

	// close the connections
	for _, conn := range conns {
		conn.Close()
	}
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
			pool.fillingStopped(err)
			return
		}
		// notify the session that this node is connected
		go pool.session.handleNodeConnected(pool.host)

		// filled one
		fillCount--
	}

	// fill the rest of the pool asynchronously
	go func() {
		err := pool.connectMany(fillCount)

		// mark the end of filling
		pool.fillingStopped(err)

		if err == nil && startCount > 0 {
			// notify the session that this node is connected again
			go pool.session.handleNodeConnected(pool.host)
		}
	}()
}

func (pool *hostConnPool) logConnectErr(err error) {
	if opErr, ok := err.(*net.OpError); ok && (opErr.Op == "dial" || opErr.Op == "read") {
		// connection refused
		// these are typical during a node outage so avoid log spam.
		if gocqlDebug {
			pool.logger.Printf("gocql: unable to dial %q: %v\n", pool.host, err)
		}
	} else if err != nil {
		// unexpected error
		pool.logger.Printf("error: failed to connect to %q due to error: %v", pool.host, err)
	}
}

// transition back to a not-filling state.
func (pool *hostConnPool) fillingStopped(err error) {
	if err != nil {
		if gocqlDebug {
			pool.logger.Printf("gocql: filling stopped %q: %v\n", pool.host.ConnectAddress(), err)
		}
		// wait for some time to avoid back-to-back filling
		// this provides some time between failed attempts
		// to fill the pool for the host to recover
		time.Sleep(time.Duration(rand.Int31n(100)+31) * time.Millisecond)
	}

	pool.mu.Lock()
	pool.filling = false
	count := len(pool.conns)
	host := pool.host
	port := pool.port
	pool.mu.Unlock()

	// if we errored and the size is now zero, make sure the host is marked as down
	// see https://github.com/apache/cassandra-gocql-driver/issues/1614
	if gocqlDebug {
		pool.logger.Printf("gocql: conns of pool after stopped %q: %v\n", host.ConnectAddress(), count)
	}
	if err != nil && count == 0 {
		if pool.session.cfg.ConvictionPolicy.AddFailure(err, host) {
			pool.session.handleNodeDown(host.ConnectAddress(), port)
		}
	}
}

// connectMany creates new connections concurrent.
func (pool *hostConnPool) connectMany(count int) error {
	if count == 0 {
		return nil
	}
	var (
		wg         sync.WaitGroup
		mu         sync.Mutex
		connectErr error
	)
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			err := pool.connect()
			pool.logConnectErr(err)
			if err != nil {
				mu.Lock()
				connectErr = err
				mu.Unlock()
			}
		}()
	}
	// wait for all connections are done
	wg.Wait()

	return connectErr
}

// create a new connection to the host and add it to the pool
func (pool *hostConnPool) connect() (err error) {
	// TODO: provide a more robust connection retry mechanism, we should also
	// be able to detect hosts that come up by trying to connect to downed ones.
	// try to connect
	var conn *Conn
	reconnectionPolicy := pool.session.cfg.ReconnectionPolicy
	for i := 0; i < reconnectionPolicy.GetMaxRetries(); i++ {
		conn, err = pool.session.connect(pool.session.ctx, pool.host, pool)
		if err == nil {
			break
		}
		if opErr, isOpErr := err.(*net.OpError); isOpErr {
			// if the error is not a temporary error (ex: network unreachable) don't
			//  retry
			if !opErr.Temporary() {
				break
			}
		}
		if gocqlDebug {
			pool.logger.Printf("gocql: connection failed %q: %v, reconnecting with %T\n",
				pool.host.ConnectAddress(), err, reconnectionPolicy)
		}
		time.Sleep(reconnectionPolicy.GetInterval(i))
	}

	if err != nil {
		return err
	}

	if pool.keyspace != "" {
		// set the keyspace
		if err = conn.UseKeyspace(pool.keyspace); err != nil {
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

	return nil
}

// handle any error from a Conn
func (pool *hostConnPool) HandleError(conn *Conn, err error, closed bool) {
	if !closed {
		// still an open connection, so continue using it
		return
	}

	// TODO: track the number of errors per host and detect when a host is dead,
	// then also have something which can detect when a host comes back.
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if pool.closed {
		// pool closed
		return
	}

	if gocqlDebug {
		pool.logger.Printf("gocql: pool connection error %q: %v\n", conn.addr, err)
	}

	// find the connection index
	for i, candidate := range pool.conns {
		if candidate == conn {
			// remove the connection, not preserving order
			pool.conns[i], pool.conns = pool.conns[len(pool.conns)-1], pool.conns[:len(pool.conns)-1]

			// lost a connection, so fill the pool
			go pool.fill()
			break
		}
	}
}

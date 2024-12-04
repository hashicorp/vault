// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cluster

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/base62"
	"go.uber.org/atomic"
)

// InmemLayer is an in-memory implementation of NetworkLayer. This is
// primarially useful for tests.
type InmemLayer struct {
	listener *inmemListener
	addr     string
	logger   log.Logger

	servConns   map[string][]net.Conn
	clientConns map[string][]net.Conn

	peers map[string]*InmemLayer
	l     sync.Mutex

	stopped *atomic.Bool
	stopCh  chan struct{}

	connectionCh chan *ConnectionInfo
	readerDelay  time.Duration
	forceTimeout string
}

// NewInmemLayer returns a new in-memory layer configured to listen on the
// provided address.
func NewInmemLayer(addr string, logger log.Logger) *InmemLayer {
	return &InmemLayer{
		addr:        addr,
		logger:      logger,
		stopped:     atomic.NewBool(false),
		stopCh:      make(chan struct{}),
		peers:       make(map[string]*InmemLayer),
		servConns:   make(map[string][]net.Conn),
		clientConns: make(map[string][]net.Conn),
	}
}

func (l *InmemLayer) SetConnectionCh(ch chan *ConnectionInfo) {
	l.l.Lock()
	l.connectionCh = ch
	l.l.Unlock()
}

func (l *InmemLayer) SetReaderDelay(delay time.Duration) {
	l.l.Lock()
	defer l.l.Unlock()

	l.readerDelay = delay

	// Update the existing server and client connections
	for _, servConns := range l.servConns {
		for _, c := range servConns {
			c.(*delayedConn).SetDelay(delay)
		}
	}

	for _, clientConns := range l.clientConns {
		for _, c := range clientConns {
			c.(*delayedConn).SetDelay(delay)
		}
	}
}

func (l *InmemLayer) SetForceTimeout(addr string) {
	l.l.Lock()
	defer l.l.Unlock()

	l.forceTimeout = addr
}

// Addrs implements NetworkLayer.
func (l *InmemLayer) Addrs() []net.Addr {
	l.l.Lock()
	defer l.l.Unlock()

	if l.listener == nil {
		return nil
	}

	return []net.Addr{l.listener.Addr()}
}

// Listeners implements NetworkLayer.
func (l *InmemLayer) Listeners() []NetworkListener {
	l.l.Lock()
	defer l.l.Unlock()

	if l.listener != nil {
		return []NetworkListener{l.listener}
	}

	l.listener = &inmemListener{
		addr:         l.addr,
		pendingConns: make(chan net.Conn, 1),

		stopped: atomic.NewBool(false),
		stopCh:  make(chan struct{}),
	}

	return []NetworkListener{l.listener}
}

// Partition forces the inmem layer to disconnect itself from peers and prevents
// creating new connections. The returned function will add all peers back
// and re-enable connections
func (l *InmemLayer) Partition() (unpartition func()) {
	l.l.Lock()
	peersCopy := make([]*InmemLayer, 0, len(l.peers))
	for _, peer := range l.peers {
		peersCopy = append(peersCopy, peer)
	}
	l.l.Unlock()
	l.DisconnectAll()
	return func() {
		for _, peer := range peersCopy {
			l.Connect(peer)
		}
	}
}

// Dial implements NetworkLayer.
func (l *InmemLayer) Dial(addr string, timeout time.Duration, tlsConfig *tls.Config) (*tls.Conn, error) {
	l.l.Lock()
	connectionCh := l.connectionCh

	if addr == l.addr {
		panic(fmt.Sprintf("%q attempted to dial itself", l.addr))
	}

	// This simulates an i/o timeout by sleeping for 20 seconds and returning
	// an error when the forceTimeout name is the same as the host we are
	// currently connecting to. Useful for checking how gRPC connections react
	// with timeouts.
	if l.forceTimeout == addr {
		l.logger.Debug("forcing timeout", "addr", addr, "me", l.addr)

		// gRPC sets a deadline of 20 seconds on the dial attempt, so
		// matching that here.
		time.Sleep(time.Second * 20)
		l.l.Unlock()
		return nil, deadlineError("i/o timeout")
	}

	peer, ok := l.peers[addr]
	l.l.Unlock()
	if !ok {
		return nil, errors.New("inmemlayer: no address found")
	}

	if timeout < 0 {
		return nil, fmt.Errorf("inmemlayer: timeout given is less than 0: %d", timeout)
	}

	alpn := ""
	if tlsConfig != nil {
		alpn = tlsConfig.NextProtos[0]
	}

	if l.logger.IsDebug() {
		l.logger.Debug("dialing connection", "node", l.addr, "remote", addr, "alpn", alpn)
	}

	if connectionCh != nil {
		select {
		case connectionCh <- &ConnectionInfo{
			Node:     l.addr,
			Remote:   addr,
			IsServer: false,
			ALPN:     alpn,
		}:
		case <-time.After(2 * time.Second):
			l.logger.Warn("failed to send connection info")
		}
	}

	conn, err := peer.clientConn(l.addr)
	if err != nil {
		return nil, err
	}

	tlsConn := tls.Client(conn, tlsConfig)

	l.l.Lock()
	l.clientConns[addr] = append(l.clientConns[addr], conn)
	l.l.Unlock()

	return tlsConn, nil
}

// clientConn is executed on a server when a new client connection comes in and
// needs to be Accepted.
func (l *InmemLayer) clientConn(addr string) (net.Conn, error) {
	l.l.Lock()

	if l.listener == nil {
		l.l.Unlock()
		return nil, errors.New("inmemlayer: listener not started")
	}

	_, ok := l.peers[addr]
	if !ok {
		l.l.Unlock()
		return nil, errors.New("inmemlayer: no peer found")
	}

	retConn, servConn := net.Pipe()

	retConn = newDelayedConn(retConn, l.readerDelay)
	servConn = newDelayedConn(servConn, l.readerDelay)

	l.servConns[addr] = append(l.servConns[addr], servConn)
	connectionCh := l.connectionCh
	pendingConns := l.listener.pendingConns
	l.l.Unlock()

	if l.logger.IsDebug() {
		l.logger.Debug("received connection", "node", l.addr, "remote", addr)
	}
	if connectionCh != nil {
		select {
		case connectionCh <- &ConnectionInfo{
			Node:     l.addr,
			Remote:   addr,
			IsServer: true,
		}:
		case <-time.After(2 * time.Second):
			l.logger.Warn("failed to send connection info")
		}
	}

	select {
	case pendingConns <- servConn:
	case <-time.After(5 * time.Second):
		return nil, errors.New("inmemlayer: timeout while accepting connection")
	}

	return retConn, nil
}

// Connect is used to connect this transport to another transport for
// a given peer name. This allows for local routing.
func (l *InmemLayer) Connect(remote *InmemLayer) {
	l.l.Lock()
	defer l.l.Unlock()
	l.peers[remote.addr] = remote
}

// Disconnect is used to remove the ability to route to a given peer.
func (l *InmemLayer) Disconnect(peer string) {
	l.l.Lock()
	defer l.l.Unlock()
	delete(l.peers, peer)

	// Remove any open connections
	servConns := l.servConns[peer]
	for _, c := range servConns {
		c.Close()
	}
	delete(l.servConns, peer)

	clientConns := l.clientConns[peer]
	for _, c := range clientConns {
		c.Close()
	}
	delete(l.clientConns, peer)
}

// DisconnectAll is used to remove all routes to peers.
func (l *InmemLayer) DisconnectAll() {
	l.l.Lock()
	defer l.l.Unlock()
	l.peers = make(map[string]*InmemLayer)

	// Close all connections
	for _, peerConns := range l.servConns {
		for _, c := range peerConns {
			c.Close()
		}
	}
	l.servConns = make(map[string][]net.Conn)

	for _, peerConns := range l.clientConns {
		for _, c := range peerConns {
			c.Close()
		}
	}
	l.clientConns = make(map[string][]net.Conn)
}

// Close is used to permanently disable the transport
func (l *InmemLayer) Close() error {
	if l.stopped.Swap(true) {
		return nil
	}

	l.DisconnectAll()
	close(l.stopCh)
	return nil
}

// inmemListener implements the NetworkListener interface.
type inmemListener struct {
	addr         string
	pendingConns chan net.Conn

	stopped *atomic.Bool
	stopCh  chan struct{}

	deadline time.Time
}

// Accept implements the NetworkListener interface.
func (ln *inmemListener) Accept() (net.Conn, error) {
	deadline := ln.deadline
	if !deadline.IsZero() {
		select {
		case conn := <-ln.pendingConns:
			return conn, nil
		case <-time.After(time.Until(deadline)):
			return nil, deadlineError("deadline")
		case <-ln.stopCh:
			return nil, errors.New("listener shut down")
		}
	}

	select {
	case conn := <-ln.pendingConns:
		return conn, nil
	case <-ln.stopCh:
		return nil, errors.New("listener shut down")
	}
}

// Close implements the NetworkListener interface.
func (ln *inmemListener) Close() error {
	if ln.stopped.Swap(true) {
		return nil
	}

	close(ln.stopCh)
	return nil
}

// Addr implements the NetworkListener interface.
func (ln *inmemListener) Addr() net.Addr {
	return inmemAddr{addr: ln.addr}
}

// SetDeadline implements the NetworkListener interface.
func (ln *inmemListener) SetDeadline(d time.Time) error {
	ln.deadline = d
	return nil
}

type inmemAddr struct {
	addr string
}

func (a inmemAddr) Network() string {
	return "inmem"
}

func (a inmemAddr) String() string {
	return a.addr
}

type deadlineError string

func (d deadlineError) Error() string   { return string(d) }
func (d deadlineError) Timeout() bool   { return true }
func (d deadlineError) Temporary() bool { return true }

// InmemLayerCluster composes a set of layers and handles connecting them all
// together. It also satisfies the NetworkLayerSet interface.
type InmemLayerCluster struct {
	layers []*InmemLayer
}

// NewInmemLayerCluster returns a new in-memory layer set that builds n nodes
// and connects them all together.
func NewInmemLayerCluster(clusterName string, nodes int, logger log.Logger) (*InmemLayerCluster, error) {
	if clusterName == "" {
		clusterID, err := base62.Random(4)
		if err != nil {
			return nil, err
		}
		clusterName = "cluster_" + clusterID
	}

	layers := make([]*InmemLayer, nodes)
	for i := 0; i < nodes; i++ {
		layers[i] = NewInmemLayer(fmt.Sprintf("%s_node_%d", clusterName, i), logger)
	}

	// Connect all the peers together
	for _, node := range layers {
		for _, peer := range layers {
			// Don't connect to itself
			if node == peer {
				continue
			}

			node.Connect(peer)
			peer.Connect(node)
		}
	}

	return &InmemLayerCluster{layers: layers}, nil
}

// ConnectCluster connects this cluster with the provided remote cluster,
// connecting all nodes to each other.
func (ic *InmemLayerCluster) ConnectCluster(remote *InmemLayerCluster) {
	for _, node := range ic.layers {
		for _, peer := range remote.layers {
			node.Connect(peer)
			peer.Connect(node)
		}
	}
}

// Layers implements the NetworkLayerSet interface.
func (ic *InmemLayerCluster) Layers() []NetworkLayer {
	ret := make([]NetworkLayer, len(ic.layers))
	for i, l := range ic.layers {
		ret[i] = l
	}

	return ret
}

func (ic *InmemLayerCluster) SetConnectionCh(ch chan *ConnectionInfo) {
	for _, node := range ic.layers {
		node.SetConnectionCh(ch)
	}
}

func (ic *InmemLayerCluster) SetReaderDelay(delay time.Duration) {
	for _, node := range ic.layers {
		node.SetReaderDelay(delay)
	}
}

func (ic *InmemLayerCluster) SetForceTimeout(addr string) {
	for _, node := range ic.layers {
		node.SetForceTimeout(addr)
	}
}

type ConnectionInfo struct {
	Node     string
	Remote   string
	IsServer bool
	ALPN     string
}

package cluster

import (
	"crypto/tls"
	"errors"
	"net"
	"sync"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/base62"
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
		pendingConns: make(chan net.Conn),

		stopped: atomic.NewBool(false),
		stopCh:  make(chan struct{}),
	}

	return []NetworkListener{l.listener}
}

// Dial implements NetworkLayer.
func (l *InmemLayer) Dial(addr string, timeout time.Duration, tlsConfig *tls.Config) (*tls.Conn, error) {
	l.l.Lock()
	defer l.l.Unlock()

	peer, ok := l.peers[addr]
	if !ok {
		return nil, errors.New("inmemlayer: no address found")
	}

	conn, err := peer.clientConn(l.addr)
	if err != nil {
		return nil, err
	}

	tlsConn := tls.Client(conn, tlsConfig)

	l.clientConns[addr] = append(l.clientConns[addr], tlsConn)

	return tlsConn, nil
}

// clientConn is executed on a server when a new client connection comes in and
// needs to be Accepted.
func (l *InmemLayer) clientConn(addr string) (net.Conn, error) {
	l.l.Lock()
	defer l.l.Unlock()

	if l.listener == nil {
		return nil, errors.New("inmemlayer: listener not started")
	}

	_, ok := l.peers[addr]
	if !ok {
		return nil, errors.New("inmemlayer: no peer found")
	}

	retConn, servConn := net.Pipe()

	l.servConns[addr] = append(l.servConns[addr], servConn)

	select {
	case l.listener.pendingConns <- servConn:
	case <-time.After(2 * time.Second):
		return nil, errors.New("inmemlayer: timeout while accepting connection")
	}

	return retConn, nil
}

// Connect is used to connect this transport to another transport for
// a given peer name. This allows for local routing.
func (l *InmemLayer) Connect(peer string, remote *InmemLayer) {
	l.l.Lock()
	defer l.l.Unlock()
	l.peers[peer] = remote
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
func NewInmemLayerCluster(nodes int, logger log.Logger) (*InmemLayerCluster, error) {
	clusterID, err := base62.Random(4)
	if err != nil {
		return nil, err
	}

	clusterName := "cluster_" + clusterID

	var layers []*InmemLayer
	for i := 0; i < nodes; i++ {
		nodeID, err := base62.Random(4)
		if err != nil {
			return nil, err
		}

		nodeName := clusterName + "_node_" + nodeID

		layers = append(layers, NewInmemLayer(nodeName, logger))
	}

	// Connect all the peers together
	for _, node := range layers {
		for _, peer := range layers {
			// Don't connect to itself
			if node == peer {
				continue
			}

			node.Connect(peer.addr, peer)
			peer.Connect(node.addr, node)
		}
	}

	return &InmemLayerCluster{layers: layers}, nil
}

// ConnectCluster connects this cluster with the provided remote cluster,
// connecting all nodes to each other.
func (ic *InmemLayerCluster) ConnectCluster(remote *InmemLayerCluster) {
	for _, node := range ic.layers {
		for _, peer := range remote.layers {
			node.Connect(peer.addr, peer)
			peer.Connect(node.addr, node)
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

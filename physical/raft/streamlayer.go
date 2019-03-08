package raft

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	fmt "fmt"
	"net"
	"sync"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/vault/helper/consts"
)

// RaftLayer implements the raft.StreamLayer interface,
// so that we can use a single RPC layer for Raft and Nomad
type raftLayer struct {
	// Addr is the listener address to return
	addr net.Addr

	// connCh is used to accept connections
	connCh chan net.Conn

	// Tracks if we are closed
	closed    bool
	closeCh   chan struct{}
	closeLock sync.Mutex

	logger log.Logger

	dialerFunc func(string, time.Duration) (net.Conn, error)
}

func NewRaftLayer(logger log.Logger, addr net.Addr) *raftLayer {
	layer := &raftLayer{
		addr:    addr,
		connCh:  make(chan net.Conn),
		closeCh: make(chan struct{}),
		logger:  logger,
	}
	return layer
}

func (l *raftLayer) SetAddr(addr net.Addr) {
	l.addr = addr
}

func (l *raftLayer) ClientLookup(context.Context, *tls.CertificateRequestInfo) (*tls.Certificate, error) {
	return nil, nil
}
func (l *raftLayer) ServerLookup(context.Context, *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return nil, nil
}
func (l *raftLayer) CALookup(context.Context) (*x509.Certificate, error) {
	return nil, nil
}

func (l *raftLayer) Stop() error {
	l.Close()
	return nil
}

// Handoff is used to hand off a connection to the
// RaftLayer. This allows it to be Accept()'ed
func (l *raftLayer) Handoff(ctx context.Context, wg *sync.WaitGroup, quit chan struct{}, conn *tls.Conn) error {
	if l.closed {
		return errors.New("raft is shutdown")
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case l.connCh <- conn:
		case <-l.closeCh:
		case <-ctx.Done():
		case <-quit:
		}
	}()

	return nil
}

// Accept is used to return connection which are
// dialed to be used with the Raft layer
func (l *raftLayer) Accept() (net.Conn, error) {
	select {
	case conn := <-l.connCh:
		return conn, nil
	case <-l.closeCh:
		return nil, fmt.Errorf("Raft RPC layer closed")
	}
}

// Close is used to stop listening for Raft connections
func (l *raftLayer) Close() error {
	l.closeLock.Lock()
	defer l.closeLock.Unlock()

	if !l.closed {
		l.closed = true
		close(l.closeCh)
	}
	return nil
}

/*
// getTLSWrapper is used to retrieve the current TLS wrapper
func (l *RaftLayer) getTLSWrapper() tlsutil.Wrapper {
	l.tlsWrapLock.RLock()
	defer l.tlsWrapLock.RUnlock()
	return l.tlsWrap
}

// ReloadTLS swaps the TLS wrapper. This is useful when upgrading or
// downgrading TLS connections.
func (l *RaftLayer) ReloadTLS(tlsWrap tlsutil.Wrapper) {
	l.tlsWrapLock.Lock()
	defer l.tlsWrapLock.Unlock()
	l.tlsWrap = tlsWrap
}
*/
// Addr is used to return the address of the listener
func (l *raftLayer) Addr() net.Addr {
	return l.addr
}

// Dial is used to create a new outgoing connection
func (l *raftLayer) Dial(address raft.ServerAddress, timeout time.Duration) (net.Conn, error) {
	tlsConfig := &tls.Config{}

	/*
		if caCert != nil {
			pool := x509.NewCertPool()
			pool.AddCert(caCert)
			tlsConfig.RootCAs = pool
			tlsConfig.ClientCAs = pool
		}*/
	l.logger.Debug("creating rpc dialer", "host", tlsConfig.ServerName)

	tlsConfig.NextProtos = []string{consts.RaftStorageALPN}
	dialer := &net.Dialer{
		Timeout: timeout,
	}
	return tls.DialWithDialer(dialer, "tcp", string(address), tlsConfig)
}

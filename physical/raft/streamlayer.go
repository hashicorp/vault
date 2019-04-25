package raft

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/tls"
	"crypto/x509"
	"errors"
	fmt "fmt"
	"net"
	"sync"
	"time"

	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	physicalstd "github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
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

	// TLS config
	certBytes  []byte
	parsedCert *x509.Certificate
	parsedKey  *ecdsa.PrivateKey
}

func NewRaftLayer(logger log.Logger, conf *physicalstd.NetworkConfig) (*raftLayer, error) {
	switch {
	case conf.Addr == nil:
		// Clustering disabled on the server, don't try to look for params
		return nil, errors.New("no raft addr found")

	case conf.KeyParams == nil:
		logger.Error("no key params found loading raft cluster TLS information")
		return nil, errors.New("no raft cluster key params found")

	case conf.KeyParams.X == nil, conf.KeyParams.Y == nil, conf.KeyParams.D == nil:
		logger.Error("failed to parse raft cluster key due to missing params")
		return nil, errors.New("failed to parse raft cluster key")

	case conf.KeyParams.Type != certutil.PrivateKeyTypeP521:
		logger.Error("unknown raft cluster key type", "key_type", conf.KeyParams.Type)
		return nil, errors.New("failed to find valid raft cluster key type")

	case len(conf.Cert) == 0:
		logger.Error("no cluster cert found")
		return nil, errors.New("no cluster cert found")

	}

	locCert := make([]byte, len(conf.Cert))
	copy(locCert, conf.Cert)

	parsedCert, err := x509.ParseCertificate(conf.Cert)
	if err != nil {
		logger.Error("failed parsing raft cluster certificate", "error", err)
		return nil, errwrap.Wrapf("error parsing raft cluster certificate: {{err}}", err)
	}

	return &raftLayer{
		addr:    conf.Addr,
		connCh:  make(chan net.Conn),
		closeCh: make(chan struct{}),
		logger:  logger,

		certBytes:  locCert,
		parsedCert: parsedCert,
		parsedKey: &ecdsa.PrivateKey{
			PublicKey: ecdsa.PublicKey{
				Curve: elliptic.P521(),
				X:     conf.KeyParams.X,
				Y:     conf.KeyParams.Y,
			},
			D: conf.KeyParams.D,
		},
	}, nil
}

func (l *raftLayer) ClientLookup(ctx context.Context, requestInfo *tls.CertificateRequestInfo) (*tls.Certificate, error) {
	for _, subj := range requestInfo.AcceptableCAs {
		if bytes.Equal(subj, l.parsedCert.RawIssuer) {
			localCert := make([]byte, len(l.certBytes))
			copy(localCert, l.certBytes)

			return &tls.Certificate{
				Certificate: [][]byte{localCert},
				PrivateKey:  l.parsedKey,
				Leaf:        l.parsedCert,
			}, nil
		}
	}

	return nil, nil
}
func (l *raftLayer) ServerLookup(context.Context, *tls.ClientHelloInfo) (*tls.Certificate, error) {
	if l.parsedKey == nil {
		return nil, errors.New("got raft connection but no local cert")
	}

	localCert := make([]byte, len(l.certBytes))
	copy(localCert, l.certBytes)

	return &tls.Certificate{
		Certificate: [][]byte{localCert},
		PrivateKey:  l.parsedKey,
		Leaf:        l.parsedCert,
	}, nil
}
func (l *raftLayer) CALookup(context.Context) (*x509.Certificate, error) {
	return l.parsedCert, nil
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

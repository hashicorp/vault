package vault

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	math "math"
	"net/http"
	"net/url"
	"sync"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/forwarding"
	"github.com/hashicorp/vault/vault/replication"
	cache "github.com/patrickmn/go-cache"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	clusterListenerAcceptDeadline = 500 * time.Millisecond

	// PerformanceReplicationALPN is the negotiated protocol used for
	// performance replication.
	PerformanceReplicationALPN = "replication_v1"

	// DRReplicationALPN is the negotiated protocol used for
	// dr replication.
	DRReplicationALPN = "replication_dr_v1"

	perfStandbyALPN = "perf_standby_v1"

	requestForwardingALPN = "req_fw_sb-act_v1"
)

var (
	// Making this a package var allows tests to modify
	HeartbeatInterval = 5 * time.Second
)

type requestForwardingHandler struct {
	fws         *http2.Server
	fwRPCServer *grpc.Server
	logger      log.Logger
	ha          bool
	core        *Core
	stopCh      chan struct{}
}

type requestForwardingClusterClient struct {
	core *Core
}

// NewRequestForwardingHandler creates a cluster handler for use with request
// forwarding.
func NewRequestForwardingHandler(c *Core, fws *http2.Server, perfStandbySlots chan struct{}, perfStandbyRepCluster *replication.Cluster, perfStandbyCache *cache.Cache) (*requestForwardingHandler, error) {
	// Resolve locally to avoid races
	ha := c.ha != nil

	fwRPCServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time: 2 * HeartbeatInterval,
		}),
		grpc.MaxRecvMsgSize(math.MaxInt32),
		grpc.MaxSendMsgSize(math.MaxInt32),
	)

	if ha && c.clusterHandler != nil {
		RegisterRequestForwardingServer(fwRPCServer, &forwardedRequestRPCServer{
			core:                  c,
			handler:               c.clusterHandler,
			perfStandbySlots:      perfStandbySlots,
			perfStandbyRepCluster: perfStandbyRepCluster,
			perfStandbyCache:      perfStandbyCache,
		})
	}

	return &requestForwardingHandler{
		fws:         fws,
		fwRPCServer: fwRPCServer,
		ha:          ha,
		logger:      c.logger.Named("request-forward"),
		core:        c,
		stopCh:      make(chan struct{}),
	}, nil
}

// ClientLookup satisfies the ClusterClient interface and returns the ha tls
// client certs.
func (c *requestForwardingClusterClient) ClientLookup(ctx context.Context, requestInfo *tls.CertificateRequestInfo) (*tls.Certificate, error) {
	parsedCert := c.core.localClusterParsedCert.Load().(*x509.Certificate)
	if parsedCert == nil {
		return nil, nil
	}
	currCert := c.core.localClusterCert.Load().([]byte)
	if len(currCert) == 0 {
		return nil, nil
	}
	localCert := make([]byte, len(currCert))
	copy(localCert, currCert)

	for _, subj := range requestInfo.AcceptableCAs {
		if bytes.Equal(subj, parsedCert.RawIssuer) {
			return &tls.Certificate{
				Certificate: [][]byte{localCert},
				PrivateKey:  c.core.localClusterPrivateKey.Load().(*ecdsa.PrivateKey),
				Leaf:        c.core.localClusterParsedCert.Load().(*x509.Certificate),
			}, nil
		}
	}

	return nil, nil
}

// ServerLookup satisfies the ClusterHandler interface and returns the server's
// tls certs.
func (rf *requestForwardingHandler) ServerLookup(ctx context.Context, clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	currCert := rf.core.localClusterCert.Load().([]byte)
	if len(currCert) == 0 {
		return nil, fmt.Errorf("got forwarding connection but no local cert")
	}

	localCert := make([]byte, len(currCert))
	copy(localCert, currCert)

	return &tls.Certificate{
		Certificate: [][]byte{localCert},
		PrivateKey:  rf.core.localClusterPrivateKey.Load().(*ecdsa.PrivateKey),
		Leaf:        rf.core.localClusterParsedCert.Load().(*x509.Certificate),
	}, nil
}

// CALookup satisfies the ClusterHandler interface and returns the ha ca cert.
func (rf *requestForwardingHandler) CALookup(ctx context.Context) (*x509.Certificate, error) {
	parsedCert := rf.core.localClusterParsedCert.Load().(*x509.Certificate)

	if parsedCert == nil {
		return nil, fmt.Errorf("forwarding connection client but no local cert")
	}

	return parsedCert, nil
}

// Handoff serves a request forwarding connection.
func (rf *requestForwardingHandler) Handoff(ctx context.Context, shutdownWg *sync.WaitGroup, closeCh chan struct{}, tlsConn *tls.Conn) error {
	if !rf.ha {
		tlsConn.Close()
		return nil
	}

	rf.logger.Debug("got request forwarding connection")

	shutdownWg.Add(2)
	// quitCh is used to close the connection and the second
	// goroutine if the server closes before closeCh.
	quitCh := make(chan struct{})
	go func() {
		select {
		case <-quitCh:
		case <-closeCh:
		case <-rf.stopCh:
		}
		tlsConn.Close()
		shutdownWg.Done()
	}()

	go func() {
		rf.fws.ServeConn(tlsConn, &http2.ServeConnOpts{
			Handler: rf.fwRPCServer,
			BaseConfig: &http.Server{
				ErrorLog: rf.logger.StandardLogger(nil),
			},
		})

		// close the quitCh which will close the connection and
		// the other goroutine.
		close(quitCh)
		shutdownWg.Done()
	}()

	return nil
}

// Stop stops the request forwarding server and closes connections.
func (rf *requestForwardingHandler) Stop() error {
	// Give some time for existing RPCs to drain.
	time.Sleep(clusterListenerAcceptDeadline)
	close(rf.stopCh)
	rf.fwRPCServer.Stop()
	return nil
}

// Starts the listeners and servers necessary to handle forwarded requests
func (c *Core) startForwarding(ctx context.Context) error {
	c.logger.Debug("cluster listener setup function")
	defer c.logger.Debug("leaving cluster listener setup function")

	// Clean up in case we have transitioned from a client to a server
	c.requestForwardingConnectionLock.Lock()
	c.clearForwardingClients()
	c.requestForwardingConnectionLock.Unlock()

	// Resolve locally to avoid races
	if c.ha == nil || c.clusterListener == nil {
		return nil
	}

	perfStandbyRepCluster, perfStandbyCache, perfStandbySlots, err := c.perfStandbyClusterHandler()
	if err != nil {
		return err
	}

	handler, err := NewRequestForwardingHandler(c, c.clusterListener.Server(), perfStandbySlots, perfStandbyRepCluster, perfStandbyCache)
	if err != nil {
		return err
	}

	c.clusterListener.AddHandler(requestForwardingALPN, handler)

	return nil
}

func (c *Core) stopForwarding() {
	if c.clusterListener != nil {
		c.clusterListener.StopHandler(requestForwardingALPN)
		c.clusterListener.StopHandler(perfStandbyALPN)
	}
}

// refreshRequestForwardingConnection ensures that the client/transport are
// alive and that the current active address value matches the most
// recently-known address.
func (c *Core) refreshRequestForwardingConnection(ctx context.Context, clusterAddr string) error {
	c.logger.Debug("refreshing forwarding connection")
	defer c.logger.Debug("done refreshing forwarding connection")

	c.requestForwardingConnectionLock.Lock()
	defer c.requestForwardingConnectionLock.Unlock()

	// Clean things up first
	c.clearForwardingClients()

	// If we don't have anything to connect to, just return
	if clusterAddr == "" {
		return nil
	}

	clusterURL, err := url.Parse(clusterAddr)
	if err != nil {
		c.logger.Error("error parsing cluster address attempting to refresh forwarding connection", "error", err)
		return err
	}

	parsedCert := c.localClusterParsedCert.Load().(*x509.Certificate)
	if parsedCert == nil {
		c.logger.Error("no request forwarding cluster certificate found")
		return errors.New("no request forwarding cluster certificate found")
	}

	if c.clusterListener != nil {
		c.clusterListener.AddClient(requestForwardingALPN, &requestForwardingClusterClient{
			core: c,
		})
	}

	// Set up grpc forwarding handling
	// It's not really insecure, but we have to dial manually to get the
	// ALPN header right. It's just "insecure" because GRPC isn't managing
	// the TLS state.
	dctx, cancelFunc := context.WithCancel(ctx)
	c.rpcClientConn, err = grpc.DialContext(dctx, clusterURL.Host,
		grpc.WithDialer(c.getGRPCDialer(ctx, requestForwardingALPN, parsedCert.Subject.CommonName, parsedCert)),
		grpc.WithInsecure(), // it's not, we handle it in the dialer
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time: 2 * HeartbeatInterval,
		}),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(math.MaxInt32),
			grpc.MaxCallSendMsgSize(math.MaxInt32),
		))
	if err != nil {
		cancelFunc()
		c.logger.Error("err setting up forwarding rpc client", "error", err)
		return err
	}
	c.rpcClientConnContext = dctx
	c.rpcClientConnCancelFunc = cancelFunc
	c.rpcForwardingClient = &forwardingClient{
		RequestForwardingClient: NewRequestForwardingClient(c.rpcClientConn),
		core:                    c,
		echoTicker:              time.NewTicker(HeartbeatInterval),
		echoContext:             dctx,
	}
	c.rpcForwardingClient.startHeartbeat()

	return nil
}

func (c *Core) clearForwardingClients() {
	c.logger.Debug("clearing forwarding clients")
	defer c.logger.Debug("done clearing forwarding clients")

	if c.rpcClientConnCancelFunc != nil {
		c.rpcClientConnCancelFunc()
		c.rpcClientConnCancelFunc = nil
	}
	if c.rpcClientConn != nil {
		c.rpcClientConn.Close()
		c.rpcClientConn = nil
	}

	c.rpcClientConnContext = nil
	c.rpcForwardingClient = nil

	if c.clusterListener != nil {
		c.clusterListener.RemoveClient(requestForwardingALPN)
	}
	c.clusterLeaderParams.Store((*ClusterLeaderParams)(nil))
}

// ForwardRequest forwards a given request to the active node and returns the
// response.
func (c *Core) ForwardRequest(req *http.Request) (int, http.Header, []byte, error) {
	c.requestForwardingConnectionLock.RLock()
	defer c.requestForwardingConnectionLock.RUnlock()

	if c.rpcForwardingClient == nil {
		return 0, nil, nil, ErrCannotForward
	}

	origPath := req.URL.Path
	defer func() {
		req.URL.Path = origPath
	}()

	req.URL.Path = req.Context().Value("original_request_path").(string)

	freq, err := forwarding.GenerateForwardedRequest(req)
	if err != nil {
		c.logger.Error("error creating forwarding RPC request", "error", err)
		return 0, nil, nil, fmt.Errorf("error creating forwarding RPC request")
	}
	if freq == nil {
		c.logger.Error("got nil forwarding RPC request")
		return 0, nil, nil, fmt.Errorf("got nil forwarding RPC request")
	}
	resp, err := c.rpcForwardingClient.ForwardRequest(c.rpcClientConnContext, freq)
	if err != nil {
		c.logger.Error("error during forwarded RPC request", "error", err)
		return 0, nil, nil, fmt.Errorf("error during forwarding RPC request")
	}

	var header http.Header
	if resp.HeaderEntries != nil {
		header = make(http.Header)
		for k, v := range resp.HeaderEntries {
			header[k] = v.Values
		}
	}

	// If we are a perf standby and the request was forwarded to the active node
	// we should attempt to wait for the WAL to ship to offer best effort read after
	// write guarantees
	if c.perfStandby && resp.LastRemoteWal > 0 {
		WaitUntilWALShipped(req.Context(), c, resp.LastRemoteWal)
	}

	return int(resp.StatusCode), header, resp.Body, nil
}

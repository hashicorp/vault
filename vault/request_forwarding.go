// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/armon/go-metrics"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/forwarding"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/cluster"
	"github.com/hashicorp/vault/vault/replication"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	NotHAMember       = "node is not in HA cluster membership"
	StatusNotHAMember = status.Error(codes.FailedPrecondition, NotHAMember)
)

const haNodeIDKey = "ha_node_id"

func haIDFromContext(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}
	res := md.Get(haNodeIDKey)
	if len(res) == 0 {
		return "", false
	}
	return res[0], true
}

// haMembershipServerCheck extracts the client's HA node ID from the context
// and checks if this client has been removed. The function returns
// StatusNotHAMember if the client has been removed
func haMembershipServerCheck(ctx context.Context, c *Core, haBackend physical.RemovableNodeHABackend) error {
	if haBackend == nil {
		return nil
	}
	nodeID, ok := haIDFromContext(ctx)
	if !ok {
		return nil
	}
	removed, err := haBackend.IsNodeRemoved(ctx, nodeID)
	if err != nil {
		c.logger.Error("failed to check if node is removed", "error", err)
		return err
	}
	if removed {
		return StatusNotHAMember
	}
	return nil
}

func haMembershipUnaryServerInterceptor(c *Core, haBackend physical.RemovableNodeHABackend) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		err = haMembershipServerCheck(ctx, c, haBackend)
		if err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func haMembershipStreamServerInterceptor(c *Core, haBackend physical.RemovableNodeHABackend) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := haMembershipServerCheck(ss.Context(), c, haBackend)
		if err != nil {
			return err
		}
		return handler(srv, ss)
	}
}

// haMembershipClientCheck checks if the given error from the server
// is StatusNotHAMember. If so, the client will mark itself as removed
// and shutdown
func haMembershipClientCheck(err error, c *Core, haBackend physical.RemovableNodeHABackend) {
	if !errors.Is(err, StatusNotHAMember) {
		return
	}
	removeErr := haBackend.RemoveSelf()
	if removeErr != nil {
		c.logger.Debug("failed to remove self", "error", removeErr)
	}
	c.shutdownRemovedNode()
}

func haMembershipUnaryClientInterceptor(c *Core, haBackend physical.RemovableNodeHABackend) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if haBackend == nil {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		ctx = metadata.AppendToOutgoingContext(ctx, haNodeIDKey, haBackend.NodeID())
		err := invoker(ctx, method, req, reply, cc, opts...)
		haMembershipClientCheck(err, c, haBackend)
		return err
	}
}

func haMembershipStreamClientInterceptor(c *Core, haBackend physical.RemovableNodeHABackend) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if haBackend == nil {
			return streamer(ctx, desc, cc, method, opts...)
		}
		ctx = metadata.AppendToOutgoingContext(ctx, haNodeIDKey, haBackend.NodeID())
		stream, err := streamer(ctx, desc, cc, method, opts...)
		haMembershipClientCheck(err, c, haBackend)
		return stream, err
	}
}

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
func NewRequestForwardingHandler(c *Core, fws *http2.Server, perfStandbySlots chan struct{}, perfStandbyRepCluster *replication.Cluster) (*requestForwardingHandler, error) {
	// Resolve locally to avoid races
	ha := c.ha != nil
	removableHABackend := c.getRemovableHABackend()

	fwRPCServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time: 2 * c.clusterHeartbeatInterval,
		}),
		grpc.MaxRecvMsgSize(math.MaxInt32),
		grpc.MaxSendMsgSize(math.MaxInt32),
		grpc.StreamInterceptor(haMembershipStreamServerInterceptor(c, removableHABackend)),
		grpc.UnaryInterceptor(haMembershipUnaryServerInterceptor(c, removableHABackend)),
	)

	if ha && c.clusterHandler != nil {
		RegisterRequestForwardingServer(fwRPCServer, &forwardedRequestRPCServer{
			core:                  c,
			handler:               c.clusterHandler,
			perfStandbySlots:      perfStandbySlots,
			perfStandbyRepCluster: perfStandbyRepCluster,
			raftFollowerStates:    c.raftFollowerStates,
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

func (c *requestForwardingClusterClient) ServerName() string {
	parsedCert := c.core.localClusterParsedCert.Load().(*x509.Certificate)
	if parsedCert == nil {
		return ""
	}

	return parsedCert.Subject.CommonName
}

func (c *requestForwardingClusterClient) CACert(ctx context.Context) *x509.Certificate {
	return c.core.localClusterParsedCert.Load().(*x509.Certificate)
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
func (rf *requestForwardingHandler) CALookup(ctx context.Context) ([]*x509.Certificate, error) {
	parsedCert := rf.core.localClusterParsedCert.Load().(*x509.Certificate)

	if parsedCert == nil {
		return nil, fmt.Errorf("forwarding connection client but no local cert")
	}

	return []*x509.Certificate{parsedCert}, nil
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
	time.Sleep(cluster.ListenerAcceptDeadline)
	close(rf.stopCh)
	rf.fwRPCServer.Stop()
	return nil
}

// Starts the listeners and servers necessary to handle forwarded requests
func (c *Core) startForwarding(ctx context.Context) error {
	c.logger.Debug("request forwarding setup function")
	defer c.logger.Debug("leaving request forwarding setup function")

	// Clean up in case we have transitioned from a client to a server
	c.requestForwardingConnectionLock.Lock()
	c.clearForwardingClients()
	c.requestForwardingConnectionLock.Unlock()

	clusterListener := c.getClusterListener()
	if c.ha == nil || clusterListener == nil {
		c.logger.Debug("request forwarding not setup")
		return nil
	}

	perfStandbyRepCluster, perfStandbySlots, err := c.perfStandbyClusterHandler()
	if err != nil {
		return err
	}

	handler, err := NewRequestForwardingHandler(c, clusterListener.Server(), perfStandbySlots, perfStandbyRepCluster)
	if err != nil {
		return err
	}

	clusterListener.AddHandler(consts.RequestForwardingALPN, handler)

	return nil
}

func (c *Core) stopForwarding() {
	clusterListener := c.getClusterListener()
	if clusterListener != nil {
		clusterListener.StopHandler(consts.RequestForwardingALPN)
		clusterListener.StopHandler(consts.PerfStandbyALPN)
	}
	c.removeAllPerfStandbySecondaries()
}

// refreshRequestForwardingConnection ensures that the client/transport are
// alive and that the current active address value matches the most
// recently-known address.
func (c *Core) refreshRequestForwardingConnection(ctx context.Context, clusterAddr string) error {
	c.logger.Debug("refreshing forwarding connection", "clusterAddr", clusterAddr)
	defer c.logger.Debug("done refreshing forwarding connection", "clusterAddr", clusterAddr)

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

	clusterListener := c.getClusterListener()
	if clusterListener == nil {
		c.logger.Error("no cluster listener configured")
		return nil
	}

	clusterListener.AddClient(consts.RequestForwardingALPN, &requestForwardingClusterClient{
		core: c,
	})

	removableHABackend := c.getRemovableHABackend()

	// Set up grpc forwarding handling
	// It's not really insecure, but we have to dial manually to get the
	// ALPN header right. It's just "insecure" because GRPC isn't managing
	// the TLS state.
	dctx, cancelFunc := context.WithCancel(ctx)
	c.rpcClientConn, err = grpc.DialContext(dctx, clusterURL.Host,
		grpc.WithDialer(clusterListener.GetDialerFunc(ctx, consts.RequestForwardingALPN)),
		grpc.WithInsecure(), // it's not, we handle it in the dialer
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time: 2 * c.clusterHeartbeatInterval,
		}),
		grpc.WithStreamInterceptor(haMembershipStreamClientInterceptor(c, removableHABackend)),
		grpc.WithUnaryInterceptor(haMembershipUnaryClientInterceptor(c, removableHABackend)),
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
	duration := c.clusterHeartbeatInterval
	if duration <= 0 {
		duration = time.Second * 5
	}
	c.rpcForwardingClient = &forwardingClient{
		RequestForwardingClient: NewRequestForwardingClient(c.rpcClientConn),
		core:                    c,
		echoTicker:              time.NewTicker(duration),
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

	clusterListener := c.getClusterListener()
	if clusterListener != nil {
		clusterListener.RemoveClient(consts.RequestForwardingALPN)
	}
	c.clusterLeaderParams.Store((*ClusterLeaderParams)(nil))
	c.rpcLastSuccessfulHeartbeat.Store(time.Time{})
}

// ForwardRequest forwards a given request to the active node and returns the
// response.
func (c *Core) ForwardRequest(req *http.Request) (int, http.Header, []byte, error) {
	// checking if the node is perfStandby here to avoid a deadlock between
	// Core.stateLock and Core.requestForwardingConnectionLock
	isPerfStandby := c.PerfStandby()
	c.requestForwardingConnectionLock.RLock()
	defer c.requestForwardingConnectionLock.RUnlock()

	if c.rpcForwardingClient == nil {
		return 0, nil, nil, ErrCannotForward
	}

	defer metrics.MeasureSince([]string{"ha", "rpc", "client", "forward"}, time.Now())

	origPath := req.URL.Path
	defer func() {
		req.URL.Path = origPath
	}()

	path, ok := logical.ContextOriginalRequestPathValue(req.Context())
	if !ok {
		return 0, nil, nil, errors.New("error extracting request path for forwarding RPC request")
	}

	req.URL.Path = path

	freq, err := forwarding.GenerateForwardedRequest(req)
	if err != nil {
		c.logger.Error("error creating forwarding RPC request", "error", err)
		return 0, nil, nil, fmt.Errorf("error creating forwarding RPC request")
	}
	if freq == nil {
		c.logger.Error("got nil forwarding RPC request")
		return 0, nil, nil, fmt.Errorf("got nil forwarding RPC request")
	}
	resp, err := c.rpcForwardingClient.ForwardRequest(req.Context(), freq)
	if err != nil {
		metrics.IncrCounter([]string{"ha", "rpc", "client", "forward", "errors"}, 1)
		c.logger.Error("error during forwarded RPC request", "error", err)

		if errors.Is(err, StatusNotHAMember) {
			return 0, nil, nil, fmt.Errorf("error during forwarding RPC request: %w", err)
		}
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
	if isPerfStandby && resp.LastRemoteWal > 0 {
		c.EntWaitUntilWALShipped(req.Context(), resp.LastRemoteWal)
	}

	return int(resp.StatusCode), header, resp.Body, nil
}

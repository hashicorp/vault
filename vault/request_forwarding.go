package vault

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	math "math"
	"net"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	cache "github.com/patrickmn/go-cache"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/forwarding"
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

type SecondaryConnsCacheVals struct {
	ID         string
	Token      string
	Connection net.Conn
	Mode       consts.ReplicationState
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
	ha := c.ha != nil

	var perfStandbyRepCluster *ReplicatedCluster
	if ha {
		id, err := uuid.GenerateUUID()
		if err != nil {
			return err
		}

		perfStandbyRepCluster = &ReplicatedCluster{
			State:              consts.ReplicationPerformanceStandby,
			ClusterID:          id,
			PrimaryClusterAddr: c.clusterAddr,
		}
		if err = c.setupReplicatedClusterPrimary(perfStandbyRepCluster); err != nil {
			return err
		}
	}

	// Get our TLS config
	tlsConfig, err := c.ClusterTLSConfig(ctx, nil, perfStandbyRepCluster)
	if err != nil {
		c.logger.Error("failed to get tls configuration when starting forwarding", "error", err)
		return err
	}

	// The server supports all of the possible protos
	tlsConfig.NextProtos = []string{"h2", requestForwardingALPN, perfStandbyALPN, PerformanceReplicationALPN, DRReplicationALPN}

	if !atomic.CompareAndSwapUint32(c.rpcServerActive, 0, 1) {
		c.logger.Warn("forwarding rpc server already running")
		return nil
	}

	fwRPCServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time: 2 * HeartbeatInterval,
		}),
		grpc.MaxRecvMsgSize(math.MaxInt32),
		grpc.MaxSendMsgSize(math.MaxInt32),
	)

	// Setup performance standby RPC servers
	perfStandbyCount := 0
	if !c.IsDRSecondary() && !c.disablePerfStandby {
		perfStandbyCount = c.perfStandbyCount()
	}
	perfStandbySlots := make(chan struct{}, perfStandbyCount)

	perfStandbyCache := cache.New(2*HeartbeatInterval, 1*time.Second)
	perfStandbyCache.OnEvicted(func(secondaryID string, _ interface{}) {
		c.logger.Debug("removing performance standby", "id", secondaryID)
		c.removePerfStandbySecondary(context.Background(), secondaryID)
		select {
		case <-perfStandbySlots:
		default:
			c.logger.Warn("perf secondary timeout hit but no slot to free")
		}
	})

	perfStandbyReplicationRPCServer := perfStandbyRPCServer(c, perfStandbyCache)

	if ha && c.clusterHandler != nil {
		RegisterRequestForwardingServer(fwRPCServer, &forwardedRequestRPCServer{
			core:                  c,
			handler:               c.clusterHandler,
			perfStandbySlots:      perfStandbySlots,
			perfStandbyRepCluster: perfStandbyRepCluster,
			perfStandbyCache:      perfStandbyCache,
		})
	}

	// Create the HTTP/2 server that will be shared by both RPC and regular
	// duties. Doing it this way instead of listening via the server and gRPC
	// allows us to re-use the same port via ALPN. We can just tell the server
	// to serve a given conn and which handler to use.
	fws := &http2.Server{
		// Our forwarding connections heartbeat regularly so anything else we
		// want to go away/get cleaned up pretty rapidly
		IdleTimeout: 5 * HeartbeatInterval,
	}

	// Shutdown coordination logic
	shutdown := new(uint32)
	shutdownWg := &sync.WaitGroup{}

	for _, addr := range c.clusterListenerAddrs {
		shutdownWg.Add(1)

		// Force a local resolution to avoid data races
		laddr := addr

		// Start our listening loop
		go func() {
			defer shutdownWg.Done()

			// closeCh is used to shutdown the spawned goroutines once this
			// function returns
			closeCh := make(chan struct{})
			defer func() {
				close(closeCh)
			}()

			if c.logger.IsInfo() {
				c.logger.Info("core/startClusterListener: starting listener", "listener_address", laddr)
			}

			// Create a TCP listener. We do this separately and specifically
			// with TCP so that we can set deadlines.
			tcpLn, err := net.ListenTCP("tcp", laddr)
			if err != nil {
				c.logger.Error("core/startClusterListener: error starting listener", "error", err)
				return
			}

			// Wrap the listener with TLS
			tlsLn := tls.NewListener(tcpLn, tlsConfig)
			defer tlsLn.Close()

			if c.logger.IsInfo() {
				c.logger.Info("core/startClusterListener: serving cluster requests", "cluster_listen_address", tlsLn.Addr())
			}

			for {
				if atomic.LoadUint32(shutdown) > 0 {
					return
				}

				// Set the deadline for the accept call. If it passes we'll get
				// an error, causing us to check the condition at the top
				// again.
				tcpLn.SetDeadline(time.Now().Add(clusterListenerAcceptDeadline))

				// Accept the connection
				conn, err := tlsLn.Accept()
				if err != nil {
					if err, ok := err.(net.Error); ok && !err.Timeout() {
						c.logger.Debug("non-timeout error accepting on cluster port", "error", err)
					}
					if conn != nil {
						conn.Close()
					}
					continue
				}
				if conn == nil {
					continue
				}

				// Type assert to TLS connection and handshake to populate the
				// connection state
				tlsConn := conn.(*tls.Conn)

				// Set a deadline for the handshake. This will cause clients
				// that don't successfully auth to be kicked out quickly.
				// Cluster connections should be reliable so being marginally
				// aggressive here is fine.
				err = tlsConn.SetDeadline(time.Now().Add(30 * time.Second))
				if err != nil {
					if c.logger.IsDebug() {
						c.logger.Debug("error setting deadline for cluster connection", "error", err)
					}
					tlsConn.Close()
					continue
				}

				err = tlsConn.Handshake()
				if err != nil {
					if c.logger.IsDebug() {
						c.logger.Debug("error handshaking cluster connection", "error", err)
					}
					tlsConn.Close()
					continue
				}

				// Now, set it back to unlimited
				err = tlsConn.SetDeadline(time.Time{})
				if err != nil {
					if c.logger.IsDebug() {
						c.logger.Debug("error setting deadline for cluster connection", "error", err)
					}
					tlsConn.Close()
					continue
				}

				switch tlsConn.ConnectionState().NegotiatedProtocol {
				case requestForwardingALPN:
					if !ha {
						tlsConn.Close()
						continue
					}

					c.logger.Debug("got request forwarding connection")

					shutdownWg.Add(2)
					// quitCh is used to close the connection and the second
					// goroutine if the server closes before closeCh.
					quitCh := make(chan struct{})
					go func() {
						select {
						case <-quitCh:
						case <-closeCh:
						}
						tlsConn.Close()
						shutdownWg.Done()
					}()

					go func() {
						fws.ServeConn(tlsConn, &http2.ServeConnOpts{
							Handler: fwRPCServer,
							BaseConfig: &http.Server{
								ErrorLog: c.logger.StandardLogger(nil),
							},
						})
						// close the quitCh which will close the connection and
						// the other goroutine.
						close(quitCh)
						shutdownWg.Done()
					}()

				case PerformanceReplicationALPN, DRReplicationALPN, perfStandbyALPN:
					handleReplicationConn(ctx, c, shutdownWg, closeCh, fws, perfStandbyReplicationRPCServer, perfStandbyCache, tlsConn)
				default:
					c.logger.Debug("unknown negotiated protocol on cluster port")
					tlsConn.Close()
					continue
				}
			}
		}()
	}

	// This is in its own goroutine so that we don't block the main thread, and
	// thus we use atomic and channels to coordinate
	// However, because you can't query the status of a channel, we set a bool
	// here while we have the state lock to know whether to actually send a
	// shutdown (e.g. whether the channel will block). See issue #2083.
	c.clusterListenersRunning = true
	go func() {
		// If we get told to shut down...
		<-c.clusterListenerShutdownCh

		// Stop the RPC server
		c.logger.Info("shutting down forwarding rpc listeners")
		fwRPCServer.Stop()

		// Set the shutdown flag. This will cause the listeners to shut down
		// within the deadline in clusterListenerAcceptDeadline
		atomic.StoreUint32(shutdown, 1)
		c.logger.Info("forwarding rpc listeners stopped")

		// Wait for them all to shut down
		shutdownWg.Wait()
		c.logger.Info("rpc listeners successfully shut down")

		// Clear us up to run this function again
		atomic.StoreUint32(c.rpcServerActive, 0)

		// Tell the main thread that shutdown is done.
		c.clusterListenerShutdownSuccessCh <- struct{}{}
	}()

	return nil
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

	// Set up grpc forwarding handling
	// It's not really insecure, but we have to dial manually to get the
	// ALPN header right. It's just "insecure" because GRPC isn't managing
	// the TLS state.
	dctx, cancelFunc := context.WithCancel(ctx)
	c.rpcClientConn, err = grpc.DialContext(dctx, clusterURL.Host,
		grpc.WithDialer(c.getGRPCDialer(ctx, requestForwardingALPN, "", nil, nil, nil)),
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

// getGRPCDialer is used to return a dialer that has the correct TLS
// configuration. Otherwise gRPC tries to be helpful and stomps all over our
// NextProtos.
func (c *Core) getGRPCDialer(ctx context.Context, alpnProto, serverName string, caCert *x509.Certificate, repClusters *ReplicatedClusters, perfStandbyCluster *ReplicatedCluster) func(string, time.Duration) (net.Conn, error) {
	return func(addr string, timeout time.Duration) (net.Conn, error) {
		tlsConfig, err := c.ClusterTLSConfig(ctx, repClusters, perfStandbyCluster)
		if err != nil {
			c.logger.Error("failed to get tls configuration", "error", err)
			return nil, err
		}
		if serverName != "" {
			tlsConfig.ServerName = serverName
		}
		if caCert != nil {
			pool := x509.NewCertPool()
			pool.AddCert(caCert)
			tlsConfig.RootCAs = pool
			tlsConfig.ClientCAs = pool
		}
		c.logger.Debug("creating rpc dialer", "host", tlsConfig.ServerName)

		tlsConfig.NextProtos = []string{alpnProto}
		dialer := &net.Dialer{
			Timeout: timeout,
		}
		return tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
	}
}

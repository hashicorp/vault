package vault

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/vault/helper/forwarding"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	clusterListenerAcceptDeadline = 500 * time.Millisecond
	heartbeatInterval             = 30 * time.Second
	requestForwardingALPN         = "req_fw_sb-act_v1"
)

// Starts the listeners and servers necessary to handle forwarded requests
func (c *Core) startForwarding() error {
	c.logger.Trace("core: cluster listener setup function")
	defer c.logger.Trace("core: leaving cluster listener setup function")

	// Clean up in case we have transitioned from a client to a server
	c.requestForwardingConnectionLock.Lock()
	c.clearForwardingClients()
	c.requestForwardingConnectionLock.Unlock()

	// Resolve locally to avoid races
	ha := c.ha != nil

	// Get our TLS config
	tlsConfig, err := c.ClusterTLSConfig()
	if err != nil {
		c.logger.Error("core: failed to get tls configuration when starting forwarding", "error", err)
		return err
	}

	// The server supports all of the possible protos
	tlsConfig.NextProtos = []string{"h2", requestForwardingALPN}

	// Create our RPC server and register the request handler server
	c.clusterParamsLock.Lock()

	if c.rpcServer != nil {
		c.logger.Warn("core: forwarding rpc server already running")
		return nil
	}

	c.rpcServer = grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time: 2 * heartbeatInterval,
		}),
	)

	if ha && c.clusterHandler != nil {
		RegisterRequestForwardingServer(c.rpcServer, &forwardedRequestRPCServer{
			core:    c,
			handler: c.clusterHandler,
		})
	}
	c.clusterParamsLock.Unlock()

	// Create the HTTP/2 server that will be shared by both RPC and regular
	// duties. Doing it this way instead of listening via the server and gRPC
	// allows us to re-use the same port via ALPN. We can just tell the server
	// to serve a given conn and which handler to use.
	fws := &http2.Server{}

	// Shutdown coordination logic
	var shutdown uint32
	shutdownWg := &sync.WaitGroup{}

	for _, addr := range c.clusterListenerAddrs {
		shutdownWg.Add(1)

		// Force a local resolution to avoid data races
		laddr := addr

		// Start our listening loop
		go func() {
			defer shutdownWg.Done()

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
				if atomic.LoadUint32(&shutdown) > 0 {
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
						c.logger.Debug("core: non-timeout error accepting on cluster port", "error", err)
					}
					if conn != nil {
						conn.Close()
					}
					continue
				}
				if conn == nil {
					continue
				}
				defer conn.Close()

				// Type assert to TLS connection and handshake to populate the
				// connection state
				tlsConn := conn.(*tls.Conn)
				err = tlsConn.Handshake()
				if err != nil {
					if c.logger.IsDebug() {
						c.logger.Debug("core: error handshaking cluster connection", "error", err)
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

					c.logger.Trace("core: got request forwarding connection")
					c.clusterParamsLock.RLock()
					rpcServer := c.rpcServer
					c.clusterParamsLock.RUnlock()
					go func() {
						fws.ServeConn(tlsConn, &http2.ServeConnOpts{
							Handler: rpcServer,
						})
						tlsConn.Close()
					}()

				default:
					c.logger.Debug("core: unknown negotiated protocol on cluster port")
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
		c.logger.Info("core: shutting down forwarding rpc listeners")
		c.clusterParamsLock.Lock()
		c.rpcServer.Stop()
		c.rpcServer = nil
		c.clusterParamsLock.Unlock()
		c.logger.Info("core: forwarding rpc listeners stopped")

		// Set the shutdown flag. This will cause the listeners to shut down
		// within the deadline in clusterListenerAcceptDeadline
		atomic.StoreUint32(&shutdown, 1)

		// Wait for them all to shut down
		shutdownWg.Wait()
		c.logger.Info("core: rpc listeners successfully shut down")

		// Tell the main thread that shutdown is done.
		c.clusterListenerShutdownSuccessCh <- struct{}{}
	}()

	return nil
}

// refreshRequestForwardingConnection ensures that the client/transport are
// alive and that the current active address value matches the most
// recently-known address.
func (c *Core) refreshRequestForwardingConnection(clusterAddr string) error {
	c.logger.Trace("core: refreshing forwarding connection")
	defer c.logger.Trace("core: done refreshing forwarding connection")

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
		c.logger.Error("core: error parsing cluster address attempting to refresh forwarding connection", "error", err)
		return err
	}

	// Set up grpc forwarding handling
	// It's not really insecure, but we have to dial manually to get the
	// ALPN header right. It's just "insecure" because GRPC isn't managing
	// the TLS state.
	ctx, cancelFunc := context.WithCancel(context.Background())
	c.rpcClientConn, err = grpc.DialContext(ctx, clusterURL.Host,
		grpc.WithDialer(c.getGRPCDialer(requestForwardingALPN, "", nil)),
		grpc.WithInsecure(), // it's not, we handle it in the dialer
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time: 2 * heartbeatInterval,
		}))
	if err != nil {
		cancelFunc()
		c.logger.Error("core: err setting up forwarding rpc client", "error", err)
		return err
	}
	c.rpcClientConnContext = ctx
	c.rpcClientConnCancelFunc = cancelFunc
	c.rpcForwardingClient = &forwardingClient{
		RequestForwardingClient: NewRequestForwardingClient(c.rpcClientConn),
		core:        c,
		echoTicker:  time.NewTicker(heartbeatInterval),
		echoContext: ctx,
	}
	c.rpcForwardingClient.startHeartbeat()

	return nil
}

func (c *Core) clearForwardingClients() {
	c.logger.Trace("core: clearing forwarding clients")
	defer c.logger.Trace("core: done clearing forwarding clients")

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

	freq, err := forwarding.GenerateForwardedRequest(req)
	if err != nil {
		c.logger.Error("core: error creating forwarding RPC request", "error", err)
		return 0, nil, nil, fmt.Errorf("error creating forwarding RPC request")
	}
	if freq == nil {
		c.logger.Error("core: got nil forwarding RPC request")
		return 0, nil, nil, fmt.Errorf("got nil forwarding RPC request")
	}
	resp, err := c.rpcForwardingClient.ForwardRequest(c.rpcClientConnContext, freq)
	if err != nil {
		c.logger.Error("core: error during forwarded RPC request", "error", err)
		return 0, nil, nil, fmt.Errorf("error during forwarding RPC request")
	}

	var header http.Header
	if resp.HeaderEntries != nil {
		header = make(http.Header)
		for k, v := range resp.HeaderEntries {
			header[k] = v.Values
		}
	}

	return int(resp.StatusCode), header, resp.Body, nil
}

// getGRPCDialer is used to return a dialer that has the correct TLS
// configuration. Otherwise gRPC tries to be helpful and stomps all over our
// NextProtos.
func (c *Core) getGRPCDialer(alpnProto, serverName string, caCert *x509.Certificate) func(string, time.Duration) (net.Conn, error) {
	return func(addr string, timeout time.Duration) (net.Conn, error) {
		tlsConfig, err := c.ClusterTLSConfig()
		if err != nil {
			c.logger.Error("core: failed to get tls configuration", "error", err)
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
		c.logger.Trace("core: creating rpc dialer", "host", tlsConfig.ServerName)

		tlsConfig.NextProtos = []string{alpnProto}
		dialer := &net.Dialer{
			Timeout: timeout,
		}
		return tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
	}
}

type forwardedRequestRPCServer struct {
	core    *Core
	handler http.Handler
}

func (s *forwardedRequestRPCServer) ForwardRequest(ctx context.Context, freq *forwarding.Request) (*forwarding.Response, error) {
	//s.core.logger.Trace("forwarding: serving rpc forwarded request")

	// Parse an http.Request out of it
	req, err := forwarding.ParseForwardedRequest(freq)
	if err != nil {
		return nil, err
	}

	// A very dummy response writer that doesn't follow normal semantics, just
	// lets you write a status code (last written wins) and a body. But it
	// meets the interface requirements.
	w := forwarding.NewRPCResponseWriter()

	resp := &forwarding.Response{}

	runRequest := func() {
		defer func() {
			// Logic here comes mostly from the Go source code
			if err := recover(); err != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				s.core.logger.Error("forwarding: panic serving request", "path", req.URL.Path, "error", err, "stacktrace", buf)
			}
		}()
		s.handler.ServeHTTP(w, req)
	}
	runRequest()
	resp.StatusCode = uint32(w.StatusCode())
	resp.Body = w.Body().Bytes()

	header := w.Header()
	if header != nil {
		resp.HeaderEntries = make(map[string]*forwarding.HeaderEntry, len(header))
		for k, v := range header {
			resp.HeaderEntries[k] = &forwarding.HeaderEntry{
				Values: v,
			}
		}
	}

	return resp, nil
}

func (s *forwardedRequestRPCServer) Echo(ctx context.Context, in *EchoRequest) (*EchoReply, error) {
	if in.ClusterAddr != "" {
		s.core.clusterPeerClusterAddrsCache.Set(in.ClusterAddr, nil, 0)
	}
	return &EchoReply{
		Message: "pong",
	}, nil
}

type forwardingClient struct {
	RequestForwardingClient

	core *Core

	echoTicker  *time.Ticker
	echoContext context.Context
}

// NOTE: we also take advantage of gRPC's keepalive bits, but as we send data
// with these requests it's useful to keep this as well
func (c *forwardingClient) startHeartbeat() {
	go func() {
		tick := func() {
			c.core.stateLock.RLock()
			clusterAddr := c.core.clusterAddr
			c.core.stateLock.RUnlock()

			ctx, cancel := context.WithTimeout(c.echoContext, 2*time.Second)
			resp, err := c.RequestForwardingClient.Echo(ctx, &EchoRequest{
				Message:     "ping",
				ClusterAddr: clusterAddr,
			})
			cancel()
			if err != nil {
				c.core.logger.Debug("forwarding: error sending echo request to active node", "error", err)
				return
			}
			if resp == nil {
				c.core.logger.Debug("forwarding: empty echo response from active node")
				return
			}
			if resp.Message != "pong" {
				c.core.logger.Debug("forwarding: unexpected echo response from active node", "message", resp.Message)
				return
			}
			c.core.logger.Trace("forwarding: successful heartbeat")
		}

		tick()

		for {
			select {
			case <-c.echoContext.Done():
				c.echoTicker.Stop()
				c.core.logger.Trace("forwarding: stopping heartbeating")
				return
			case <-c.echoTicker.C:
				tick()
			}
		}
	}()
}

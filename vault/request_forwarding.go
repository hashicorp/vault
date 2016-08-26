package vault

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/vault/helper/forwarding"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
)

const (
	clusterListenerAcceptDeadline = 500 * time.Millisecond
)

// Starts the listeners and servers necessary to handle forwarded requests
func (c *Core) startForwarding() error {
	// Clean up in case we have transitioned from a client to a server
	c.clearForwardingClients()

	// Get our base handler (for our RPC server) and our wrapped handler (for
	// straight HTTP/2 forwarding)
	baseHandler, wrappedHandler := c.clusterHandlerSetupFunc()

	// Get our TLS config
	tlsConfig, err := c.ClusterTLSConfig()
	if err != nil {
		c.logger.Error("core/startClusterListener: failed to get tls configuration", "error", err)
		return err
	}

	// The server supports all of the possible protos
	tlsConfig.NextProtos = []string{"h2", "req_fw_sb-act_v1"}

	// Create our RPC server and register the request handler server
	c.rpcServer = grpc.NewServer()
	RegisterRequestForwardingServer(c.rpcServer, &forwardedRequestRPCServer{
		core:    c,
		handler: baseHandler,
	})

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

			c.logger.Info("core/startClusterListener: starting listener")

			// Create a TCP listener. We do this separately and specifically
			// with TCP so that we can set deadlines.
			tcpLn, err := net.ListenTCP("tcp", laddr)
			if err != nil {
				c.logger.Error("core/startClusterListener: error starting listener", "error", err)
				return
			}

			// Wrap the listener with TLS
			tlsLn := tls.NewListener(tcpLn, tlsConfig)

			if c.logger.IsInfo() {
				c.logger.Info("core/startClusterListener: serving cluster requests", "cluster_listen_address", tlsLn.Addr())
			}

			for {
				if atomic.LoadUint32(&shutdown) > 0 {
					tlsLn.Close()
					return
				}

				// Set the deadline for the accept call. If it passes we'll get
				// an error, causing us to check the condition at the top
				// again.
				tcpLn.SetDeadline(time.Now().Add(clusterListenerAcceptDeadline))

				// Accept the connection
				conn, err := tlsLn.Accept()
				if err != nil {
					if conn != nil {
						conn.Close()
					}
					continue
				}

				// Type assert to TLS connection and handshake to populate the
				// connection state
				tlsConn := conn.(*tls.Conn)
				err = tlsConn.Handshake()
				if err != nil {
					if c.logger.IsDebug() {
						c.logger.Debug("core/startClusterListener/Accept: error handshaking", "error", err)
					}
					if conn != nil {
						conn.Close()
					}
					continue
				}

				switch tlsConn.ConnectionState().NegotiatedProtocol {
				case "h2":
					c.logger.Debug("core/startClusterListener/Accept: got h2 connection")
					go fws.ServeConn(conn, &http2.ServeConnOpts{
						Handler: wrappedHandler,
					})

				case "req_fw_sb-act_v1":
					c.logger.Debug("core/startClusterListener/Accept: got req_fw_sb-act_v1 connection")
					go fws.ServeConn(conn, &http2.ServeConnOpts{
						Handler: c.rpcServer,
					})

				default:
					c.logger.Debug("core/startClusterListener/Accept: unknown negotiated protocol")
					conn.Close()
					continue
				}
			}
		}()
	}

	// This is in its own goroutine so that we don't block the main thread, and
	// thus we use atomic and channels to coordinate
	go func() {
		// If we get told to shut down...
		<-c.clusterListenerShutdownCh

		// Stop the RPC server
		c.rpcServer.Stop()
		c.logger.Info("core/startClusterListener: shutting down listeners")

		// Set the shutdown flag. This will cause the listeners to shut down
		// within the deadline in clusterListenerAcceptDeadline
		atomic.StoreUint32(&shutdown, 1)

		// Wait for them all to shut down
		shutdownWg.Wait()
		c.logger.Info("core/startClusterListener: listeners successfully shut down")

		// Tell the main thread that shutdown is done.
		c.clusterListenerShutdownSuccessCh <- struct{}{}
	}()

	return nil
}

// refreshRequestForwardingConnection ensures that the client/transport are
// alive and that the current active address value matches the most
// recently-known address.
func (c *Core) refreshRequestForwardingConnection(clusterAddr string) error {
	c.requestForwardingConnectionLock.Lock()
	defer c.requestForwardingConnectionLock.Unlock()

	// It's nil but we don't have an address anyways, so exit
	if c.requestForwardingConnection == nil && clusterAddr == "" {
		return nil
	}

	// NOTE: We don't fast path the case where we have a connection because the
	// address is the same, because the cert/key could have changed if the
	// active node ended up being the same node. Before we hit this function in
	// Leader() we'll have done a hash on the advertised info to ensure that we
	// won't hit this function unnecessarily anyways.

	// Disabled, potentially, so clean up anything that might be around.
	if clusterAddr == "" {
		c.clearForwardingClients()
		return nil
	}

	clusterURL, err := url.Parse(clusterAddr)
	if err != nil {
		c.logger.Error("core/refreshRequestForwardingConnection: error parsing cluster address", "error", err)
		return err
	}

	switch os.Getenv("VAULT_USE_GRPC_REQUEST_FORWARDING") {
	case "":
		// Set up normal HTTP forwarding handling
		tlsConfig, err := c.ClusterTLSConfig()
		if err != nil {
			c.logger.Error("core/refreshRequestForwardingConnection: error fetching cluster tls configuration", "error", err)
			return err
		}
		tp := &http2.Transport{
			TLSClientConfig: tlsConfig,
		}
		c.requestForwardingConnection = &activeConnection{
			transport:   tp,
			clusterAddr: clusterAddr,
		}

	default:
		// Set up grpc forwarding handling
		// It's not really insecure, but we have to dial manually to get the
		// ALPN header right. It's just "insecure" because GRPC isn't managing
		// the TLS state.
		ctx, cancelFunc := context.WithCancel(context.Background())
		c.rpcClientConnCancelFunc = cancelFunc
		c.rpcClientConn, err = grpc.DialContext(ctx, clusterURL.Host, grpc.WithDialer(c.getGRPCDialer()), grpc.WithInsecure())
		if err != nil {
			c.logger.Error("core/refreshRequestForwardingConnection: err setting up rpc client", "error", err)
			return err
		}
		c.rpcForwardingClient = NewRequestForwardingClient(c.rpcClientConn)
	}

	return nil
}

func (c *Core) clearForwardingClients() {
	if c.requestForwardingConnection != nil {
		c.requestForwardingConnection.transport.CloseIdleConnections()
		c.requestForwardingConnection = nil
	}

	c.rpcForwardingClient = nil

	if c.rpcClientConnCancelFunc != nil {
		c.rpcClientConnCancelFunc()
		c.rpcClientConnCancelFunc = nil
	}

	if c.rpcClientConn != nil {
		c.rpcClientConn.Close()
		c.rpcClientConn = nil
	}
}

// ForwardRequest forwards a given request to the active node and returns the
// response.
func (c *Core) ForwardRequest(req *http.Request) (int, http.Header, []byte, error) {
	c.requestForwardingConnectionLock.RLock()
	defer c.requestForwardingConnectionLock.RUnlock()

	switch os.Getenv("VAULT_USE_GRPC_REQUEST_FORWARDING") {
	case "":
		if c.requestForwardingConnection == nil {
			return 0, nil, nil, ErrCannotForward
		}

		if c.requestForwardingConnection.clusterAddr == "" {
			return 0, nil, nil, ErrCannotForward
		}

		freq, err := forwarding.GenerateForwardedHTTPRequest(req, c.requestForwardingConnection.clusterAddr+"/cluster/local/forwarded-request")
		if err != nil {
			c.logger.Error("core/ForwardRequest: error creating forwarded request", "error", err)
			return 0, nil, nil, fmt.Errorf("error creating forwarding request")
		}

		//resp, err := c.requestForwardingConnection.Do(freq)
		resp, err := c.requestForwardingConnection.transport.RoundTrip(freq)
		if err != nil {
			return 0, nil, nil, err
		}
		defer resp.Body.Close()

		// Read the body into a buffer so we can write it back out to the
		// original requestor
		buf := bytes.NewBuffer(nil)
		_, err = buf.ReadFrom(resp.Body)
		if err != nil {
			return 0, nil, nil, err
		}
		return resp.StatusCode, resp.Header, buf.Bytes(), nil

	default:
		if c.rpcForwardingClient == nil {
			return 0, nil, nil, ErrCannotForward
		}

		freq, err := forwarding.GenerateForwardedRequest(req)
		if err != nil {
			c.logger.Error("core/ForwardRequest: error creating forwarding RPC request", "error", err)
			return 0, nil, nil, fmt.Errorf("error creating forwarding RPC request")
		}
		if freq == nil {
			c.logger.Error("core/ForwardRequest: got nil forwarding RPC request")
			return 0, nil, nil, fmt.Errorf("got nil forwarding RPC request")
		}
		resp, err := c.rpcForwardingClient.HandleRequest(context.Background(), freq, grpc.FailFast(true))
		if err != nil {
			c.logger.Error("core/ForwardRequest: error during forwarded RPC request", "error", err)
			return 0, nil, nil, fmt.Errorf("error during forwarding RPC request")
		}

		var header http.Header
		if resp.HeaderEntries != nil {
			header = make(http.Header)
			for k, v := range resp.HeaderEntries {
				for _, j := range v.Values {
					header.Add(k, j)
				}
			}
		}

		return int(resp.StatusCode), header, resp.Body, nil
	}
}

// getGRPCDialer is used to return a dialer that has the correct TLS
// configuration. Otherwise gRPC tries to be helpful and stomps all over our
// NextProtos.
func (c *Core) getGRPCDialer() func(string, time.Duration) (net.Conn, error) {
	return func(addr string, timeout time.Duration) (net.Conn, error) {
		tlsConfig, err := c.ClusterTLSConfig()
		if err != nil {
			c.logger.Error("core/getGRPCDialer: failed to get tls configuration", "error", err)
			return nil, err
		}
		tlsConfig.NextProtos = []string{"req_fw_sb-act_v1"}
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

func (s *forwardedRequestRPCServer) HandleRequest(ctx context.Context, freq *forwarding.Request) (*forwarding.Response, error) {
	// Parse an http.Request out of it
	req, err := forwarding.ParseForwardedRequest(freq)
	if err != nil {
		return nil, err
	}

	// A very dummy response writer that doesn't follow normal semantics, just
	// lets you write a status code (last written wins) and a body. But it
	// meets the interface requirements.
	w := forwarding.NewRPCResponseWriter()

	s.handler.ServeHTTP(w, req)

	resp := &forwarding.Response{
		StatusCode: uint32(w.StatusCode()),
		Body:       w.Body().Bytes(),
	}

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

package vault

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/hashicorp/vault/helper/forwarding"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
)

// Starts the listeners and servers necessary to handle forwarded requests
func (c *Core) startForwarding(lns []net.Listener) error {
	baseHandler, wrappedHandler := c.clusterHandlerSetupFunc()

	tlsConfig, err := c.ClusterTLSConfig()
	if err != nil {
		c.logger.Printf("[ERR] core/startClusterListener: failed to get tls configuration: %v", err)
		return err
	}
	tlsConfig.NextProtos = []string{"h2", "req_fw_sb-act_v1"}

	c.rpcServer = grpc.NewServer()
	RegisterForwardedRequestHandlerServer(c.rpcServer, &forwardedRequestRPCServer{
		core:    c,
		handler: baseHandler,
	})

	tlsLns := make([]net.Listener, 0, len(lns))
	for _, ln := range lns {
		tlsLn := tls.NewListener(ln, tlsConfig)
		tlsLns = append(tlsLns, tlsLn)
		c.logger.Printf("[TRACE] core/startClusterListener: serving cluster requests on %s", tlsLn.Addr())

		fws := &http2.Server{}

		go func() {
			for {
				select {
				case <-c.clusterListenerShutdownCh:
					return
				default:
					conn, err := tlsLn.Accept()
					if err != nil {
						if conn != nil {
							conn.Close()
						}
						continue
					}
					tlsConn := conn.(*tls.Conn)
					err = tlsConn.Handshake()
					if err != nil {
						c.logger.Printf("[TRACE] core/startClusterListener/Accept: error handshaking: %v", err)
						if conn != nil {
							conn.Close()
						}
						continue
					}
					switch tlsConn.ConnectionState().NegotiatedProtocol {
					case "h2":
						c.logger.Printf("[TRACE] core/startClusterListener/Accept: got h2 connection")
						go fws.ServeConn(conn, &http2.ServeConnOpts{
							Handler: wrappedHandler,
						})

					case "req_fw_sb-act_v1":
						c.logger.Printf("[TRACE] core/startClusterListener/Accept: got req_fw_sb-act_v1 connection")
						go fws.ServeConn(conn, &http2.ServeConnOpts{
							Handler: c.rpcServer,
						})

					default:
						c.logger.Printf("[TRACE] core/startClusterListener/Accept: unknown negotiated protocol")
						conn.Close()
						continue
					}
				}
			}
		}()
	}

	c.clusterListenerShutdownCh = make(chan struct{})
	c.clusterListenerShutdownSuccessCh = make(chan struct{})

	go func() {
		<-c.clusterListenerShutdownCh
		c.logger.Printf("[TRACE] core/startClusterListener: shutting down listeners")
		c.rpcServer.Stop()
		c.rpcServer = nil
		for _, tlsLn := range tlsLns {
			tlsLn.Close()
		}
		close(c.clusterListenerShutdownSuccessCh)
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

	// Disabled, potentially
	if clusterAddr == "" {
		c.requestForwardingConnection = nil
		c.forwardingClient = nil
		return nil
	}

	clusterURL, err := url.Parse(clusterAddr)
	if err != nil {
		c.logger.Printf("[ERR] core/refreshRequestForwardingConnection: error parsing cluster address: %v", err)
		return err
	}

	switch os.Getenv("VAULT_USE_GRPC_REQUEST_FORWARDING") {
	case "":
		// Set up normal HTTP forwarding handling
		tlsConfig, err := c.ClusterTLSConfig()
		if err != nil {
			c.logger.Printf("[ERR] core/refreshRequestForwardingConnection: error fetching cluster tls configuration: %v", err)
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
		cc, err := grpc.Dial(clusterURL.Host, grpc.WithDialer(c.getGRPCDialer()), grpc.WithInsecure())
		if err != nil {
			c.logger.Printf("[ERR] core/refreshRequestForwardingConnection: err setting up rpc client: %v", err)
			return err
		}
		c.forwardingClient = NewForwardedRequestHandlerClient(cc)
	}

	return nil
}

// ForwardRequest forwards a given request to the active node and returns the
// response.
func (c *Core) ForwardRequest(req *http.Request) (int, []byte, error) {
	c.requestForwardingConnectionLock.RLock()
	defer c.requestForwardingConnectionLock.RUnlock()

	switch os.Getenv("VAULT_USE_GRPC_REQUEST_FORWARDING") {
	case "":
		if c.requestForwardingConnection == nil {
			return 0, nil, ErrCannotForward
		}

		if c.requestForwardingConnection.clusterAddr == "" {
			return 0, nil, ErrCannotForward
		}

		freq, err := forwarding.GenerateForwardedHTTPRequest(req, c.requestForwardingConnection.clusterAddr+"/cluster/local/forwarded-request")
		if err != nil {
			c.logger.Printf("[ERR] core/ForwardRequest: error creating forwarded request: %v", err)
			return 0, nil, fmt.Errorf("error creating forwarding request")
		}

		//resp, err := c.requestForwardingConnection.Do(freq)
		resp, err := c.requestForwardingConnection.transport.RoundTrip(freq)
		if err != nil {
			return 0, nil, err
		}
		defer resp.Body.Close()

		// Read the body into a buffer so we can write it back out to the
		// original requestor
		buf := bytes.NewBuffer(nil)
		_, err = buf.ReadFrom(resp.Body)
		if err != nil {
			return 0, nil, err
		}
		return resp.StatusCode, buf.Bytes(), nil

	default:
		if c.forwardingClient == nil {
			return 0, nil, ErrCannotForward
		}

		freq, err := forwarding.GenerateForwardedRequest(req)
		if err != nil {
			c.logger.Printf("[ERR] core/ForwardRequest: error creating forwarding RPC request: %v", err)
			return 0, nil, fmt.Errorf("error creating forwarding RPC request")
		}
		if freq == nil {
			c.logger.Printf("[ERR] core/ForwardRequest: got nil forwarding RPC request")
			return 0, nil, fmt.Errorf("got nil forwarding RPC request")
		}
		resp, err := c.forwardingClient.HandleRequest(context.Background(), freq, grpc.FailFast(true))
		if err != nil {
			c.logger.Printf("[ERR] core/ForwardRequest: error during forwarded RPC request: %v", err)
			return 0, nil, fmt.Errorf("error during forwarding RPC request")
		}
		return int(resp.StatusCode), resp.Body, nil
	}
}

func (c *Core) getGRPCDialer() func(string, time.Duration) (net.Conn, error) {
	return func(addr string, timeout time.Duration) (net.Conn, error) {
		tlsConfig, err := c.ClusterTLSConfig()
		if err != nil {
			c.logger.Printf("[ERR] core/getGRPCDialer: failed to get tls configuration: %v", err)
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
	req, err := forwarding.ParseForwardedRequest(freq)
	if err != nil {
		return nil, err
	}

	w := forwarding.NewRPCResponseWriter()
	s.handler.ServeHTTP(w, req)

	return &forwarding.Response{
		StatusCode: uint32(w.StatusCode()),
		Body:       w.Body().Bytes(),
	}, nil
}

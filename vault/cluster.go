package vault

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/logical"
	"golang.org/x/net/http2"
)

const (
	// Storage path where the local cluster name and identifier are stored
	coreLocalClusterInfoPath = "core/cluster/local/info"

	corePrivateKeyTypeP521    = "p521"
	corePrivateKeyTypeED25519 = "ed25519"

	// Internal so as not to log a trace message
	IntNoForwardingHeaderName = "X-Vault-Internal-No-Request-Forwarding"
)

var (
	ErrCannotForward = errors.New("cannot forward request; no connection or address not known")
)

type ClusterLeaderParams struct {
	LeaderUUID         string
	LeaderRedirectAddr string
	LeaderClusterAddr  string
}

// Structure representing the storage entry that holds cluster information
type Cluster struct {
	// Name of the cluster
	Name string `json:"name" structs:"name" mapstructure:"name"`

	// Identifier of the cluster
	ID string `json:"id" structs:"id" mapstructure:"id"`
}

// Cluster fetches the details of the local cluster. This method errors out
// when Vault is sealed.
func (c *Core) Cluster(ctx context.Context) (*Cluster, error) {
	var cluster Cluster

	// Fetch the storage entry. This call fails when Vault is sealed.
	entry, err := c.barrier.Get(ctx, coreLocalClusterInfoPath)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return &cluster, nil
	}

	// Decode the cluster information
	if err = jsonutil.DecodeJSON(entry.Value, &cluster); err != nil {
		return nil, errwrap.Wrapf("failed to decode cluster details: {{err}}", err)
	}

	// Set in config file
	if c.clusterName != "" {
		cluster.Name = c.clusterName
	}

	return &cluster, nil
}

// This sets our local cluster cert and private key based on the advertisement.
// It also ensures the cert is in our local cluster cert pool.
func (c *Core) loadLocalClusterTLS(adv activeAdvertisement) (retErr error) {
	defer func() {
		if retErr != nil {
			c.localClusterCert.Store(([]byte)(nil))
			c.localClusterParsedCert.Store((*x509.Certificate)(nil))
			c.localClusterPrivateKey.Store((*ecdsa.PrivateKey)(nil))

			c.requestForwardingConnectionLock.Lock()
			c.clearForwardingClients()
			c.requestForwardingConnectionLock.Unlock()
		}
	}()

	switch {
	case adv.ClusterAddr == "":
		// Clustering disabled on the server, don't try to look for params
		return nil

	case adv.ClusterKeyParams == nil:
		c.logger.Error("no key params found loading local cluster TLS information")
		return fmt.Errorf("no local cluster key params found")

	case adv.ClusterKeyParams.X == nil, adv.ClusterKeyParams.Y == nil, adv.ClusterKeyParams.D == nil:
		c.logger.Error("failed to parse local cluster key due to missing params")
		return fmt.Errorf("failed to parse local cluster key")

	case adv.ClusterKeyParams.Type != corePrivateKeyTypeP521:
		c.logger.Error("unknown local cluster key type", "key_type", adv.ClusterKeyParams.Type)
		return fmt.Errorf("failed to find valid local cluster key type")

	case adv.ClusterCert == nil || len(adv.ClusterCert) == 0:
		c.logger.Error("no local cluster cert found")
		return fmt.Errorf("no local cluster cert found")

	}

	c.localClusterPrivateKey.Store(&ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P521(),
			X:     adv.ClusterKeyParams.X,
			Y:     adv.ClusterKeyParams.Y,
		},
		D: adv.ClusterKeyParams.D,
	})

	locCert := make([]byte, len(adv.ClusterCert))
	copy(locCert, adv.ClusterCert)
	c.localClusterCert.Store(locCert)

	cert, err := x509.ParseCertificate(adv.ClusterCert)
	if err != nil {
		c.logger.Error("failed parsing local cluster certificate", "error", err)
		return errwrap.Wrapf("error parsing local cluster certificate: {{err}}", err)
	}

	c.localClusterParsedCert.Store(cert)

	return nil
}

// setupCluster creates storage entries for holding Vault cluster information.
// Entries will be created only if they are not already present. If clusterName
// is not supplied, this method will auto-generate it.
func (c *Core) setupCluster(ctx context.Context) error {
	// Prevent data races with the TLS parameters
	c.clusterParamsLock.Lock()
	defer c.clusterParamsLock.Unlock()

	// Check if storage index is already present or not
	cluster, err := c.Cluster(ctx)
	if err != nil {
		c.logger.Error("failed to get cluster details", "error", err)
		return err
	}

	var modified bool

	if cluster == nil {
		cluster = &Cluster{}
	}

	if cluster.Name == "" {
		// If cluster name is not supplied, generate one
		if c.clusterName == "" {
			c.logger.Debug("cluster name not found/set, generating new")
			clusterNameBytes, err := uuid.GenerateRandomBytes(4)
			if err != nil {
				c.logger.Error("failed to generate cluster name", "error", err)
				return err
			}

			c.clusterName = fmt.Sprintf("vault-cluster-%08x", clusterNameBytes)
		}

		cluster.Name = c.clusterName
		if c.logger.IsDebug() {
			c.logger.Debug("cluster name set", "name", cluster.Name)
		}
		modified = true
	}

	if cluster.ID == "" {
		c.logger.Debug("cluster ID not found, generating new")
		// Generate a clusterID
		cluster.ID, err = uuid.GenerateUUID()
		if err != nil {
			c.logger.Error("failed to generate cluster identifier", "error", err)
			return err
		}
		if c.logger.IsDebug() {
			c.logger.Debug("cluster ID set", "id", cluster.ID)
		}
		modified = true
	}

	// If we're using HA, generate server-to-server parameters
	if c.ha != nil {
		// Create a private key
		if c.localClusterPrivateKey.Load().(*ecdsa.PrivateKey) == nil {
			c.logger.Debug("generating cluster private key")
			key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
			if err != nil {
				c.logger.Error("failed to generate local cluster key", "error", err)
				return err
			}

			c.localClusterPrivateKey.Store(key)
		}

		// Create a certificate
		if c.localClusterCert.Load().([]byte) == nil {
			c.logger.Debug("generating local cluster certificate")

			host, err := uuid.GenerateUUID()
			if err != nil {
				return err
			}
			host = fmt.Sprintf("fw-%s", host)
			template := &x509.Certificate{
				Subject: pkix.Name{
					CommonName: host,
				},
				DNSNames: []string{host},
				ExtKeyUsage: []x509.ExtKeyUsage{
					x509.ExtKeyUsageServerAuth,
					x509.ExtKeyUsageClientAuth,
				},
				KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement | x509.KeyUsageCertSign,
				SerialNumber: big.NewInt(mathrand.Int63()),
				NotBefore:    time.Now().Add(-30 * time.Second),
				// 30 years of single-active uptime ought to be enough for anybody
				NotAfter:              time.Now().Add(262980 * time.Hour),
				BasicConstraintsValid: true,
				IsCA:                  true,
			}

			certBytes, err := x509.CreateCertificate(rand.Reader, template, template, c.localClusterPrivateKey.Load().(*ecdsa.PrivateKey).Public(), c.localClusterPrivateKey.Load().(*ecdsa.PrivateKey))
			if err != nil {
				c.logger.Error("error generating self-signed cert", "error", err)
				return errwrap.Wrapf("unable to generate local cluster certificate: {{err}}", err)
			}

			parsedCert, err := x509.ParseCertificate(certBytes)
			if err != nil {
				c.logger.Error("error parsing self-signed cert", "error", err)
				return errwrap.Wrapf("error parsing generated certificate: {{err}}", err)
			}

			c.localClusterCert.Store(certBytes)
			c.localClusterParsedCert.Store(parsedCert)
		}
	}

	if modified {
		// Encode the cluster information into as a JSON string
		rawCluster, err := json.Marshal(cluster)
		if err != nil {
			c.logger.Error("failed to encode cluster details", "error", err)
			return err
		}

		// Store it
		err = c.barrier.Put(ctx, &logical.StorageEntry{
			Key:   coreLocalClusterInfoPath,
			Value: rawCluster,
		})
		if err != nil {
			c.logger.Error("failed to store cluster details", "error", err)
			return err
		}
	}

	return nil
}

// ClusterClient is used to lookup a client certificate.
type ClusterClient interface {
	ClientLookup(context.Context, *tls.CertificateRequestInfo) (*tls.Certificate, error)
}

// ClusterHandler exposes functions for looking up TLS configuration and handing
// off a connection for a cluster listener application.
type ClusterHandler interface {
	ServerLookup(context.Context, *tls.ClientHelloInfo) (*tls.Certificate, error)
	CALookup(context.Context) (*x509.Certificate, error)

	// Handoff is used to pass the connection lifetime off to
	// the handler
	Handoff(context.Context, *sync.WaitGroup, chan struct{}, *tls.Conn) error
	Stop() error
}

// ClusterListener is the source of truth for cluster handlers and connection
// clients. It dynamically builds the cluster TLS information. It's also
// responsible for starting tcp listeners and accepting new cluster connections.
type ClusterListener struct {
	handlers   map[string]ClusterHandler
	clients    map[string]ClusterClient
	shutdown   *uint32
	shutdownWg *sync.WaitGroup
	server     *http2.Server

	clusterListenerAddrs []*net.TCPAddr
	clusterCipherSuites  []uint16
	logger               log.Logger
	l                    sync.RWMutex
}

// AddClient adds a new client for an ALPN name
func (cl *ClusterListener) AddClient(alpn string, client ClusterClient) {
	cl.l.Lock()
	cl.clients[alpn] = client
	cl.l.Unlock()
}

// RemoveClient removes the client for the specified ALPN name
func (cl *ClusterListener) RemoveClient(alpn string) {
	cl.l.Lock()
	delete(cl.clients, alpn)
	cl.l.Unlock()
}

// AddHandler registers a new cluster handler for the provided ALPN name.
func (cl *ClusterListener) AddHandler(alpn string, handler ClusterHandler) {
	cl.l.Lock()
	cl.handlers[alpn] = handler
	cl.l.Unlock()
}

// StopHandler stops the cluster handler for the provided ALPN name, it also
// calls stop on the handler.
func (cl *ClusterListener) StopHandler(alpn string) {
	cl.l.Lock()
	handler, ok := cl.handlers[alpn]
	delete(cl.handlers, alpn)
	cl.l.Unlock()
	if ok {
		handler.Stop()
	}
}

// Server returns the http2 server that the cluster listener is using
func (cl *ClusterListener) Server() *http2.Server {
	return cl.server
}

// TLSConfig returns a tls config object that uses dynamic lookups to correctly
// authenticate registered handlers/clients
func (cl *ClusterListener) TLSConfig(ctx context.Context) (*tls.Config, error) {
	serverLookup := func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		cl.logger.Debug("performing server cert lookup")

		cl.l.RLock()
		defer cl.l.RUnlock()
		for _, v := range clientHello.SupportedProtos {
			if handler, ok := cl.handlers[v]; ok {
				return handler.ServerLookup(ctx, clientHello)
			}
		}

		cl.logger.Warn("no TLS certs found for ALPN", "ALPN", clientHello.SupportedProtos)
		return nil, errors.New("unsupported protocol")
	}

	clientLookup := func(requestInfo *tls.CertificateRequestInfo) (*tls.Certificate, error) {
		cl.logger.Debug("performing client cert lookup")

		cl.l.RLock()
		defer cl.l.RUnlock()
		for _, client := range cl.clients {
			cert, err := client.ClientLookup(ctx, requestInfo)
			if err == nil && cert != nil {
				return cert, nil
			}
		}

		cl.logger.Warn("no client information found")
		return nil, errors.New("no client cert found")
	}

	serverConfigLookup := func(clientHello *tls.ClientHelloInfo) (*tls.Config, error) {
		caPool := x509.NewCertPool()

		ret := &tls.Config{
			ClientAuth:           tls.RequireAndVerifyClientCert,
			GetCertificate:       serverLookup,
			GetClientCertificate: clientLookup,
			MinVersion:           tls.VersionTLS12,
			RootCAs:              caPool,
			ClientCAs:            caPool,
			NextProtos:           clientHello.SupportedProtos,
			CipherSuites:         cl.clusterCipherSuites,
		}

		cl.l.RLock()
		defer cl.l.RUnlock()
		for _, v := range clientHello.SupportedProtos {
			if handler, ok := cl.handlers[v]; ok {
				ca, err := handler.CALookup(ctx)
				if err != nil {
					return nil, err
				}

				caPool.AddCert(ca)
				return ret, nil
			}
		}

		cl.logger.Warn("no TLS config found for ALPN", "ALPN", clientHello.SupportedProtos)
		return nil, errors.New("unsupported protocol")
	}

	return &tls.Config{
		ClientAuth:           tls.RequireAndVerifyClientCert,
		GetCertificate:       serverLookup,
		GetClientCertificate: clientLookup,
		GetConfigForClient:   serverConfigLookup,
		MinVersion:           tls.VersionTLS12,
		CipherSuites:         cl.clusterCipherSuites,
	}, nil
}

// Run starts the tcp listeners and will accept connections until stop is
// called.
func (cl *ClusterListener) Run(ctx context.Context) error {
	// Get our TLS config
	tlsConfig, err := cl.TLSConfig(ctx)
	if err != nil {
		cl.logger.Error("failed to get tls configuration when starting cluster listener", "error", err)
		return err
	}

	// The server supports all of the possible protos
	tlsConfig.NextProtos = []string{"h2", requestForwardingALPN, perfStandbyALPN, PerformanceReplicationALPN, DRReplicationALPN}

	for i, laddr := range cl.clusterListenerAddrs {
		// closeCh is used to shutdown the spawned goroutines once this
		// function returns
		closeCh := make(chan struct{})

		if cl.logger.IsInfo() {
			cl.logger.Info("starting listener", "listener_address", laddr)
		}

		// Create a TCP listener. We do this separately and specifically
		// with TCP so that we can set deadlines.
		tcpLn, err := net.ListenTCP("tcp", laddr)
		if err != nil {
			cl.logger.Error("error starting listener", "error", err)
			continue
		}
		if laddr.String() != tcpLn.Addr().String() {
			// If we listened on port 0, record the port the OS gave us.
			cl.clusterListenerAddrs[i] = tcpLn.Addr().(*net.TCPAddr)
		}

		// Wrap the listener with TLS
		tlsLn := tls.NewListener(tcpLn, tlsConfig)

		if cl.logger.IsInfo() {
			cl.logger.Info("serving cluster requests", "cluster_listen_address", tlsLn.Addr())
		}

		cl.shutdownWg.Add(1)
		// Start our listening loop
		go func(closeCh chan struct{}, tlsLn net.Listener) {
			defer func() {
				cl.shutdownWg.Done()
				tlsLn.Close()
				close(closeCh)
			}()

			for {
				if atomic.LoadUint32(cl.shutdown) > 0 {
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
						cl.logger.Debug("non-timeout error accepting on cluster port", "error", err)
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
					if cl.logger.IsDebug() {
						cl.logger.Debug("error setting deadline for cluster connection", "error", err)
					}
					tlsConn.Close()
					continue
				}

				err = tlsConn.Handshake()
				if err != nil {
					if cl.logger.IsDebug() {
						cl.logger.Debug("error handshaking cluster connection", "error", err)
					}
					tlsConn.Close()
					continue
				}

				// Now, set it back to unlimited
				err = tlsConn.SetDeadline(time.Time{})
				if err != nil {
					if cl.logger.IsDebug() {
						cl.logger.Debug("error setting deadline for cluster connection", "error", err)
					}
					tlsConn.Close()
					continue
				}

				cl.l.RLock()
				handler, ok := cl.handlers[tlsConn.ConnectionState().NegotiatedProtocol]
				cl.l.RUnlock()
				if !ok {
					cl.logger.Debug("unknown negotiated protocol on cluster port")
					tlsConn.Close()
					continue
				}

				if err := handler.Handoff(ctx, cl.shutdownWg, closeCh, tlsConn); err != nil {
					cl.logger.Error("error handling cluster connection", "error", err)
					continue
				}
			}
		}(closeCh, tlsLn)
	}

	return nil
}

// Stop stops the cluster listner
func (cl *ClusterListener) Stop() {
	// Set the shutdown flag. This will cause the listeners to shut down
	// within the deadline in clusterListenerAcceptDeadline
	atomic.StoreUint32(cl.shutdown, 1)
	cl.logger.Info("forwarding rpc listeners stopped")

	// Wait for them all to shut down
	cl.shutdownWg.Wait()
	cl.logger.Info("rpc listeners successfully shut down")
}

// startClusterListener starts cluster request listeners during unseal. It
// is assumed that the state lock is held while this is run. Right now this
// only starts cluster listeners. Once the listener is started handlers/clients
// can start being registered to it.
func (c *Core) startClusterListener(ctx context.Context) error {
	if c.clusterAddr == "" {
		c.logger.Info("clustering disabled, not starting listeners")
		return nil
	}

	if c.clusterListenerAddrs == nil || len(c.clusterListenerAddrs) == 0 {
		c.logger.Warn("clustering not disabled but no addresses to listen on")
		return fmt.Errorf("cluster addresses not found")
	}

	c.logger.Debug("starting cluster listeners")

	// Create the HTTP/2 server that will be shared by both RPC and regular
	// duties. Doing it this way instead of listening via the server and gRPC
	// allows us to re-use the same port via ALPN. We can just tell the server
	// to serve a given conn and which handler to use.
	h2Server := &http2.Server{
		// Our forwarding connections heartbeat regularly so anything else we
		// want to go away/get cleaned up pretty rapidly
		IdleTimeout: 5 * HeartbeatInterval,
	}

	c.clusterListener = &ClusterListener{
		handlers:   make(map[string]ClusterHandler),
		clients:    make(map[string]ClusterClient),
		shutdown:   new(uint32),
		shutdownWg: &sync.WaitGroup{},
		server:     h2Server,

		clusterListenerAddrs: c.clusterListenerAddrs,
		clusterCipherSuites:  c.clusterCipherSuites,
		logger:               c.logger.Named("cluster-listener"),
	}

	err := c.clusterListener.Run(ctx)
	if err != nil {
		return err
	}
	if strings.HasSuffix(c.clusterAddr, ":0") {
		// If we listened on port 0, record the port the OS gave us.
		c.clusterAddr = fmt.Sprintf("https://%s", c.clusterListener.clusterListenerAddrs[0])
	}
	return nil
}

// stopClusterListener stops any existing listeners during seal. It is
// assumed that the state lock is held while this is run.
func (c *Core) stopClusterListener() {
	if c.clusterListener == nil {
		c.logger.Debug("clustering disabled, not stopping listeners")
		return
	}

	c.logger.Info("stopping cluster listeners")

	c.clusterListener.Stop()

	c.logger.Info("cluster listeners successfully shut down")
}

func (c *Core) SetClusterListenerAddrs(addrs []*net.TCPAddr) {
	c.clusterListenerAddrs = addrs
	if c.clusterAddr == "" && len(addrs) == 1 {
		c.clusterAddr = fmt.Sprintf("https://%s", addrs[0].String())
	}
}

func (c *Core) SetClusterHandler(handler http.Handler) {
	c.clusterHandler = handler
}

// getGRPCDialer is used to return a dialer that has the correct TLS
// configuration. Otherwise gRPC tries to be helpful and stomps all over our
// NextProtos.
func (c *Core) getGRPCDialer(ctx context.Context, alpnProto, serverName string, caCert *x509.Certificate) func(string, time.Duration) (net.Conn, error) {
	return func(addr string, timeout time.Duration) (net.Conn, error) {
		if c.clusterListener == nil {
			return nil, errors.New("clustering disabled")
		}

		tlsConfig, err := c.clusterListener.TLSConfig(ctx)
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

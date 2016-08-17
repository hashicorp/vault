package vault

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	mathrand "math/rand"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"

	"golang.org/x/net/context"
	"golang.org/x/net/http2"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/forwarding"
	"github.com/hashicorp/vault/helper/jsonutil"
)

const (
	// Storage path where the local cluster name and identifier are stored
	coreLocalClusterInfoPath = "core/cluster/local/info"

	corePrivateKeyTypeP521 = "p521"

	// Internal so as not to log a trace message
	IntNoForwardingHeaderName = "X-Vault-Internal-No-Request-Forwarding"
)

var (
	ErrCannotForward = errors.New("cannot forward request; no connection or address not known")
)

type clusterKeyParams struct {
	Type string   `json:"type"`
	X    *big.Int `json:"x"`
	Y    *big.Int `json:"y"`
	D    *big.Int `json:"d"`
}

type activeConnection struct {
	*http.Client
	clusterAddr string
}

// Structure representing the storage entry that holds cluster information
type Cluster struct {
	// Name of the cluster
	Name string `json:"name" structs:"name" mapstructure:"name"`

	// Identifier of the cluster
	ID string `json:"id" structs:"id" mapstructure:"id"`
}

// Cluster fetches the details of either local or global cluster based on the
// input. This method errors out when Vault is sealed.
func (c *Core) Cluster() (*Cluster, error) {
	var cluster Cluster

	// Fetch the storage entry. This call fails when Vault is sealed.
	entry, err := c.barrier.Get(coreLocalClusterInfoPath)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return &cluster, nil
	}

	// Decode the cluster information
	if err = jsonutil.DecodeJSON(entry.Value, &cluster); err != nil {
		return nil, fmt.Errorf("failed to decode cluster details: %v", err)
	}

	// Set in config file
	if c.clusterName != "" {
		cluster.Name = c.clusterName
	}

	return &cluster, nil
}

// This is idempotent, so we return nil if there is no entry yet (say, because
// the active node has not yet generated this)
func (c *Core) loadClusterTLS(adv activeAdvertisement) error {
	c.clusterParamsLock.Lock()
	defer c.clusterParamsLock.Unlock()

	switch {
	case adv.ClusterKeyParams.X == nil, adv.ClusterKeyParams.Y == nil, adv.ClusterKeyParams.D == nil:
		c.logger.Printf("[ERR] core/loadClusterPrivateKey: failed to parse local cluster key due to missing params")
		return fmt.Errorf("failed to parse local cluster key")
	case adv.ClusterKeyParams.Type == corePrivateKeyTypeP521:
	default:
		c.logger.Printf("[ERR] core/loadClusterPrivateKey: unknown local cluster key type %v", adv.ClusterKeyParams.Type)
		return fmt.Errorf("failed to find valid local cluster key type")
	}
	c.localClusterPrivateKey = &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P521(),
			X:     adv.ClusterKeyParams.X,
			Y:     adv.ClusterKeyParams.Y,
		},
		D: adv.ClusterKeyParams.D,
	}

	c.localClusterCert = adv.ClusterCert

	cert, err := x509.ParseCertificate(c.localClusterCert)
	if err != nil {
		c.logger.Printf("[ERR] core/loadClusterPrivateKey: failed parsing local cluster certificate: %v", err)
		return fmt.Errorf("error parsing local cluster certificate: %v", err)
	}

	c.localClusterCertPool.AddCert(cert)

	return nil
}

// setupCluster creates storage entries for holding Vault cluster information.
// Entries will be created only if they are not already present. If clusterName
// is not supplied, this method will auto-generate it.
func (c *Core) setupCluster() error {
	c.clusterParamsLock.Lock()
	defer c.clusterParamsLock.Unlock()

	// Check if storage index is already present or not
	cluster, err := c.Cluster()
	if err != nil {
		c.logger.Printf("[ERR] core: failed to get cluster details: %v", err)
		return err
	}

	var modified bool

	if cluster == nil {
		cluster = &Cluster{}
	}

	if cluster.Name == "" {
		// If cluster name is not supplied, generate one
		if c.clusterName == "" {
			c.logger.Printf("[TRACE] core: cluster name not found/set, generating new")
			clusterNameBytes, err := uuid.GenerateRandomBytes(4)
			if err != nil {
				c.logger.Printf("[ERR] core: failed to generate cluster name: %v", err)
				return err
			}

			c.clusterName = fmt.Sprintf("vault-cluster-%08x", clusterNameBytes)
		}

		cluster.Name = c.clusterName
		c.logger.Printf("[DEBUG] core: cluster name set to %s", cluster.Name)
		modified = true
	}

	if cluster.ID == "" {
		c.logger.Printf("[TRACE] core: cluster ID not found, generating new")
		// Generate a clusterID
		cluster.ID, err = uuid.GenerateUUID()
		if err != nil {
			c.logger.Printf("[ERR] core: failed to generate cluster identifier: %v", err)
			return err
		}
		c.logger.Printf("[DEBUG] core: cluster ID set to %s", cluster.ID)
		modified = true
	}

	// Create a private key
	{
		c.logger.Printf("[TRACE] core: generating cluster private key")
		key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		if err != nil {
			c.logger.Printf("[ERR] core: failed to generate local cluster key: %v", err)
			return err
		}

		c.localClusterPrivateKey = key
	}

	// Create a certificate
	{
		c.logger.Printf("[TRACE] core: generating local cluster certificate")

		host, err := uuid.GenerateUUID()
		if err != nil {
			return err
		}

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
			IsCA: true,
		}

		certBytes, err := x509.CreateCertificate(rand.Reader, template, template, c.localClusterPrivateKey.Public(), c.localClusterPrivateKey)
		if err != nil {
			c.logger.Printf("[ERR] core: error generating self-signed cert: %v", err)
			return fmt.Errorf("unable to generate local cluster certificate: %v", err)
		}

		_, err = x509.ParseCertificate(certBytes)
		if err != nil {
			c.logger.Printf("[ERR] core: error parsing self-signed cert: %v", err)
			return fmt.Errorf("error parsing generated certificate: %v", err)
		}

		c.localClusterCert = certBytes
	}

	if modified {
		// Encode the cluster information into as a JSON string
		rawCluster, err := json.Marshal(cluster)
		if err != nil {
			c.logger.Printf("[ERR] core: failed to encode cluster details: %v", err)
			return err
		}

		// Store it
		err = c.barrier.Put(&Entry{
			Key:   coreLocalClusterInfoPath,
			Value: rawCluster,
		})
		if err != nil {
			c.logger.Printf("[ERR] core: failed to store cluster details: %v", err)
			return err
		}
	}

	return nil
}

// SetClusterListenerSetupFunc sets the listener setup func, which is used to
// know which ports to listen on and a handler to use.
func (c *Core) SetClusterListenerSetupFunc(setupFunc func() ([]net.Listener, http.Handler, error)) {
	c.clusterListenerSetupFunc = setupFunc
}

// startClusterListener starts cluster request listeners during postunseal. It
// is assumed that the state lock is held while this is run.
func (c *Core) startClusterListener() error {
	if c.clusterListenerShutdownCh != nil {
		c.logger.Printf("[ERR] core/startClusterListener: attempt to set up cluster listeners when already set up")
		return fmt.Errorf("cluster listeners already setup")
	}

	if c.clusterListenerSetupFunc == nil {
		c.logger.Printf("[ERR] core/startClusterListener: cluster listener setup function has not been set")
		return fmt.Errorf("cluster listener setup function has not been set")
	}

	if c.clusterAddr == "" {
		c.logger.Printf("[TRACE] core/startClusterListener: clustering disabled, starting listeners")
		return nil
	}

	c.logger.Printf("[TRACE] core/startClusterListener: starting listeners")

	lns, handler, err := c.clusterListenerSetupFunc()
	if err != nil {
		return err
	}

	tlsConfig, err := c.ClusterTLSConfig()
	if err != nil {
		c.logger.Printf("[ERR] core/startClusterListener: failed to get tls configuration: %v", err)
		return err
	}

	tlsLns := make([]net.Listener, 0, len(lns))
	for _, ln := range lns {
		tlsLn := tls.NewListener(ln, tlsConfig)
		tlsLns = append(tlsLns, tlsLn)
		server := &http.Server{
			Handler: handler,
		}
		http2.ConfigureServer(server, nil)
		server.TLSNextProto["forwarding_v1"] = c.handleRPCForwardingRequest

		c.logger.Printf("[TRACE] core/startClusterListener: serving cluster requests on %s", tlsLn.Addr())

		c.forwardingService = grpc.NewServer()
		RegisterForwardedRequestHandlerServer(c.forwardingService, &forwardedRequestRPCServer{})

		go server.Serve(tlsLn)
	}

	c.clusterListenerShutdownCh = make(chan struct{})
	c.clusterListenerShutdownSuccessCh = make(chan struct{})

	go func() {
		<-c.clusterListenerShutdownCh
		c.logger.Printf("[TRACE] core/startClusterListener: shutting down listeners")
		c.forwardingService.Stop()
		c.forwardingService = nil
		for _, tlsLn := range tlsLns {
			tlsLn.Close()
		}
		close(c.clusterListenerShutdownSuccessCh)
	}()

	return nil
}

// stopClusterListener stops any existing listeners during preseal. It is
// assumed that the state lock is held while this is run.
func (c *Core) stopClusterListener() {
	c.logger.Printf("[TRACE] core/stopClusterListener: stopping listeners")
	if c.clusterListenerShutdownCh != nil {
		close(c.clusterListenerShutdownCh)
		defer func() { c.clusterListenerShutdownCh = nil }()
	}

	// The reason for this loop-de-loop is that we may be unsealing again
	// quickly, and if the listeners are not yet closed, we will get socket
	// bind errors. This ensures proper ordering.
	if c.clusterListenerShutdownSuccessCh == nil {
		return
	}
	<-c.clusterListenerShutdownSuccessCh
	defer func() { c.clusterListenerShutdownSuccessCh = nil }()
}

// ClusterTLSConfig generates a TLS configuration based on the local cluster
// key and cert. This isn't called often and we lock because the CertPool is
// not concurrency-safe.
func (c *Core) ClusterTLSConfig() (*tls.Config, error) {
	c.clusterParamsLock.Lock()
	defer c.clusterParamsLock.Unlock()

	cluster, err := c.Cluster()
	if err != nil {
		return nil, err
	}
	if cluster == nil {
		return nil, fmt.Errorf("cluster information is nil")
	}
	if c.localClusterCert == nil || len(c.localClusterCert) == 0 {
		return nil, fmt.Errorf("cluster certificate is nil")
	}

	parsedCert, err := x509.ParseCertificate(c.localClusterCert)
	if err != nil {
		return nil, fmt.Errorf("error parsing local cluster certificate: %v", err)
	}

	// This is idempotent, so be sure it's been added
	c.localClusterCertPool.AddCert(parsedCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{
			tls.Certificate{
				Certificate: [][]byte{c.localClusterCert},
				PrivateKey:  c.localClusterPrivateKey,
			},
		},
		RootCAs:    c.localClusterCertPool,
		ServerName: parsedCert.Subject.CommonName,
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  c.localClusterCertPool,
		NextProtos: []string{
			"forwarding_v1",
		},
	}

	return tlsConfig, nil
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
		return nil
	}

	tlsConfig, err := c.ClusterTLSConfig()
	if err != nil {
		c.logger.Printf("[ERR] core/refreshRequestForwardingConnection: error fetching cluster tls configuration: %v", err)
		return err
	}
	tp := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	err = http2.ConfigureTransport(tp)
	if err != nil {
		c.logger.Printf("[ERR] core/refreshRequestForwardingConnection: error configuring transport: %v", err)
		return err
	}
	c.requestForwardingConnection = &activeConnection{
		Client: &http.Client{
			Transport: tp,
		},
		clusterAddr: clusterAddr,
	}

	return nil
}

// ForwardRequest forwards a given request to the active node and returns the
// response.
func (c *Core) ForwardRequest(req *http.Request) (*http.Response, error) {
	c.requestForwardingConnectionLock.RLock()
	defer c.requestForwardingConnectionLock.RUnlock()
	if c.requestForwardingConnection == nil {
		return nil, ErrCannotForward
	}

	if c.requestForwardingConnection.clusterAddr == "" {
		return nil, ErrCannotForward
	}

	freq, err := forwarding.GenerateForwardedRequest(req, c.requestForwardingConnection.clusterAddr+"/cluster/local/forwarded-request")
	if err != nil {
		c.logger.Printf("[ERR] core/ForwardRequest: error creating forwarded request: %v", err)
		return nil, fmt.Errorf("error creating forwarding request")
	}

	return c.requestForwardingConnection.Do(freq)
}

// WrapListenersForClustering takes in Vault's listeners and original HTTP
// handler, creates a new handler that handles forwarded requests, and returns
// the cluster setup function that creates the new listners and assigns to the
// new handler
func WrapListenersForClustering(addrs []string, handler http.Handler, logger *log.Logger) func() ([]net.Listener, http.Handler, error) {
	// This mux handles cluster functions (right now, only forwarded requests)
	mux := http.NewServeMux()
	mux.HandleFunc("/cluster/local/forwarded-request", func(w http.ResponseWriter, req *http.Request) {
		freq, err := forwarding.ParseForwardedRequest(req)
		if err != nil {
			if logger != nil {
				logger.Printf("[ERR] http/ForwardedRequestHandler: error parsing forwarded request: %v", err)
			}

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)

			type errorResponse struct {
				Errors []string
			}
			resp := &errorResponse{
				Errors: []string{
					err.Error(),
				},
			}

			enc := json.NewEncoder(w)
			enc.Encode(resp)
			return
		}

		// To avoid the risk of a forward loop in some pathological condition,
		// set the no-forward header
		freq.Header.Set(IntNoForwardingHeaderName, "true")
		handler.ServeHTTP(w, freq)
	})

	return func() ([]net.Listener, http.Handler, error) {
		ret := make([]net.Listener, 0, len(addrs))
		// Loop over the existing listeners and start listeners on appropriate ports
		for _, addr := range addrs {
			ln, err := net.Listen("tcp", addr)
			if err != nil {
				return nil, nil, err
			}
			ret = append(ret, ln)
		}

		return ret, mux, nil
	}
}

type forwardedRequestRPCServer struct{}

func (s *forwardedRequestRPCServer) HandleRequest(context.Context, *forwarding.Request) (*forwarding.Response, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *Core) handleRPCForwardingRequest(server *http.Server, conn *tls.Conn, handler http.Handler) {
	h2s := &http2.Server{}
	h2s.ServeConn(conn, &http2.ServeConnOpts{
		Handler: c.forwardingService,
	})
}

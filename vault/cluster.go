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

	"golang.org/x/net/http2"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/requestutil"
)

const (
	// Storage path where the local cluster name and identifier are stored
	coreLocalClusterInfoPath = "core/cluster/local/info"
	coreLocalClusterKeyPath  = "core/cluster/local/key"

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

	// Certificate corresponding to the private key
	Certificate []byte `json:"certificate" structs:"certificate" mapstructure:"certificate"`
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
	c.clusterParamsLock.RLock()
	if c.localClusterPrivateKey != nil && c.localClusterCert != nil && len(c.localClusterCert) != 0 {
		c.clusterParamsLock.RUnlock()
		return nil
	}

	c.clusterParamsLock.RUnlock()
	c.clusterParamsLock.Lock()
	defer c.clusterParamsLock.Unlock()

	// Verify no modification
	if c.localClusterPrivateKey != nil && c.localClusterCert != nil && len(c.localClusterCert) != 0 {
		return nil
	}

	if c.localClusterPrivateKey == nil {
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
	// Check if storage index is already present or not
	cluster, err := c.Cluster()
	if err != nil {
		c.logger.Printf("[ERR] core: failed to get cluster details: %v", err)
		return err
	}

	var modified bool
	var needNewCert bool

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

	// Check for and optionally generate the local cluster private key
	if c.localClusterPrivateKey == nil {
		c.logger.Printf("[TRACE] core: local cluster private key not loaded, loading now")
		entry, err := c.barrier.Get(coreLocalClusterKeyPath)
		if err != nil {
			c.logger.Printf("[ERR] core: failed to read local cluster key path: %v", err)
			return err
		}
		if entry == nil {
			c.logger.Printf("[TRACE] core: local cluster private key not found, generating new")
			key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
			if err != nil {
				c.logger.Printf("[ERR] core: failed to generate local cluster key: %v", err)
				return err
			}

			// Encode the cluster information into as a JSON string
			rawPrivKeyParams, err := json.Marshal(&clusterKeyParams{
				Type: corePrivateKeyTypeP521,
				X:    key.X,
				Y:    key.Y,
				D:    key.D,
			})
			if err != nil {
				c.logger.Printf("[ERR] core: failed to encode local cluster key: %v", err)
				return err
			}

			// Store it
			err = c.barrier.Put(&Entry{
				Key:   coreLocalClusterKeyPath,
				Value: rawPrivKeyParams,
			})
			if err != nil {
				c.logger.Printf("[ERR] core: failed to store local cluster key: %v", err)
				return err
			}

			needNewCert = true
			c.localClusterPrivateKey = key
		} else {
			c.logger.Printf("[TRACE] core: local cluster private key found")
			var params clusterKeyParams
			if err = jsonutil.DecodeJSON(entry.Value, &params); err != nil {
				c.logger.Printf("[ERR] core: failed to decode local cluster key: %v", err)
				return err
			}
			switch {
			case params.X == nil, params.Y == nil, params.D == nil:
				c.logger.Printf("[ERR] core: failed to parse local cluster key due to missing params")
				return fmt.Errorf("failed to parse local cluster key")
			case params.Type == corePrivateKeyTypeP521:
			default:
				c.logger.Printf("[ERR] core: unknown local cluster key type %v", params.Type)
				return fmt.Errorf("failed to find valid local cluster key type")
			}
			c.localClusterPrivateKey = &ecdsa.PrivateKey{
				PublicKey: ecdsa.PublicKey{
					Curve: elliptic.P521(),
					X:     params.X,
					Y:     params.Y,
				},
				D: params.D,
			}
		}
	}

	if needNewCert || cluster.Certificate == nil || len(cluster.Certificate) == 0 {
		c.logger.Printf("[TRACE] core: generating new local cluster certificate")
		if c.localClusterPrivateKey == nil {
			c.logger.Printf("[ERR] core: generating a new local cluster certificate but private key is nil")
			return fmt.Errorf("failed to find private key when generating new local cluster certificate")
		}

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
			KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement | x509.KeyUsageCertSign,
			SerialNumber:          big.NewInt(mathrand.Int63()),
			NotBefore:             time.Now().Add(-30 * time.Second),
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

		cluster.Certificate = certBytes
		modified = true
	}
	c.localClusterCert = cluster.Certificate

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
		c.logger.Printf("[TRACE] core/startClusterListener: serving cluster requests on %s", tlsLn.Addr())
		go server.Serve(tlsLn)
	}

	c.clusterListenerShutdownCh = make(chan struct{})
	c.clusterListenerShutdownSuccessCh = make(chan struct{})

	go func() {
		<-c.clusterListenerShutdownCh
		c.logger.Printf("[TRACE] core/startClusterListener: shutting down listeners")
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
			"h2",
		},
	}

	return tlsConfig, nil
}

// refreshRequestForwardingConnection ensures that the client/transport are
// alive and that the current active address value matches the most
// recently-known address.
func (c *Core) refreshRequestForwardingConnection(clusterAddr string) error {
	c.logger.Printf("[TRACE] core/refreshRequestForwardingConnection: cluster address %s", clusterAddr)
	c.requestForwardingConnectionLock.RLock()

	if c.requestForwardingConnection == nil && clusterAddr == "" {
		c.requestForwardingConnectionLock.RUnlock()
		c.logger.Printf("[TRACE] core/refreshRequestForwardingConnection: no change (nil and empty)")
		return nil
	}
	if c.requestForwardingConnection != nil &&
		c.requestForwardingConnection.clusterAddr == clusterAddr {
		c.requestForwardingConnectionLock.RUnlock()
		c.logger.Printf("[TRACE] core/refreshRequestForwardingConnection: no changes")
		return nil
	}

	// Give up the read lock and get a write lock, then verify it's still required
	c.requestForwardingConnectionLock.RUnlock()
	c.requestForwardingConnectionLock.Lock()
	defer c.requestForwardingConnectionLock.Unlock()
	if c.requestForwardingConnection != nil &&
		c.requestForwardingConnection.clusterAddr == clusterAddr {
		return nil
	}

	// Disabled, potentially
	if clusterAddr == "" {
		c.requestForwardingConnection = nil
		return nil
	}

	if c.requestForwardingConnection == nil {
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
		}
	}

	if c.requestForwardingConnection.clusterAddr != clusterAddr {
		c.requestForwardingConnection.clusterAddr = clusterAddr
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

	freq, err := requestutil.GenerateForwardedRequest(req, c.requestForwardingConnection.clusterAddr+"/cluster/forwarded-request")
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
func WrapListenersForClustering(lns []net.Listener, handler http.Handler, logger *log.Logger) func() ([]net.Listener, http.Handler, error) {
	// This mux handles cluster functions (right now, only forwarded requests)
	mux := http.NewServeMux()
	mux.HandleFunc("/cluster/forwarded-request", func(w http.ResponseWriter, req *http.Request) {
		freq, err := requestutil.ParseForwardedRequest(req)
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
		ret := make([]net.Listener, 0, len(lns))
		// Loop over the existing listeners and start listeners on appropriate ports
		for _, ln := range lns {
			tcpAddr, ok := ln.Addr().(*net.TCPAddr)
			if !ok {
				if logger != nil {
					logger.Printf("[TRACE] http/WrapClusterListener: %s not a candidate for cluster request handling", ln.Addr().String())
				}
				continue
			}
			if logger != nil {
				logger.Printf("[TRACE] http/WrapClusterListener: %s is a candidate for cluster request handling at addr %s and port %d", tcpAddr.String(), tcpAddr.IP.String(), tcpAddr.Port+1)
			}

			ipStr := tcpAddr.IP.String()
			if len(tcpAddr.IP) == net.IPv6len {
				ipStr = fmt.Sprintf("[%s]", ipStr)
			}
			ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ipStr, tcpAddr.Port+1))
			if err != nil {
				return nil, nil, err
			}
			ret = append(ret, ln)
		}

		return ret, mux, nil
	}
}

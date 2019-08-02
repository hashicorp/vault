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
	"time"

	"github.com/hashicorp/errwrap"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/cluster"
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

// startClusterListener starts cluster request listeners during unseal. It
// is assumed that the state lock is held while this is run. Right now this
// only starts cluster listeners. Once the listener is started handlers/clients
// can start being registered to it.
func (c *Core) startClusterListener(ctx context.Context) error {
	if c.ClusterAddr() == "" {
		c.logger.Info("clustering disabled, not starting listeners")
		return nil
	}

	if c.clusterListener != nil {
		c.logger.Warn("cluster listener is already started")
		return nil
	}

	if c.clusterListenerAddrs == nil || len(c.clusterListenerAddrs) == 0 {
		c.logger.Warn("clustering not disabled but no addresses to listen on")
		return fmt.Errorf("cluster addresses not found")
	}

	c.logger.Debug("starting cluster listeners")

	c.clusterListener = cluster.NewListener(c.clusterListenerAddrs, c.clusterCipherSuites, c.logger.Named("cluster-listener"))

	err := c.clusterListener.Run(ctx)
	if err != nil {
		return err
	}
	if strings.HasSuffix(c.ClusterAddr(), ":0") {
		// If we listened on port 0, record the port the OS gave us.
		c.clusterAddr.Store(fmt.Sprintf("https://%s", c.clusterListener.Addrs()[0]))
	}
	return nil
}

func (c *Core) ClusterAddr() string {
	return c.clusterAddr.Load().(string)
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
	c.clusterListener = nil

	c.logger.Info("cluster listeners successfully shut down")
}

func (c *Core) SetClusterListenerAddrs(addrs []*net.TCPAddr) {
	c.clusterListenerAddrs = addrs
	if c.ClusterAddr() == "" && len(addrs) == 1 {
		c.clusterAddr.Store(fmt.Sprintf("https://%s", addrs[0].String()))
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

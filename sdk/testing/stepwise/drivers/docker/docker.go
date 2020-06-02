package docker

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	mathrand "math/rand"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/internalshared/reloadutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/testing/stepwise"
	"github.com/hashicorp/vault/vault"
	"golang.org/x/net/http2"

	docker "github.com/docker/docker/client"
	uuid "github.com/hashicorp/go-uuid"
)

var _ stepwise.StepDriver = &DockerCluster{}

// DockerCluster is used to managing the lifecycle of the test Vault cluster
type DockerCluster struct {
	// PluginName is the input from the test case
	PluginName string
	// ClusterName is a UUID name of the cluster. Docker ID?
	CluterName string

	// DriverOptions are a set of options from the Stepwise test using this
	// cluster
	DriverOptions stepwise.DriverOptions

	RaftStorage        bool
	ClientAuthRequired bool
	BarrierKeys        [][]byte
	RecoveryKeys       [][]byte
	CACertBytes        []byte
	CACertPEM          []byte
	CAKeyPEM           []byte
	CACertPEMFile      string
	ID                 string
	TempDir            string
	ClusterName        string
	RootCAs            *x509.CertPool
	CACert             *x509.Certificate
	CAKey              *ecdsa.PrivateKey
	CleanupFunc        func()
	SetupFunc          func()
	ClusterNodes       []*DockerClusterNode

	rootToken string
}

// Teardown stops all the containers.
// TODO: error/logging
func (rc *DockerCluster) Teardown() error {
	for _, node := range rc.ClusterNodes {
		node.Cleanup()
	}
	// TODO should return something
	return nil
}

// ExpandPath accepts a user supplied path and prefixes it with the required
// mount, auth, and namespace parts as needed. This enables drivers to create
// random/unique mount and namespace paths while allowing tests to be written
// with relative paths.
// TODO should this be in stepwise and not per-driver? Maybe have Drivers give
// mount, namespace returns and have stepwise do this so people aren't
// re-inventing this
func (dc *DockerCluster) ExpandPath(path string) string {
	// TODO mount point
	newPath := path
	newPath = fmt.Sprintf("%s/%s", dc.DriverOptions.MountPath, newPath)
	if dc.DriverOptions.PluginType == stepwise.PluginTypeCredential {
		newPath = fmt.Sprintf("%s/%s", "auth", newPath)
		// TODO prefix with namespace
	}
	return newPath
}

// MountPath returns the path that the plugin under test is mounted at
// TODO: include namespace support
func (dc *DockerCluster) MountPath() string {
	if dc.DriverOptions.PluginType == stepwise.PluginTypeCredential {
		return fmt.Sprintf("%s/%s", "auth", dc.DriverOptions.MountPath)
	}
	return dc.DriverOptions.MountPath
}

// RootToken returns the root token of the cluster, if set
func (dc *DockerCluster) RootToken() string {
	return dc.rootToken
}

func (dc *DockerCluster) Name() string {
	// TODO return UUID cluster name
	return dc.PluginName
}

func (dc *DockerCluster) Client() (*api.Client, error) {
	if len(dc.ClusterNodes) > 0 {
		if dc.ClusterNodes[0].Client != nil {
			// TODO is clone needed here?
			c, err := dc.ClusterNodes[0].Client.Clone()
			if err != nil {
				return nil, err
			}
			c.SetToken(dc.ClusterNodes[0].Client.Token())
			return c, nil
			// return dc.ClusterNodes[0].Client, nil
		}
	}

	return nil, errors.New("no configured client found")
}

func (rc *DockerCluster) GetBarrierOrRecoveryKeys() [][]byte {
	return rc.GetBarrierKeys()
}

func (rc *DockerCluster) GetCACertPEMFile() string {
	return rc.CACertPEMFile
}

func (rc *DockerCluster) ClusterID() string {
	return rc.ID
}

func (n *DockerClusterNode) Name() string {
	return n.Cluster.ClusterName + "-" + n.NodeID
}

type VaultClusterNode interface {
	Name() string
	APIClient() *api.Client
}

func (rc *DockerCluster) Nodes() []VaultClusterNode {
	ret := make([]VaultClusterNode, len(rc.ClusterNodes))
	for i, core := range rc.ClusterNodes {
		ret[i] = core
	}
	return ret
}

func (rc *DockerCluster) GetBarrierKeys() [][]byte {
	ret := make([][]byte, len(rc.BarrierKeys))
	for i, k := range rc.BarrierKeys {
		ret[i] = vault.TestKeyCopy(k)
	}
	return ret
}

func (rc *DockerCluster) GetRecoveryKeys() [][]byte {
	ret := make([][]byte, len(rc.RecoveryKeys))
	for i, k := range rc.RecoveryKeys {
		ret[i] = vault.TestKeyCopy(k)
	}
	return ret
}

func (rc *DockerCluster) SetBarrierKeys(keys [][]byte) {
	rc.BarrierKeys = make([][]byte, len(keys))
	for i, k := range keys {
		rc.BarrierKeys[i] = vault.TestKeyCopy(k)
	}
}

func (rc *DockerCluster) SetRecoveryKeys(keys [][]byte) {
	rc.RecoveryKeys = make([][]byte, len(keys))
	for i, k := range keys {
		rc.RecoveryKeys[i] = vault.TestKeyCopy(k)
	}
}

func (rc *DockerCluster) Initialize(ctx context.Context) error {
	client, err := rc.ClusterNodes[0].CreateAPIClient()
	if err != nil {
		return err
	}

	var resp *api.InitResponse
	for ctx.Err() == nil {
		resp, err = client.Sys().Init(&api.InitRequest{
			SecretShares:    3,
			SecretThreshold: 3,
		})
		if err == nil && resp != nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		return err
	}
	if resp == nil {
		return fmt.Errorf("nil response to init request")
	}

	for _, k := range resp.Keys {
		raw, err := hex.DecodeString(k)
		if err != nil {
			return err
		}
		rc.BarrierKeys = append(rc.BarrierKeys, raw)
	}
	for _, k := range resp.RecoveryKeys {
		raw, err := hex.DecodeString(k)
		if err != nil {
			return err
		}
		rc.RecoveryKeys = append(rc.RecoveryKeys, raw)
	}
	rc.rootToken = resp.RootToken

	// Write root token and barrier keys
	err = ioutil.WriteFile(filepath.Join(rc.TempDir, "root_token"), []byte(rc.rootToken), 0755)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	for _, key := range rc.BarrierKeys {
		// TODO handle errors?
		_, _ = buf.Write(key)
		_, _ = buf.WriteRune('\n')
	}
	err = ioutil.WriteFile(filepath.Join(rc.TempDir, "barrier_keys"), buf.Bytes(), 0755)
	if err != nil {
		return err
	}
	for _, key := range rc.RecoveryKeys {
		// TODO handle errors?
		_, _ = buf.Write(key)
		_, _ = buf.WriteRune('\n')
	}
	err = ioutil.WriteFile(filepath.Join(rc.TempDir, "recovery_keys"), buf.Bytes(), 0755)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Unseal
	for j, node := range rc.ClusterNodes {
		// copy the index value, so we're not reusing it in deeper scopes
		i := j
		client, err := node.CreateAPIClient()
		if err != nil {
			return err
		}
		node.Client = client

		if i > 0 && rc.RaftStorage {
			leader := rc.ClusterNodes[0]
			resp, err := client.Sys().RaftJoin(&api.RaftJoinRequest{
				LeaderAPIAddr:    fmt.Sprintf("https://%s:%d", rc.ClusterNodes[0].Name(), leader.Address.Port),
				LeaderCACert:     string(rc.CACertPEM),
				LeaderClientCert: string(node.ServerCertPEM),
				LeaderClientKey:  string(node.ServerKeyPEM),
			})
			if err != nil {
				return err
			}
			if resp == nil || !resp.Joined {
				return fmt.Errorf("nil or negative response from raft join request: %v", resp)
			}
		}

		var unsealed bool
		for _, key := range rc.BarrierKeys {
			resp, err := client.Sys().Unseal(hex.EncodeToString(key))
			if err != nil {
				return err
			}
			unsealed = !resp.Sealed
		}
		if i == 0 && !unsealed {
			return fmt.Errorf("could not unseal node %d", i)
		}
		client.SetToken(rc.rootToken)

		err = TestWaitHealthMatches(ctx, node.Client, func(health *api.HealthResponse) error {
			if health.Sealed {
				return fmt.Errorf("node %d is sealed: %#v", i, health)
			}
			if health.ClusterID == "" {
				return fmt.Errorf("node %d has no cluster ID", i)
			}

			rc.ID = health.ClusterID
			return nil
		})
		if err != nil {
			return err
		}

		if i == 0 {
			err = TestWaitLeaderMatches(ctx, node.Client, func(leader *api.LeaderResponse) error {
				if !leader.IsSelf {
					return fmt.Errorf("node %d leader=%v, expected=%v", i, leader.IsSelf, true)
				}

				return nil
			})
			if err != nil {
				return err
			}
		}
	}

	for i, node := range rc.ClusterNodes {
		expectLeader := i == 0
		err = TestWaitLeaderMatches(ctx, node.Client, func(leader *api.LeaderResponse) error {
			if expectLeader != leader.IsSelf {
				return fmt.Errorf("node %d leader=%v, expected=%v", i, leader.IsSelf, expectLeader)
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (rc *DockerCluster) setupCA(opts *DockerClusterOptions) error {
	var err error

	certIPs := []net.IP{
		net.IPv6loopback,
		net.ParseIP("127.0.0.1"),
	}

	var caKey *ecdsa.PrivateKey
	if opts != nil && opts.CAKey != nil {
		caKey = opts.CAKey
	} else {
		caKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return err
		}
	}
	rc.CAKey = caKey

	var caBytes []byte
	if opts != nil && len(opts.CACert) > 0 {
		caBytes = opts.CACert
	} else {
		caCertTemplate := &x509.Certificate{
			Subject: pkix.Name{
				CommonName: "localhost",
			},
			DNSNames:              []string{"localhost"},
			IPAddresses:           certIPs,
			KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
			SerialNumber:          big.NewInt(mathrand.Int63()),
			NotBefore:             time.Now().Add(-30 * time.Second),
			NotAfter:              time.Now().Add(262980 * time.Hour),
			BasicConstraintsValid: true,
			IsCA:                  true,
		}
		caBytes, err = x509.CreateCertificate(rand.Reader, caCertTemplate, caCertTemplate, caKey.Public(), caKey)
		if err != nil {
			return err
		}
	}
	caCert, err := x509.ParseCertificate(caBytes)
	if err != nil {
		return err
	}
	rc.CACert = caCert
	rc.CACertBytes = caBytes

	rc.RootCAs = x509.NewCertPool()
	rc.RootCAs.AddCert(caCert)

	caCertPEMBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	}
	rc.CACertPEM = pem.EncodeToMemory(caCertPEMBlock)

	rc.CACertPEMFile = filepath.Join(rc.TempDir, "ca", "ca.pem")
	err = ioutil.WriteFile(rc.CACertPEMFile, rc.CACertPEM, 0755)
	if err != nil {
		return err
	}

	marshaledCAKey, err := x509.MarshalECPrivateKey(caKey)
	if err != nil {
		return err
	}
	caKeyPEMBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: marshaledCAKey,
	}
	rc.CAKeyPEM = pem.EncodeToMemory(caKeyPEMBlock)

	// We don't actually need this file, but it may be helpful for debugging.
	err = ioutil.WriteFile(filepath.Join(rc.TempDir, "ca", "ca_key.pem"), rc.CAKeyPEM, 0755)
	if err != nil {
		return err
	}

	return nil
}

// TODO: unused at this point
// func (rc *DockerCluster) raftJoinConfig() []api.RaftJoinRequest {
// 	ret := make([]api.RaftJoinRequest, len(rc.ClusterNodes))
// 	for _, node := range rc.ClusterNodes {
// 		ret = append(ret, api.RaftJoinRequest{
// 			LeaderAPIAddr:    fmt.Sprintf("https://%s:%d", node.Address.IP, node.Address.Port),
// 			LeaderCACert:     string(rc.CACertPEM),
// 			LeaderClientCert: string(node.ServerCertPEM),
// 			LeaderClientKey:  string(node.ServerKeyPEM),
// 		})
// 	}
// 	return ret
// }

// Don't call this until n.Address.IP is populated
func (n *DockerClusterNode) setupCert() error {
	var err error

	n.ServerKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	certTemplate := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: n.Name(),
		},
		// Include host.docker.internal for the sake of benchmark-vault running on MacOS/Windows.
		// This allows Prometheus running in docker to scrape the cluster for metrics.
		DNSNames:    []string{"localhost", "host.docker.internal", n.Name()},
		IPAddresses: []net.IP{net.IPv6loopback, net.ParseIP("127.0.0.1")}, // n.Address.IP,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
		SerialNumber: big.NewInt(mathrand.Int63()),
		NotBefore:    time.Now().Add(-30 * time.Second),
		NotAfter:     time.Now().Add(262980 * time.Hour),
	}
	n.ServerCertBytes, err = x509.CreateCertificate(rand.Reader, certTemplate, n.Cluster.CACert, n.ServerKey.Public(), n.Cluster.CAKey)
	if err != nil {
		return err
	}
	n.ServerCert, err = x509.ParseCertificate(n.ServerCertBytes)
	if err != nil {
		return err
	}
	n.ServerCertPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: n.ServerCertBytes,
	})

	marshaledKey, err := x509.MarshalECPrivateKey(n.ServerKey)
	if err != nil {
		return err
	}
	n.ServerKeyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: marshaledKey,
	})

	n.ServerCertPEMFile = filepath.Join(n.WorkDir, "cert.pem")
	err = ioutil.WriteFile(n.ServerCertPEMFile, n.ServerCertPEM, 0755)
	if err != nil {
		return err
	}

	n.ServerKeyPEMFile = filepath.Join(n.WorkDir, "key.pem")
	err = ioutil.WriteFile(n.ServerKeyPEMFile, n.ServerKeyPEM, 0755)
	if err != nil {
		return err
	}

	tlsCert, err := tls.X509KeyPair(n.ServerCertPEM, n.ServerKeyPEM)
	if err != nil {
		return err
	}

	certGetter := reloadutil.NewCertificateGetter(n.ServerCertPEMFile, n.ServerKeyPEMFile, "")
	if err := certGetter.Reload(nil); err != nil {
		// TODO error handle or panic?
		panic(err)
	}
	tlsConfig := &tls.Config{
		Certificates:   []tls.Certificate{tlsCert},
		RootCAs:        n.Cluster.RootCAs,
		ClientCAs:      n.Cluster.RootCAs,
		ClientAuth:     tls.RequestClientCert,
		NextProtos:     []string{"h2", "http/1.1"},
		GetCertificate: certGetter.GetCertificate,
	}

	if n.Cluster.ClientAuthRequired {
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}
	n.TLSConfig = tlsConfig

	return nil
}

type DockerClusterNode struct {
	NodeID            string
	Address           *net.TCPAddr
	HostPort          string
	Client            *api.Client
	ServerCert        *x509.Certificate
	ServerCertBytes   []byte
	ServerCertPEM     []byte
	ServerCertPEMFile string
	ServerKey         *ecdsa.PrivateKey
	ServerKeyPEM      []byte
	ServerKeyPEMFile  string
	TLSConfig         *tls.Config
	WorkDir           string
	Cluster           *DockerCluster
	container         *types.ContainerJSON
	dockerAPI         *docker.Client
}

func (n *DockerClusterNode) APIClient() *api.Client {
	return n.Client
}

func (n *DockerClusterNode) CreateAPIClient() (*api.Client, error) {
	transport := cleanhttp.DefaultPooledTransport()
	transport.TLSClientConfig = n.TLSConfig.Clone()
	if err := http2.ConfigureTransport(transport); err != nil {
		return nil, err
	}
	client := &http.Client{
		Transport: transport,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			// This can of course be overridden per-test by using its own client
			return fmt.Errorf("redirects not allowed in these tests")
		},
	}
	config := api.DefaultConfig()
	if config.Error != nil {
		return nil, config.Error
	}
	config.Address = fmt.Sprintf("https://127.0.0.1:%s", n.HostPort)
	config.HttpClient = client
	config.MaxRetries = 0
	apiClient, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	apiClient.SetToken(n.Cluster.RootToken())
	return apiClient, nil
}

func (n *DockerClusterNode) Cleanup() {
	if err := n.dockerAPI.ContainerKill(context.Background(), n.container.ID, "KILL"); err != nil {
		panic(err)
	}
}

func (n *DockerClusterNode) Start(cli *docker.Client, caDir, netName string, netCIDR *DockerClusterNode, pluginBinPath string) error {
	n.dockerAPI = cli

	err := n.setupCert()
	if err != nil {
		return err
	}

	vaultCfg := map[string]interface{}{
		"listener": map[string]interface{}{
			"tcp": map[string]interface{}{
				"address":       fmt.Sprintf("%s:%d", "0.0.0.0", 8200),
				"tls_cert_file": "/vault/config/cert.pem",
				"tls_key_file":  "/vault/config/key.pem",
				"telemetry": map[string]interface{}{
					"unauthenticated_metrics_access": true,
				},
			},
		},
		"telemetry": map[string]interface{}{
			"disable_hostname": true,
		},
		"storage": map[string]interface{}{
			"raft": map[string]interface{}{
				"path":    "/vault/file",
				"node_id": n.NodeID,
				//"retry_join": string(joinConfigStr),
			},
		},
		"cluster_name":         netName,
		"log_level":            "TRACE",
		"raw_storage_endpoint": true,
		"plugin_directory":     "/vault/config",
		//TODO: mlock
		"disable_mlock": true,
		// These are being provided by docker-entrypoint now, since we don't know
		// the address before the container starts.
		//"api_addr": fmt.Sprintf("https://%s:%d", n.Address.IP, n.Address.Port),
		//"cluster_addr": fmt.Sprintf("https://%s:%d", n.Address.IP, n.Address.Port+1),
	}
	cfgJSON, err := json.Marshal(vaultCfg)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(n.WorkDir, "local.json"), cfgJSON, 0644)
	if err != nil {
		return err
	}
	// setup plugin bin copy if needed
	copyFromTo := map[string]string{
		n.WorkDir: "/vault/config",
		caDir:     "/usr/local/share/ca-certificates/",
	}
	if pluginBinPath != "" {
		base := path.Base(pluginBinPath)
		copyFromTo[pluginBinPath] = filepath.Join("/vault/config", base)
	}

	r := &Runner{
		dockerAPI: cli,
		ContainerConfig: &container.Config{
			Image: "vault",
			Entrypoint: []string{"/bin/sh", "-c", "update-ca-certificates && " +
				"exec /usr/local/bin/docker-entrypoint.sh vault server -log-level=trace -dev-plugin-dir=./vault/config -config /vault/config/local.json"},
			Env: []string{
				"VAULT_CLUSTER_INTERFACE=eth0",
				// TODO: api addr set is funny
				"VAULT_API_ADDR=https://127.0.0.1:8200",
				fmt.Sprintf("VAULT_REDIRECT_ADDR=https://%s:8200", n.Name()),
			},
			Labels:       nil,
			ExposedPorts: nat.PortSet{"8200/tcp": {}, "8201/tcp": {}},
		},
		ContainerName: n.Name(),
		NetName:       netName,
		CopyFromTo:    copyFromTo,
	}

	n.container, err = r.Start(context.Background())
	if err != nil {
		return err
	}

	n.Address = &net.TCPAddr{
		IP:   net.ParseIP(n.container.NetworkSettings.IPAddress),
		Port: 8200,
	}
	ports := n.container.NetworkSettings.NetworkSettingsBase.Ports[nat.Port("8200/tcp")]
	if len(ports) == 0 {
		n.Cleanup()
		return fmt.Errorf("could not find port binding for 8200/tcp")
	}
	n.HostPort = ports[0].HostPort

	return nil
}

// DockerClusterOptions has options for setting up the docker cluster
type DockerClusterOptions struct {
	KeepStandbysSealed bool
	RequireClientAuth  bool
	SkipInit           bool
	CACert             []byte
	NumCores           int
	TempDir            string
	PluginTestBin      string
	// SetupFunc is called after the cluster is started.
	SetupFunc func(t testing.T, c *DockerCluster)
	CAKey     *ecdsa.PrivateKey
	// TODO: plugin source dir here?
}

//
// test methods/functions
//

// TestWaitHealthMatches checks health TODO: update docs
func TestWaitHealthMatches(ctx context.Context, client *api.Client, ready func(response *api.HealthResponse) error) error {
	var health *api.HealthResponse
	var err error
	for ctx.Err() == nil {
		// TODO ideally Health method would take a context
		health, err = client.Sys().Health()
		switch {
		case err != nil:
		case health == nil:
			err = fmt.Errorf("nil response to health check")
		default:
			err = ready(health)
			if err == nil {
				return nil
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("error checking health: %v", err)
}

func TestWaitLeaderMatches(ctx context.Context, client *api.Client, ready func(response *api.LeaderResponse) error) error {
	var leader *api.LeaderResponse
	var err error
	for ctx.Err() == nil {
		// TODO ideally Leader method would take a context
		leader, err = client.Sys().Leader()
		switch {
		case err != nil:
		case leader == nil:
			err = fmt.Errorf("nil response to leader check")
		default:
			err = ready(leader)
			if err == nil {
				return nil
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("error checking leader: %v", err)
}

// end test helper methods

// TODO: change back to 3
// var DefaultNumCores = 3

var DefaultNumCores = 3

// creates a managed docker container running Vault
func (cluster *DockerCluster) setupDockerCluster(base *vault.CoreConfig, opts *DockerClusterOptions) error {

	if opts != nil && opts.TempDir != "" {
		if _, err := os.Stat(opts.TempDir); os.IsNotExist(err) {
			if err := os.MkdirAll(opts.TempDir, 0700); err != nil {
				return err
			}
		}
		cluster.TempDir = opts.TempDir
	} else {
		tempDir, err := ioutil.TempDir("", "vault-test-cluster-")
		if err != nil {
			return err
		}
		cluster.TempDir = tempDir
	}
	caDir := filepath.Join(cluster.TempDir, "ca")
	if err := os.MkdirAll(caDir, 0755); err != nil {
		return err
	}

	var numCores int
	if opts == nil || opts.NumCores == 0 {
		numCores = DefaultNumCores
	} else {
		numCores = opts.NumCores
	}

	if opts != nil && opts.RequireClientAuth {
		cluster.ClientAuthRequired = true
	}

	cidr := "192.168.128.0/20"
	for i := 0; i < numCores; i++ {
		nodeID := fmt.Sprintf("vault-%d", i)
		node := &DockerClusterNode{
			NodeID:  nodeID,
			Cluster: cluster,
			WorkDir: filepath.Join(cluster.TempDir, nodeID),
		}
		cluster.ClusterNodes = append(cluster.ClusterNodes, node)
		if err := os.MkdirAll(node.WorkDir, 0700); err != nil {
			return err
		}
	}

	err := cluster.setupCA(opts)
	if err != nil {
		return err
	}

	cli, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithVersion("1.40"))
	if err != nil {
		return err
	}
	netName := "vault-test"
	_, err = SetupNetwork(cli, netName, cidr)
	if err != nil {
		return err
	}

	for _, node := range cluster.ClusterNodes {
		pluginBinPath := ""
		if opts != nil {
			pluginBinPath = opts.PluginTestBin
		}

		err := node.Start(cli, caDir, netName, node, pluginBinPath)
		if err != nil {
			return err
		}
	}

	if opts == nil || !opts.SkipInit {
		if err := cluster.Initialize(context.Background()); err != nil {
			return err
		}
	}

	return nil
}

// Docker networking functions
// SetupNetwork establishes networking for the Docker container
func SetupNetwork(cli *docker.Client, netName, cidr string) (string, error) {
	ctx := context.Background()

	netResources, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return "", err
	}
	for _, netRes := range netResources {
		if netRes.Name == netName {
			if len(netRes.IPAM.Config) > 0 && netRes.IPAM.Config[0].Subnet == cidr {
				return netRes.ID, nil
			}
			err = cli.NetworkRemove(ctx, netRes.ID)
			if err != nil {
				return "", err
			}
		}
	}

	id, err := createNetwork(cli, netName, cidr)
	if err != nil {
		return "", fmt.Errorf("couldn't create network %s on %s: %w", netName, cidr, err)
	}
	return id, nil
}

func createNetwork(cli *docker.Client, netName, cidr string) (string, error) {
	resp, err := cli.NetworkCreate(context.Background(), netName, types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         "bridge",
		Options:        map[string]string{},
		IPAM: &network.IPAM{
			Driver:  "default",
			Options: map[string]string{},
			Config: []network.IPAMConfig{
				{
					Subnet: cidr,
				},
			},
		},
	})
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

// NewDockerDriver creats a new Stepwise Driver for executing tests
func NewDockerDriver(name string, do *stepwise.DriverOptions) *DockerCluster {
	// TODO name here should be name of the test?
	clusterUUID, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}

	if do == nil {
		// set empty values
		do = &stepwise.DriverOptions{}
	}
	return &DockerCluster{
		// TODO should this be the driver options plugin name?
		PluginName:    name,
		ClusterName:   fmt.Sprintf("test-%s-%s", name, clusterUUID),
		RaftStorage:   true,
		DriverOptions: *do,
	}
}

// Setup creates any temp dir and compiles the binary for copying to Docker
func (dc *DockerCluster) Setup() error {
	// TODO many not use name here
	name := dc.DriverOptions.PluginName
	// TODO make PluginName give random name with prefix if given
	pluginName := dc.DriverOptions.PluginName
	// get the working directory of the plugin being tested.
	srcDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// tmpDir gets cleaned up when the cluster is cleaned up
	tmpDir, err := ioutil.TempDir("", "bin")
	if err != nil {
		log.Fatal(err)
	}

	binName, binPath, sha256value, err := stepwise.CompilePlugin(name, pluginName, srcDir, tmpDir)
	if err != nil {
		panic(err)
	}

	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
	}

	dOpts := &DockerClusterOptions{PluginTestBin: binPath}
	err = dc.setupDockerCluster(coreConfig, dOpts)
	if err != nil {
		panic(err)
	}

	cores := dc.ClusterNodes
	client := cores[0].Client

	// use client to mount plugin
	err = client.Sys().RegisterPlugin(&api.RegisterPluginInput{
		Name:    name,
		Type:    consts.PluginType(dc.DriverOptions.PluginType),
		Command: binName,
		SHA256:  sha256value,
	})
	if err != nil {
		panic(err)
	}

	var mountErr error
	switch dc.DriverOptions.PluginType {
	case stepwise.PluginTypeCredential:
		mountErr = client.Sys().EnableAuthWithOptions(dc.DriverOptions.MountPath, &api.EnableAuthOptions{
			Type: name,
		})
	case stepwise.PluginTypeDatabase:
		// TODO database type
		return errors.New("plugin type database not yet supported")
	case stepwise.PluginTypeSecrets:
		mountErr = client.Sys().Mount(dc.DriverOptions.MountPath, &api.MountInput{
			Type: name,
		})
	default:
		return fmt.Errorf("unknown plugin type: %s", dc.DriverOptions.PluginType.String())
	}
	if mountErr != nil {
		panic(mountErr)
	}
	return mountErr
}

// Runner manages the lifecycle of the Docker container
type Runner struct {
	dockerAPI       *docker.Client
	ContainerConfig *container.Config
	ContainerName   string
	NetName         string
	IP              string
	CopyFromTo      map[string]string
}

func (d *Runner) Start(ctx context.Context) (*types.ContainerJSON, error) {
	hostConfig := &container.HostConfig{
		PublishAllPorts: true,
		// TODO: configure auto remove
		// AutoRemove:      false,
		AutoRemove: true,
	}

	networkingConfig := &network.NetworkingConfig{}
	switch d.NetName {
	case "":
	case "host":
		hostConfig.NetworkMode = "host"
	default:
		es := &network.EndpointSettings{
			//Links:               nil,
			Aliases: []string{d.ContainerName},
		}
		if len(d.IP) != 0 {
			es.IPAMConfig = &network.EndpointIPAMConfig{
				IPv4Address: d.IP,
			}
		}
		networkingConfig.EndpointsConfig = map[string]*network.EndpointSettings{
			d.NetName: es,
		}
	}

	// best-effort pull
	resp, _ := d.dockerAPI.ImageCreate(ctx, d.ContainerConfig.Image, types.ImageCreateOptions{})
	if resp != nil {
		_, _ = ioutil.ReadAll(resp)
	}

	cfg := *d.ContainerConfig
	hostConfig.CapAdd = strslice.StrSlice{"IPC_LOCK"}
	cfg.Hostname = d.ContainerName
	//fullName := d.NetName + "." + d.ContainerName
	fullName := d.ContainerName
	container, err := d.dockerAPI.ContainerCreate(ctx, &cfg, hostConfig, networkingConfig, fullName)
	if err != nil {
		return nil, fmt.Errorf("container create failed: %v", err)
	}

	for from, to := range d.CopyFromTo {
		srcInfo, err := archive.CopyInfoSourcePath(from, false)
		if err != nil {
			return nil, fmt.Errorf("error copying from source %q: %v", from, err)
		}

		srcArchive, err := archive.TarResource(srcInfo)
		if err != nil {
			return nil, fmt.Errorf("error creating tar from source %q: %v", from, err)
		}
		defer srcArchive.Close()

		dstInfo := archive.CopyInfo{Path: to}

		dstDir, content, err := archive.PrepareArchiveCopy(srcArchive, srcInfo, dstInfo)
		if err != nil {
			return nil, fmt.Errorf("error preparing copy from %q -> %q: %v", from, to, err)
		}
		defer content.Close()
		err = d.dockerAPI.CopyToContainer(ctx, container.ID, dstDir, content, types.CopyToContainerOptions{})
		if err != nil {
			return nil, fmt.Errorf("error copying from %q -> %q: %v", from, to, err)
		}
	}

	err = d.dockerAPI.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, fmt.Errorf("container start failed: %v", err)
	}

	inspect, err := d.dockerAPI.ContainerInspect(ctx, container.ID)
	if err != nil {
		return nil, err
	}
	return &inspect, nil
}
